package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

type model struct {
	dirInfo         fileTracker
	cursor          int
	oneDirBack      dirBackInfo
	width           int
	height          int
	exitCmd         string
	err             error
	minRowToDisplay int
}

// Information about the files that should currently be displayed
type fileTracker struct {
	files []os.DirEntry
	path  string
}

// Saves the state of where the cursor was if the user hits 'b' to go back one directory
type dirBackInfo struct {
	cursor          int
	minRowToDisplay int
}

// For messages that contain errors. This is the idiomatic way of wrapping errors with Bubble Tea.

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func getInitialFiles() tea.Msg {
	dirEntries, err := os.ReadDir(".")
	if err != nil {
		return errMsg{err}
	}
	currDir, err := os.Getwd()
	if err != nil {
		return errMsg{err}
	}
	return fileTracker{files: dirEntries, path: currDir}
}

func (m model) Init() tea.Cmd {
	return getInitialFiles
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
