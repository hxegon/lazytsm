# A TUI for navigating tmux sessions.

`lazytsm` treats a set of preconfigured directories as a list of (potential or
active) tmux sessions. Select a directory containing git repos, and you can
fuzzy search them. Press enter to either create a new tmux session with that
repo as the working directory, or reattach to an existing session associated
with it.

![](https://github.com/hxegon/lazytsm/blob/main/show.gif)

Extremely pre-alpha, but it sets the groundwork for a lot of functionality I've
wanted to add to the old `tms` script
https://github.com/ThePrimeagen/.dotfiles/blob/master/bin/.local/scripts/tmux-sessionizer

This is written by a Go newbie so any feedback / code review would be much
appreciated.

## Installation

Check the justfile for build/run instructions. You can also install this as a
nix package.

For dev dependencies there's a `flake.nix` and a `.envrc` with `use flake` if
you're into that sort of thing

### TMux popup window shortcut

The suggested way to use this inside tmux is to use a shortcut for opening
`lazytsm` inside a tmux popup window:

```tmux
# in ~/.tmux.conf or ~/config/tmux/tmux.conf
bind-key C-g popup -d "#{pane_current_path}" -h 90% -w 100% -E lazygit
```

You can then use your tmux leader, then control-g to open it.

## Configuration

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
`extraDirs` will add those folders as is.

## Usage

You can use Ctrl+n or Ctrl+p to navigate up and down the list of candidates, or
just start typing to search for something by name.
