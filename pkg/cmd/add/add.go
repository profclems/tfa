package add

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/profclems/go-dotenv"
	"github.com/profclems/tfa/totp"
	"github.com/profclems/tfa/utils"
	"github.com/profclems/tfa/utils/iomanip"
	"github.com/spf13/cobra"
)

type Options struct {
	Name     string
	Code     string
	Override bool
	IO       *iomanip.IO
	Config   *dotenv.DotEnv
}

func NewAddCmd(iom *iomanip.IO, cfg *dotenv.DotEnv, runFunc func() error) *cobra.Command {
	opts := &Options{
		IO:     iom,
		Config: cfg,
	}

	var addCmd = &cobra.Command{
		Use:     "add",
		Short:   "Add a new account",
		Long:    ``,
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"new", "create"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				opts.Code = args[0]
			}

			if runFunc != nil {
				return runFunc()
			}

			if opts.Name == "" {
				err := survey.AskOne(&survey.Input{
					Message: "Name:",
				}, &opts.Name, survey.WithValidator(survey.Required))
				if err != nil {
					return fmt.Errorf("could not prompt: %w", err)
				}
			}

			if opts.Code == "" {
				err := survey.AskOne(&survey.Input{
					Message: "Secret Code:",
				}, &opts.Code, survey.WithValidator(survey.Required))
				if err != nil {
					return fmt.Errorf("could not prompt: %w", err)
				}
			}
			return addRun(opts)
		},
	}

	addCmd.Flags().StringVarP(&opts.Name, "name", "n", "", "Name of account to save with")
	addCmd.Flags().BoolVarP(&opts.Override, "yes", "y", false, "Force override existing")

	return addCmd
}

func addRun(opts *Options) error {
	otp, err := totp.New(opts.Name, opts.Code)
	if err != nil {
		return err
	}
	if value, isSet := opts.Config.Config[otp.Name]; isSet && value != nil {
		prompt := &survey.Confirm{
			Message: "Do you want to override?",
			Default: false,
		}
		err = survey.AskOne(prompt, &opts.Override)
		if err != nil {
			return fmt.Errorf("could not prompt: %w", err)
		}
	} else {
		opts.Override = true // there is no existing account so save
	}

	fmt.Fprintf(opts.IO.StdOut, "%06d (%s remaining)\n", otp.Password, utils.Pluralize(int(otp.Timer), "second"))

	if opts.Override {
		opts.Config.Set(otp.Name, otp.Code)
		if err = opts.Config.Save(); err != nil {
			return err
		}
	} else {
		fmt.Fprintln(opts.IO.StdErr, "Aborted")
	}
	return nil
}
