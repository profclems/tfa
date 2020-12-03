package get

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/atotto/clipboard"
	"github.com/profclems/go-dotenv"
	"github.com/profclems/tfa/totp"
	"github.com/profclems/tfa/utils"
	"github.com/profclems/tfa/utils/iomanip"
	"github.com/spf13/cobra"
	"strings"
)

type Options struct {
	Name   string
	Copy   bool
	IO     *iomanip.IO
	Config *dotenv.DotEnv
}

func NewGetCmd(iom *iomanip.IO, cfg *dotenv.DotEnv, runFunc func() error) *cobra.Command {
	opts := &Options{
		IO:     iom,
		Config: cfg,
	}

	var getCmd = &cobra.Command{
		Use:     "get",
		Short:   "Get the OTP for an account",
		Long:    ``,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				opts.Name = strings.ToUpper(args[0])
			}

			if runFunc != nil {
				return runFunc()
			}

			if opts.Name == "" {
				var promptOptions []string
				for name := range opts.Config.Config {
					promptOptions = append(promptOptions, name)
				}
				if len(promptOptions) > 0 {
					prompt := &survey.Select{
						Message: "Select Account",
						Options: promptOptions,
					}
					err := survey.AskOne(prompt, &opts.Name, survey.WithValidator(survey.Required))
					if err != nil {
						return fmt.Errorf("could not prompt: %w", err)
					}
				}
			}
			return getRun(opts)
		},
	}

	getCmd.Flags().BoolVarP(&opts.Copy, "copy", "c", false, "Copy OTP to Clipboard")

	return getCmd
}

func getRun(opts *Options) error {
	code, isSet := opts.Config.Config[opts.Name]
	if !isSet {
		fmt.Fprintf(opts.IO.StdErr, "There's no account for %q in the config file\n", opts.Name)
		return nil
	}

	otp, err := totp.New(opts.Name, code.(string))
	if err != nil {
		return err
	}

	fmt.Fprintf(opts.IO.StdOut, "%06d (%s remaining)\n\n", otp.Password, utils.Pluralize(int(otp.Timer), "second"))

	if opts.Copy {
		err = clipboard.WriteAll(fmt.Sprintf("%06d", otp.Password))
		if err != nil {
			return fmt.Errorf("failed to copy to clipboard: %w", err)
		}
		fmt.Fprintln(opts.IO.StdErr, "Copied to Clipboard!")
	}
	return nil
}
