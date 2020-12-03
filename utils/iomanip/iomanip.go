package iomanip

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/profclems/tfa/utils/color"
	"github.com/profclems/tfa/utils/terminal"
)

type IO struct {
	In     io.ReadCloser
	StdOut io.Writer
	StdErr io.Writer

	IsaTTY   bool // stdout is a tty
	IsErrTTY bool // stderr is a tty
	IsInTTY  bool // stdin is a tty
}

func InitIO() *IO {
	stdoutIsTTY := terminal.IsTerminal(os.Stdout)
	stderrIsTTY := terminal.IsTerminal(os.Stderr)

	ioStream := &IO{
		In:       os.Stdin,
		StdOut:   color.NewColorable(os.Stdout),
		StdErr:   color.NewColorable(os.Stderr),
		IsaTTY:   stdoutIsTTY,
		IsErrTTY: stderrIsTTY,
	}

	if stdin, ok := ioStream.In.(*os.File); ok {
		ioStream.IsInTTY = terminal.IsTerminal(stdin)
	}

	color.IsColorEnabled = color.IsEnabled() && stdoutIsTTY && stderrIsTTY

	return ioStream
}

func (s *IO) TerminalWidth() int {
	return terminal.Width(s.StdOut)
}

func TestIO() (*IO, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	return &IO{
		In:     ioutil.NopCloser(in),
		StdOut: out,
		StdErr: errOut,
	}, in, out, errOut
}
