# A TUI for managing tmux sessions.

![](https://github.com/hxegon/lazytsm/blob/main/show.gif)

Extremely pre-alpha, but it sets the groundwork for a lot of functionality I've
wanted to add to the old `tms` script
https://github.com/ThePrimeagen/.dotfiles/blob/master/bin/.local/scripts/tmux-sessionizer

This is written by a Go newbie so any feedback / code review would be much
appreciated.

## Installation

Just clone and `go build cmd/lazytsm.go`, or `just build` There's a `flake.nix`
and a `.envrc` with `use flake` if you're into that sort of thing

## Usage

Make a `.lazytsm.toml` in your `$HOME` folder that looks something like this:

```toml
gitSearchDirs = [
"/home/hxegon/Code"
]
extraDirs = [
"/home/hxegon/dotfiles",
"/home/hxegon/scripts",
"/home/hxegon/.config/nvim",
"/home/hxegon/Notes"
]
```

`gitSearchDirs` will search for any directories containing a `.git` folder,
`extraDirs` will list those folders as is.

`lazytsm` will get you into a tmux session in that folder, either starting a new
one or attaching/switching to one with the same name

I use this through a tmux popup window so I can switch quickly from any program
through a key binding:

```
bind-key C-p popup -d "#{pane_current_path}" -h 80% -w 60% -E /home/hxegon/Code/hxegon/lazytsm/lazytsm
```

When you're in the program you can use `/` to filter, and you select what
session you want to open with `enter`

## V1 roadmap

- [x] Config file
- [x] extra directories
- [x] Sort by .git/objects modtime
- [x] Shorten $HOME parts of path to ~
- [x] Write a proper readme
- [x] nix build
- [x] Automatically filter when user starts typing non command chars
- [x] Make filter state style a bit prettier
- [ ] command line args
  - [ ] help option
- [ ] nix _Installation instructions_
- [ ] Instructions for adding popup window shortcut

- [ ] Store timestamp when a selecting a session, and use it to sort by most
      recently used
- [ ] Modes: default to search, but also have an ibuffer like mode where you can
      mark sessions for deletion
- [ ] List project tmux sessions, differentiate list items ( No duplicates of
      standard projects! )
- [ ] Parse ssh host config and have a way to start / reattach to a session on a
      remote host
- [ ] A way to save / restore sessions, or start sessions based on a config file
