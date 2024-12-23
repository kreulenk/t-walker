# T-Walker

T-Walker is a simple terminal user interface designed to be the quickest way to navigate directories and quickly edit files.

## Installation
Installing T-Walker is simple. Just clone the repository and run the install directive from the MakeFile.

The install directive assumes that you are using zsh as your shell. If you are using a different shell, you will need to
manually add the `t` command to your shell's sourcing file.

```bash
git clone https://github.com/kreulenk/t-walker.git
cd t-walker
make install
```

The final `make install` command may require root permissions.

## Usage
To use T-Walker, simply run the `t` command from the terminal. The interface will open and you can navigate directories
with the arrow keys and enter into directories files with the enter key.

![demo.gif](./docs/assets/demo.gif)


The following keybindings are available:
- `Up Arrow` - Move up one file.
- `Down Arrow` - Move down one file.
- `Left Arrow` - Move left one file.
- `Right Arrow` - Move right one file.
- `Enter` - Enter into a directory.
- `b` - Go back one directory.
- `e` - Edit the selected file. Defaults to vim but this can be overriden using the `EDITOR` environment variable.
- `c` - Change into the selected directory from your shell. This will cause T-walker to exit.
- `r` - Remove the selected file.
- `s` - Filter the files in the current directory.
- `q` - Quit the program.