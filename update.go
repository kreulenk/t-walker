package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (m model) keyPressNormalMode(msg tea.KeyMsg) (model, tea.Cmd) {
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
		if m.cursor < len(m.dirInfo.searchFilteredFiles)-1 {
			m.cursor = m.cursor + numColumns
		}

	case "right":
		if m.cursor < len(m.dirInfo.searchFilteredFiles)-1 {
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
			m.searchStr = ""
		}

	case "c":
		m.exitCmd = fmt.Sprintf("cd %s", m.dirInfo.path)
		return m, tea.Quit
	case "e":
		if m.dirInfo.searchFilteredFiles[m.cursor].IsDir() {
			m.err = fmt.Errorf("cannot open file %s", m.dirInfo.searchFilteredFiles[m.cursor].Name())
			return m, nil
		}
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}
		openEditorCmd := exec.Command(editor, filepath.Join(m.dirInfo.path, m.dirInfo.searchFilteredFiles[m.cursor].Name()))
		openEditorCmd.Stdout = os.Stdout
		openEditorCmd.Stdin = os.Stdin
		openEditorCmd.Stderr = os.Stderr
		err := openEditorCmd.Run()
		if err != nil {
			m.err = fmt.Errorf("error opening editor: %v", err)
		}
	case "s":
		m.mode = search
	case "r":
		m.mode = removeFileConfirm
	case "enter":
		if m.dirInfo.searchFilteredFiles[m.cursor].IsDir() {
			m.dirInfo, m.err = m.walkIntoDir()
			if m.err == nil {
				m.oneDirBack.cursor = m.cursor
				m.oneDirBack.minRowToDisplay = m.minRowToDisplay
				m.cursor = 0
				m.searchStr = ""
			}
		}
	}
	return m, nil
}

func (m model) keyPressSearchMode(msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = normal
	case "ctrl+c":
		return m, tea.Quit
	case "enter":
		m.mode = normal
	case "backspace":
		if len(m.searchStr) > 0 {
			m.searchStr = m.searchStr[:len(m.searchStr)-1]
		}
		m.cursor = 0
	default:
		m.searchStr = m.searchStr + msg.String()
		m.cursor = 0
	}
	m.dirInfo.searchFilteredFiles = filterFilesBySearch(m.dirInfo.files, m.searchStr)
	return m, nil
}

func (m model) keyPressRemoveFileConfirmMode(msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "y":
		err := os.Remove(filepath.Join(m.dirInfo.path, m.dirInfo.searchFilteredFiles[m.cursor].Name()))
		if err != nil {
			m.err = errMsg{err}
		}
		dirEntries, err := os.ReadDir(m.dirInfo.path)
		if err != nil {
			m.err = errMsg{err}
			return m, nil
		}
		m.dirInfo.files = dirEntries
		m.dirInfo.searchFilteredFiles = dirEntries
		m.mode = normal
	case "ctrl+c":
		return m, tea.Quit
	default:
		m.mode = normal
	}
	return m, nil
}

func filterFilesBySearch(allFiles []os.DirEntry, searchTerm string) []os.DirEntry {
	if searchTerm == "" {
		return allFiles
	}
	var filteredFiles []os.DirEntry
	for _, file := range allFiles {
		if strings.Contains(file.Name(), searchTerm) {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles
}

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
		return m, nil

	case tea.KeyMsg:
		// After cursor updates, make sure the cursor is within bounds
		switch m.mode {
		case normal:
			newModel, cmd := m.keyPressNormalMode(msg)
			if cmd != nil {
				return m, cmd
			}
			newModel.minRowToDisplay = getNewMinRowToDisplay(m.height, m.width, m.minRowToDisplay, m.cursor)
			return newModel, nil
		case search:
			newModel, cmd := m.keyPressSearchMode(msg)
			if cmd != nil {
				return m, cmd
			}
			return newModel, nil
		case removeFileConfirm:
			newModel, cmd := m.keyPressRemoveFileConfirmMode(msg)
			if cmd != nil {
				return m, cmd
			}
			return newModel, nil
		default:
			panic("unhandled default case")
		}
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

	return fileTracker{files: dirEntries, searchFilteredFiles: dirEntries, path: oneBackFilePath}, nil
}

func (m model) walkIntoDir() (fileTracker, error) {
	newPath := filepath.Join(m.dirInfo.path, m.dirInfo.searchFilteredFiles[m.cursor].Name())
	dirEntries, err := os.ReadDir(newPath)
	if err != nil {
		return m.dirInfo, errMsg{err}
	}
	return fileTracker{files: dirEntries, searchFilteredFiles: dirEntries, path: newPath}, nil
}

func getNewMinRowToDisplay(windowHeight, windowWidth, currentMinRowToDisplay, currentCursor int) int {
	numColumns := windowWidth / colWidth
	newCursorRow := currentCursor / numColumns
	currentMaxRow := currentMinRowToDisplay + windowHeight - rowsAlwaysTaken

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
