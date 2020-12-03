package root

import (
	"github.com/profclems/go-dotenv"
	"github.com/profclems/tfa/pkg/cmd/add"
	"github.com/profclems/tfa/pkg/cmd/get"
	"github.com/profclems/tfa/pkg/cmd/list"
	"github.com/profclems/tfa/pkg/cmd/version"
	"github.com/profclems/tfa/utils/iomanip"
	"github.com/spf13/cobra"
)

func NewRootCmd(iom *iomanip.IO, cfg *dotenv.DotEnv, buildVersion, buildDate string) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "tfa",
		Short: "TFA is a two-factor authentication CLI tool",
	}
	rootCmd.SetOut(iom.StdOut)
	rootCmd.SetErr(iom.StdErr)

	formattedVersion := version.Scheme(buildVersion, buildDate)
	rootCmd.SetVersionTemplate(formattedVersion)
	rootCmd.Version = formattedVersion

	rootCmd.AddCommand(version.NewCmdVersion(iom, buildVersion, buildDate))
	rootCmd.AddCommand(add.NewAddCmd(iom, cfg, nil))
	rootCmd.AddCommand(get.NewGetCmd(iom, cfg, nil))
	rootCmd.AddCommand(list.NewListCmd(iom, cfg, nil))

	return rootCmd
}
