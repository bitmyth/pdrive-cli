package file

import (
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	cmdCreate "github.com/bitmyth/pdrive-cli/cli/cmd/file/create"
	cmdUpload "github.com/bitmyth/pdrive-cli/cli/cmd/file/upload"
	"github.com/spf13/cobra"
)

func NewCmdFile(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file <command>",
		Short: "Manage Files",
		Long:  "Manage Files.",
	}

	cmd.AddCommand(cmdUpload.NewCmdUpload(f))
	cmd.AddCommand(cmdCreate.NewCmdCreate(f))

	return cmd
}
