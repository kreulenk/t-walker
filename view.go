package main

import (
	"fmt"
	"strings"
)

const (
	colorBlue      = "\033[34m"
	colorPurple    = "\033[35m"
	colorGreen     = "\033[32m"
	colorLightBlue = "\033[36m"
	colorReset     = "\033[0m"
)

const (
	colWidth        = 50
	rowsAlwaysTaken = 5 // The number of rows to always subtract from the window height to determine what should be displayed
)

func (m model) View() string {
	var s strings.Builder

	if len(m.dirInfo.searchFilteredFiles) > 0 {
		// Calculate the number of columns that can fit in the terminal width
		numColumns := m.width / colWidth
		if numColumns == 0 {
			numColumns = 1
		}

		for rowStartIndex := 0; rowStartIndex < len(m.dirInfo.searchFilteredFiles); rowStartIndex += numColumns { // Iterate over files one row at a time
			if !(shouldPrintRow(rowStartIndex/numColumns, m.height, m.minRowToDisplay)) {
				continue
			}
			for colIndex := 0; colIndex < numColumns; colIndex++ { // Iterate over every file in a single row
				if rowStartIndex+colIndex >= len(m.dirInfo.searchFilteredFiles) {
					break
				}

				file := m.dirInfo.searchFilteredFiles[rowStartIndex+colIndex]
				cursorText := " "
				if m.cursor == rowStartIndex+colIndex {
					cursorText = ">"
				}

				var colorToUse string
				if file.IsDir() {
					colorToUse = colorBlue
				} else if isSymLinkDir(file, m.dirInfo.path) {
					colorToUse = colorPurple
				} else {
					colorToUse = colorLightBlue
				}

				fileInfo, err := file.Info()

				var fileName string
				if len(file.Name()) > 29 {
					fileName = file.Name()[0:28] + "…"
				} else {
					fileName = file.Name()
				}
				if file.IsDir() || isSymLinkDir(file, m.dirInfo.path) {
					fileName = fileName + "/"
				}

				if err != nil {
					s.WriteString(fmt.Sprintf("%s%s %-30s %s%s", cursorText, colorToUse, fileName, "error reading permission bits", colorReset))
				} else {
					s.WriteString(fmt.Sprintf("%s%s %-30s %s%s", cursorText, colorToUse, fileName, fileInfo.Mode(), colorReset))
				}

				if colIndex < numColumns-1 {
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

	if m.mode == search {
		s.WriteString(fmt.Sprintf("Current Directory: %s | Press 'esc' to exit search mode\n", m.dirInfo.path))
	} else {
		s.WriteString(fmt.Sprintf("Current Directory: %s\n", m.dirInfo.path))
	}

	if m.searchStr != "" && (m.mode == normal || m.mode == search) {
		s.WriteString(fmt.Sprintf("Search: %s\n", m.searchStr))
	} else if m.mode == removeFileConfirm {
		s.WriteString("Remove file? (y/n): ")
	}

	return s.String()
}

func shouldPrintRow(currRow, windowHeight, minRowToDisplay int) bool {
	maxRowToDisplay := minRowToDisplay + windowHeight - rowsAlwaysTaken
	return currRow >= minRowToDisplay && currRow <= maxRowToDisplay
}
