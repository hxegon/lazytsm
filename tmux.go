package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// TODO: Refactor so that if something changes about the TMux session during runtime, it won't cause issues
// TODO: figure out how to replace pgrep dependency

type tmuxStatus int

const (
	runningInside tmuxStatus = iota
	runningNotInside
	notRunning
	unknown
)

type Tmux struct {
	error   error
	binPath string
	status  tmuxStatus
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
		binPath, _ = filepath.Abs(binPath)
	} else {
		return Tmux{}, fmt.Errorf("couldn't find tmux executable in $PATH: %v", err)
	}

	st, err := getTmuxStatus()
	tm := Tmux{
		status:  st,
		binPath: binPath,
		error:   err,
	}

	return tm, err
}

func (t Tmux) SessionNames() []string {
	return []string{}
}

// Checks if there's a tmux session with a matching name. Always returns false if tmux isn't running
func (t Tmux) HasSession(query string) (bool, error) {
	if t.status == notRunning || t.status == unknown {
		return false, nil
	}

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

// Returns the name of the current session, or an empty string if there isn't one.
//
// i.e. We aren't in a session or tmux isn't running or something
func (t *Tmux) CurrentSessionName() string {
	// There's no current session name if we aren't in a tmux session.
	if t.status != runningInside {
		return ""
	}

	cmd := exec.Command("tmux", "display-message", "-p", "#S")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic("Command to get tmux session named errored in an unexpected way.")
	}

	return strings.TrimRight(string(out), "\n")
}

/*
This function will end the current process

NOTE: This sets the name of the tmux session to the name of target

!!!: THIS ENDS THE CURRENT GO PROGRAM AND BECOMES THE TMUX PROGRAM!!!
*/
func (t Tmux) OpenOrSwitchTmuxSession(target, cwdPath string) error {
	// attach/new-session logic in the tms script
	var tmuxCmd []string
	switch t.status {
	case notRunning:
		tmuxCmd = []string{"tmux", "new-session", "-s", target, "-c", cwdPath}
	case runningNotInside, runningInside:
		if has, _ := t.HasSession(target); !has {
			// Make session without switching to it
			exec.Command("tmux", "new-session", "-ds", target, "-c", cwdPath).Run() // TODO: Error handling
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
		panic(fmt.Sprintf("idk, unknown Tmux.status %v", t.error))
	}

	// low level call so the go program "becomes" the tmux command
	// I don't think tmux will inherit your env without the os.Environ() here
	err := syscall.Exec(t.binPath, tmuxCmd, os.Environ())

	return fmt.Errorf("lazytsm couldn't become tmux command. This should never happen!: %v", err)
}
