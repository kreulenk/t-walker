package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.err = nil // Reset the error message if the user tries a new action
	switch msg := msg.(type) {
	case fileTracker:
		m.dirInfo = msg
		return m, nil

	case errMsg:
		// There was an error. Note it in the model. And tell the runtime
		// we're done and want to quit.
		m.err = msg
		return m, nil

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.KeyMsg:
		numColumns := m.width / colWidth
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

			// The "up" and "k" keys move the cursor up
		case "up":
			if m.cursor > numColumns-1 {
				m.cursor = m.cursor - numColumns
			}

		// The "down" and "j" keys move the cursor down
		case "down":
			if m.cursor < len(m.dirInfo.files)-1 {
				m.cursor = m.cursor + numColumns
			}

		case "right":
			if m.cursor < len(m.dirInfo.files)-1 {
				m.cursor++
			}
		case "left":
			if m.cursor > 0 {
				m.cursor--
			}
		case "b":
			if m.dirInfo.path != "/" {
				m.dirInfo, m.err = m.getCurrDirOneBack()
				m.cursor = m.oneDirBack.cursor
				m.minRowToDisplay = m.oneDirBack.minRowToDisplay
				m.oneDirBack.cursor = 0
			}

		case "c":
			m.exitCmd = fmt.Sprintf("cd %s", m.dirInfo.path)
			return m, tea.Quit
		case "e":
			if m.dirInfo.files[m.cursor].IsDir() {
				m.err = fmt.Errorf("cannot open file %s", m.dirInfo.files[m.cursor].Name())
				return m, nil
			}
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vi"
			}
			openEditorCmd := exec.Command(editor, filepath.Join(m.dirInfo.path, m.dirInfo.files[m.cursor].Name()))
			openEditorCmd.Stdout = os.Stdout
			openEditorCmd.Stdin = os.Stdin
			openEditorCmd.Stderr = os.Stderr
			err := openEditorCmd.Run()
			if err != nil {
				m.err = fmt.Errorf("error opening editor: %v", err)
			}
		case "enter":
			if m.dirInfo.files[m.cursor].IsDir() {
				m.dirInfo, m.err = m.walkIntoDir()
				if m.err == nil {
					m.oneDirBack.cursor = m.cursor
					m.oneDirBack.minRowToDisplay = m.minRowToDisplay
					m.cursor = 0
				}
			}
		}
		// After cursor updates, make sure the cursor is within bounds
		m.minRowToDisplay = getNewMinRowToDisplay(m.height, m.minRowToDisplay, numColumns, m.cursor)
	}

	// If we happen to get any other messages, don't do anything.
	return m, nil
}

func (m model) getCurrDirOneBack() (fileTracker, error) {
	splitPath := strings.Split(m.dirInfo.path, "/")
	oneBackFilePath := "/" + filepath.Join(splitPath[0:len(splitPath)-1]...)
	dirEntries, err := os.ReadDir(oneBackFilePath)
	if err != nil {
		return m.dirInfo, errMsg{err}
	}

	return fileTracker{files: dirEntries, path: oneBackFilePath}, nil
}

func (m model) walkIntoDir() (fileTracker, error) {
	newPath := filepath.Join(m.dirInfo.path, m.dirInfo.files[m.cursor].Name())
	dirEntries, err := os.ReadDir(newPath)
	if err != nil {
		return m.dirInfo, errMsg{err}
	}
	return fileTracker{files: dirEntries, path: newPath}, nil
}

func getNewMinRowToDisplay(windowHeight, currentMinRowToDisplay, numColumns, currentCursor int) int {
	newCursorRow := currentCursor / numColumns
	currentMaxRow := currentMinRowToDisplay + windowHeight - rowsAlreadyTaken

	var newMinRowToDisplay int
	if newCursorRow < currentMinRowToDisplay {
		newMinRowToDisplay = currentMinRowToDisplay - 1
	} else if newCursorRow > currentMaxRow {
		newMinRowToDisplay = currentMinRowToDisplay + 1
	} else {
		newMinRowToDisplay = currentMinRowToDisplay
	}

	return newMinRowToDisplay
}
