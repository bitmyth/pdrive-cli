package token

import (
	"fmt"
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	"github.com/bitmyth/pdrive-cli/cli/config"
	"github.com/bitmyth/pdrive-cli/cli/iostreams"
	"github.com/spf13/cobra"
)

type TokenOptions struct {
	IO     *iostreams.IOStreams
	Config func() (config.Config, error)

	Hostname string
}

func NewCmdToken(f *factory.Factory, runF func(*TokenOptions) error) *cobra.Command {
	opts := &TokenOptions{
		IO:     f.IOStreams,
		Config: f.Config,
	}

	cmd := &cobra.Command{
		Use:   "token",
		Short: "Print the auth token pd is configured to use",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}

			return tokenRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Hostname, "hostname", "h", "", "The hostname of the PDrive instance authenticated with")

	return cmd
}

func tokenRun(opts *TokenOptions) error {
	hostname := opts.Hostname

	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	if hostname == "" {
		hostname, _ = cfg.DefaultHost()
	}
	key := "oauth_token"
	val, err := cfg.GetOrDefault(hostname, key)
	if err != nil {
		return fmt.Errorf("no oauth token")
	}

	if val != "" {
		fmt.Fprintf(opts.IO.Out, "%s\n", val)
	}
	return nil
}
