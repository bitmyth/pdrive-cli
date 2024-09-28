package crypto

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/bitmyth/pdrive-cli/cli/cmd/crypto/rsa"
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	"github.com/spf13/cobra"
)

func NewCmdEncrypt(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "encrypt <data>",
		Args:    cobra.ExactArgs(1),
		Short:   "encrypt",
		Long:    `encrypt data`,
		Aliases: []string{"enc"},

		SilenceErrors: true,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
			$ pd encrypt "hello world"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			plain := args[len(args)-1]
			fmt.Printf("%q", plain)
			return rsa.Encrypt(f, plain)
		},
	}

	return cmd
}
