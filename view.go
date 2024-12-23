package main

import (
	"fmt"
	"strings"
)

const (
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorGreen  = "\033[32m"
	colorReset  = "\033[0m"
)

const (
	colWidth         = 50
	rowsAlreadyTaken = 5 // The number of rows to subtract from the window height to determine what should be displayed
)

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
			if !(shouldPrintRow(currFile/numColumns, m.height, m.minRowToDisplay)) {
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

func shouldPrintRow(currRow, windowHeight, minRowToDisplay int) bool {
	maxRowToDisplay := minRowToDisplay + windowHeight - rowsAlreadyTaken
	return currRow >= minRowToDisplay && currRow <= maxRowToDisplay
}
