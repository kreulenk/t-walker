package display

// The number of rows to subtract from the window height to determine what should be displayed. This is due to other
// elements that are displayed on the screen.
const rowsToSubtract = 4

func ShouldPrintRow(currRow, windowHeight, minRowToDisplay int) bool {
	maxRowToDisplay := minRowToDisplay + windowHeight - rowsToSubtract
	return currRow >= minRowToDisplay && currRow <= maxRowToDisplay
}

func GetNewMinRowToDisplay(windowHeight, currentMinRowToDisplay, numColumns, currentCursor int) int {
	newCursorRow := currentCursor / numColumns
	currentMaxRow := currentMinRowToDisplay + windowHeight - rowsToSubtract

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
