package list

import (
	"fmt"
	"time"

	"github.com/profclems/go-dotenv"

	"github.com/gosuri/uilive"
	"github.com/profclems/glab/pkg/tableprinter"
	"github.com/profclems/tfa/totp"
	"github.com/profclems/tfa/utils"
	"github.com/profclems/tfa/utils/color"
	"github.com/profclems/tfa/utils/iomanip"
	"github.com/spf13/cobra"
)

type Options struct {
	Name string
	Code string
	OTP  *totp.OTP
}

var (
	IO *iomanip.IO
)

func NewListCmd(iom *iomanip.IO, cfg *dotenv.DotEnv, runFunc func() error) *cobra.Command {
	var opts []*Options
	IO = iom

	var listCmd = &cobra.Command{
		Use:     "list",
		Short:   "List the available 2fa codes",
		Long:    ``,
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"view"},
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			// TODO: get from config file
			// tests := []string{"ayus", "wedd", "ABCDEFGHIJ234567", "234"}
			cfgCodes := cfg.Config
			for name, code := range cfgCodes {
				opt := &Options{
					Code: code.(string),
					Name: name,
				}
				opt.OTP, err = totp.New(opt.Name, opt.Code)
				if err != nil {
					return err
				}
				opts = append(opts, opt)
			}
			if runFunc != nil {
				return runFunc()
			}
			return getRun(opts)
		},
	}

	return listCmd
}

func getRun(opts []*Options) (err error) {
	writer := uilive.New()
	timer := 0
	timerMessage := ""
	writer.Out = IO.StdOut
	writer.RefreshInterval = time.Second
	table := tableprinter.NewTablePrinter()
	table.TerminalWidth = IO.TerminalWidth()
	// start listening for updates and render
	writer.Start()
	defer writer.Stop()
	if len(opts) > 0 {
		for {
			timer = int(opts[0].OTP.Timer)
			fmt.Fprintln(writer.Newline())
			fmt.Fprintln(writer, "Two Factor authentication codes")
			fmt.Fprintln(writer.Newline())
			for _, opt := range opts {
				err = opt.OTP.Refresh()
				if err != nil {
					return err
				}
				table.AddCell(opt.OTP.Name)
				table.AddCell(" ")
				table.AddCellf("%06d", opt.OTP.Password)
				table.EndRow()
				// fmt.Fprintln(writer.Newline(), opt.OTP.Name, " ", opt.OTP.Password)
			}
			fmt.Fprintln(writer, table.Render())
			table.Rows = nil // clear the table rows
			fmt.Fprintln(writer.Newline())
			timerMessage = utils.Pluralize(timer, "second")
			if timer > 15 {
				timerMessage = color.Green(timerMessage)
			} else if timer <= 15 && timer > 5 {
				timerMessage = color.Yellow(timerMessage)
			} else {
				timerMessage = color.Red(timerMessage)
			}
			fmt.Fprintf(writer.Newline(), "Timer: (%s remaining)\n", timerMessage)
			fmt.Fprint(writer.Newline(), "Press Ctrl-C to quit\n")
			time.Sleep(time.Second) // sleep for a second
		}
	} else {
		return nil
	}
}
