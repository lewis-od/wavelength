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
		endRow:   0,
		action:   progress.Build,
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
	endRow   int
	action   progress.Action
}

func (d *ansiDisplay) Init(action progress.Action) {
	d.action = action
	d.numRows = 0
	d.endRow = 0
	d.statuses = make(map[string]*buildStatus)
}

func (d *ansiDisplay) Started(lambdaName string) {
	d.statuses[lambdaName] = &buildStatus{
		row:       d.numRows,
		completed: false,
	}
	message := fmt.Sprintf(d.action.InProgress, lambdaName)
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
		message = fmt.Sprintf(d.action.Success, lambdaName)
	} else {
		message = fmt.Sprintf(d.action.Error, lambdaName)
	}

	d.term.Print(message)
	d.term.MoveTo(d.endRow, 0)
}
