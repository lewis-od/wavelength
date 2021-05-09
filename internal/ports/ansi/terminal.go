package ansi

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type CursorPosition struct {
	Row, Col int
}

type Terminal interface {
	GetPosition() (CursorPosition, error)
	MoveTo(row, col int)
	Print(text string)
}

func NewAnsiTerminal() Terminal {
	return &ansiTerminal{}
}

type ansiTerminal struct {}

var outputRegexp = regexp.MustCompile(`(?P<row>\d+);(?P<col>\d+)`)

func (t *ansiTerminal) GetPosition() (CursorPosition, error) {
	toggleRawMode(on)

	// Same as $ echo -e "\033[6n"
	cmd := exec.Command("echo", fmt.Sprintf("%c[6n", 27))
	randomBytes := &bytes.Buffer{}
	cmd.Stdout = randomBytes
	_ = cmd.Start()

	// Capture output from echo command
	reader := bufio.NewReader(os.Stdin)
	cmd.Wait()

	// Trigger input
	fmt.Print(randomBytes)
	rawOutput, _ := reader.ReadSlice('R')
	output := string(rawOutput)

	toggleRawMode(off)

	if strings.Contains(output, ";") {
		return parseOutput(output)
	}

	return CursorPosition{}, fmt.Errorf("unable to determine cursor location")
}

func parseOutput(output string) (CursorPosition, error)  {
	matches := outputRegexp.FindStringSubmatch(output)
	rowIndex, colIndex := outputRegexp.SubexpIndex("row"), outputRegexp.SubexpIndex("col")
	rowStr, colStr := matches[rowIndex], matches[colIndex]

	row, err := strconv.Atoi(rowStr)
	if err != nil {
		return CursorPosition{}, err
	}

	col, err := strconv.Atoi(colStr)
	if err != nil {
		return CursorPosition{}, err
	}

	return CursorPosition{
		Row: row,
		Col: col,
	}, nil
}

type state bool
const on state = true
const off state = false

func toggleRawMode(desiredState state) {
	argument := "raw"
	if desiredState == off {
		argument = "-raw"
	}
	rawMode := exec.Command("/bin/stty", argument)
	rawMode.Stdin = os.Stdin
	_ = rawMode.Run()
	rawMode.Wait()
}

func (t *ansiTerminal) MoveTo(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

func (t *ansiTerminal) Print(text string) {
	fmt.Print(text)
}

func StdoutIsTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
