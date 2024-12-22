package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"t-walker/pkg/display"
)

const (
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorGreen  = "\033[32m"
	colorReset  = "\033[0m"
)

const colWidth = 50

type fileTracker struct {
	files []os.DirEntry
	path  string
}

// Saves the state of where the cursor was if the user hits 'b' to go back one directory
type dirBackInfo struct {
	cursor          int
	minRowToDisplay int
}

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

type errMsg struct{ err error }

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errMsg) Error() string { return e.err.Error() }

func (m model) Init() tea.Cmd {
	return getInitialFiles
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
		m.minRowToDisplay = display.GetNewMinRowToDisplay(m.height, m.minRowToDisplay, numColumns, m.cursor)
	}

	// If we happen to get any other messages, don't do anything.
	return m, nil
}

func (m model) View() string {
	if m.exitCmd != "" {
		return fmt.Sprintln("Executing command:", colorGreen, m.exitCmd, colorReset)
	}

	var s strings.Builder

	if len(m.dirInfo.files) > 0 {
		// Calculate the number of columns that can fit in the terminal width
		numColumns := m.width / colWidth
		if numColumns == 0 {
			numColumns = 1
		}

		for currFile := 0; currFile < len(m.dirInfo.files); currFile += numColumns {
			if !(display.ShouldPrintRow(currFile/numColumns, m.height, m.minRowToDisplay)) {
				continue
			}
			for currColumn := 0; currColumn < numColumns; currColumn++ {
				if currFile+currColumn >= len(m.dirInfo.files) {
					break
				}

				file := m.dirInfo.files[currFile+currColumn]
				cursorText := " " // no cursor
				if m.cursor == currFile+currColumn {
					cursorText = ">" // cursor!
				}

				var colorToUse string
				if file.IsDir() {
					colorToUse = colorBlue
				} else {
					colorToUse = colorPurple
				}

				fileInfo, err := file.Info()

				var fileName string
				if len(file.Name()) > 29 {
					fileName = file.Name()[0:28] + "…"
				} else {
					fileName = file.Name()
				}
				if file.IsDir() {
					fileName = fileName + "/"
				}

				if err != nil {
					s.WriteString(fmt.Sprintf("%s%s %-30s %s%s", cursorText, colorToUse, fileName, "error reading permission bits", colorReset))
				} else {
					s.WriteString(fmt.Sprintf("%s%s %-30s %s%s", cursorText, colorToUse, fileName, fileInfo.Mode(), colorReset))
				}

				if currColumn < numColumns-1 {
					s.WriteString(" | ")
				}
			}
			s.WriteString("\n")
		}
	} else {
		s.WriteString("No items to display in this directory...\n")
	}

	if m.err != nil {
		s.WriteString(fmt.Sprintf("Error: %v\n", m.err))
	} else {
		s.WriteString("\n")
	}
	s.WriteString(fmt.Sprintf("Current Directory: %s\n", m.dirInfo.path))

	return s.String()
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
