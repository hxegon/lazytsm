package project

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

type tmuxStatus int

const (
	runningInside tmuxStatus = iota
	runningNotInside
	notRunning
	unknown
)

type Tmux struct {
	status  tmuxStatus
	binPath string
	error   error
}

// Wrap the pgrep command
func pgrep(query string) string {
	cmd := exec.Command("pgrep", query)
	out, _ := cmd.CombinedOutput()
	return string(out)
}

func getTmuxStatus() (tmuxStatus, error) {
	tmuxRunning := pgrep("tmux") != ""

	if tmuxRunning {
		if os.Getenv("TMUX") != "" {
			return runningInside, nil
		} else {
			return runningNotInside, nil
		}
	}

	return notRunning, nil
}

func NewTmux() (Tmux, error) {
	binPath, err := exec.LookPath("tmux")
	if err == nil {
		binPath, _ = filepath.Abs(binPath) // TODO
	} else {
		return Tmux{}, fmt.Errorf("Couldn't find tmux executable in $PATH: %v", err)
	}

	st, err := getTmuxStatus()
	tm := Tmux{
		status:  st,
		binPath: binPath,
		error:   err,
	}

	return tm, err
}

// Checks if there's a tmux session with a matching name. Always returns false if tmux isn't running
func (t Tmux) HasSession(query string) (bool, error) {
	if t.status == notRunning || t.status == unknown {
		return false, nil
	}

	// TODO: injection vuln issue?
	tArg := fmt.Sprintf("-t=%v", query)

	cmd := exec.Command("tmux", "has-session", tArg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}

	if len(out) == 0 {
		return true, nil
	}
	return false, nil
}

// This function will end the current process
// NOTE: This sets the name of the tmux session to the name of target
// !!!: THIS ENDS THE CURRENT GO PROGRAM AND BECOMES THE TMUX PROGRAM!!!
func (t Tmux) OpenOrSwitchTmuxSession(target, cwdPath string) error {
	// TODO: Break this up once I understand it better, as of now it's a ripoff of the old
	// attach/new-session logic in the tms script
	var tmuxCmd []string
	switch t.status {
	case notRunning:
		tmuxCmd = []string{"tmux", "new-session", "-s", target, "-c", cwdPath}
	case runningNotInside, runningInside:
		if has, _ := t.HasSession(target); !has {
			// Make session without switching to it
			exec.Command("tmux", "new-session", "-ds", target, "-c", cwdPath).Run() // .Error()
		}
		// and switch to it
		if t.status == runningInside {
			fmt.Println("switch-client")
			tmuxCmd = []string{"tmux", "switch-client", "-t", target}
		} else {
			fmt.Println("attach")
			tmuxCmd = []string{"tmux", "attach", "-t", target}
		}
	case unknown:
		panic(fmt.Sprintf("idk, unknown Tmux.status", t.error))
	}

	// low level call so the go program "becomes" the tmux command
	// I don't think tmux will inherit your env without the os.Environ() here
	err := syscall.Exec(t.binPath, tmuxCmd, os.Environ())

	return fmt.Errorf("lazytsm couldn't become tmux command. This should never happen!: %v", err)
}
