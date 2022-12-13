package auth

import (
	loginCmd "github.com/bitmyth/pdrive-cli/cli/cmd/auth/login"
	registerCmd "github.com/bitmyth/pdrive-cli/cli/cmd/auth/register"
	authTokenCmd "github.com/bitmyth/pdrive-cli/cli/cmd/auth/token"
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	"github.com/spf13/cobra"
)

func NewCmdAuth(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Authenticate with PDrive",
		Annotations: map[string]string{
			"IsCore": "true",
		},
	}

	cmd.AddCommand(loginCmd.NewCmdLogin(f, nil))
	cmd.AddCommand(authTokenCmd.NewCmdToken(f, nil))
	cmd.AddCommand(registerCmd.NewCmdRegister(f, nil))
	//cmd.AddCommand(authLogoutCmd.NewCmdLogout(f, nil))

	return cmd
}
