package version

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/bitmyth/pdrive-cli/cli/build"
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	"github.com/spf13/cobra"
)

func NewCmdVersion(f *factory.Factory) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "version",
		Short:   "v",
		Long:    `show version`,
		Aliases: []string{"v"},

		SilenceErrors: true,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
			$ pd version
		`),
		Run: func(cmd *cobra.Command, args []string) {
			runVersion(f)
		},
	}

	cmd.PersistentFlags().Bool("help", false, "Show help for command")
	return cmd
}

func runVersion(f *factory.Factory) {
	io := f.IOStreams
	fmt.Fprintf(io.Out, "Version: %s", build.Version)
}
