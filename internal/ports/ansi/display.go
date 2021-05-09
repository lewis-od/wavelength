package ansi

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/progress"
)

func NewAnsiDisplay(term Terminal) progress.BuildDisplay {
	return &ansiDisplay{
		term:     term,
		statuses: make(map[string]*buildStatus),
		numRows:  0,
		endRow: 0,
	}
}

type buildStatus struct {
	row       int
	completed bool
}

type ansiDisplay struct {
	term     Terminal
	statuses map[string]*buildStatus
	numRows  int
	endRow int
}

func (d *ansiDisplay) Started(lambdaName string) {
	d.statuses[lambdaName] = &buildStatus{
		row:       d.numRows,
		completed: false,
	}
	message := fmt.Sprintf("🔨 Building %s...\n", lambdaName)
	d.term.Print(message)
	d.numRows++
}

func (d *ansiDisplay) Completed(lambdaName string, wasSuccessful bool) {
	status := d.statuses[lambdaName]
	status.completed = true

	if d.endRow == 0 {
		// Calling GetPosition moves the cursor down 1 line; hack to get around that
		position, err := d.term.GetPosition()
		if err != nil {
			panic(err)
		}
		d.endRow = position.Row - 1
		d.term.MoveTo(d.endRow, 0)
	}

	targetRow := d.endRow - d.numRows + status.row
	d.term.MoveTo(targetRow, 0)

	message := ""
	if wasSuccessful {
		message = fmt.Sprintf("✅  Building %s... done", lambdaName)
	} else {
		message = fmt.Sprintf("❌ Building %s... error", lambdaName)
	}

	d.term.Print(message)
	d.term.MoveTo(d.endRow, 0)
}
