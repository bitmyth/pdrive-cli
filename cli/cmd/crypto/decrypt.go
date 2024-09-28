package crypto

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/bitmyth/pdrive-cli/cli/cmd/crypto/rsa"
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	"github.com/spf13/cobra"
)

func NewCmdDecrypt(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "decrypt <data>",
		Short:   "decrypt",
		Long:    `decrypt data`,
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"dec"},

		SilenceErrors: true,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
			$ pd decrypt "0ba728f..."
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return rsa.Decrypt(f, args[0])
		},
	}

	return cmd
}
