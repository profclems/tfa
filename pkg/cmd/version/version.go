package version

import (
	"fmt"
	"strings"

	"github.com/profclems/tfa/utils/iomanip"
	"github.com/spf13/cobra"
)

func NewCmdVersion(iom *iomanip.IO, buildVersion, buildDate string) *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "show 2fa version information",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprint(iom.StdOut, Scheme(buildVersion, buildDate))
			return nil
		},
	}
	return versionCmd
}

func Scheme(version, buildDate string) string {
	version = strings.TrimPrefix(version, "v")

	if buildDate != "" {
		version = fmt.Sprintf("%s (%s)", version, buildDate)
	}

	return fmt.Sprintf("2fa version %s\n", version)
}
