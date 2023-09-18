package root

import (
	"github.com/MakeNowJust/heredoc"
	authCmd "github.com/bitmyth/pdrive-cli/cli/cmd/auth"
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	fileCmd "github.com/bitmyth/pdrive-cli/cli/cmd/file"
	catCmd "github.com/bitmyth/pdrive-cli/cli/cmd/file/cat"
	"github.com/bitmyth/pdrive-cli/cli/cmd/file/cd"
	fileDeleteCmd "github.com/bitmyth/pdrive-cli/cli/cmd/file/delete"
	fileDownloadCmd "github.com/bitmyth/pdrive-cli/cli/cmd/file/download"
	fileLsCmd "github.com/bitmyth/pdrive-cli/cli/cmd/file/ls"
	"github.com/bitmyth/pdrive-cli/cli/cmd/file/mkdir"
	fileShareCmd "github.com/bitmyth/pdrive-cli/cli/cmd/file/share"
	"github.com/bitmyth/pdrive-cli/cli/cmd/file/upload"
	"github.com/bitmyth/pdrive-cli/cli/cmd/version"
	"github.com/spf13/cobra"
	"strings"
)

func NewCmdRoot(f *factory.Factory) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "pd <command> <subcommand> [flags]",
		Short: "PDrive CLI",
		Long:  `Work seamlessly with PDrive from the command line.`,

		SilenceErrors: false,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
			$ pd auth login
			$ pd auth token
			$ pd file upload
		`),
		Annotations: map[string]string{
			"help:feedback": heredoc.Doc(`
				Open an issue using 'gh issue create -R github.com/bitmyth/go-pdrive'
			`),
		},
		Run: func(cmd *cobra.Command, args []string) {
			runRoot(f)
		},
	}

	cmd.PersistentFlags().Bool("help", false, "Show help for command")

	// Child commands
	cmd.AddCommand(fileCmd.NewCmdFile(f))
	cmd.AddCommand(authCmd.NewCmdAuth(f))
	cmd.AddCommand(version.NewCmdVersion(f))

	return cmd
}

func runRoot(f *factory.Factory) {
	for {
		input, err := f.Prompter.Input("", "")
		if err != nil {
			return
		}
		op := strings.Split(input, " ")
		switch op[0] {
		case "cat":
			catCmd.Cat(f, op[1])
		case "d":
			fileDownloadCmd.Download(f, op[1])
		case "u":
			upload.RunUploadFile(f, op[1], fileLsCmd.Dir)
		case "exit":
			return
		case "ls", "l":
			fileLsCmd.Ls(f)
		case "n":
			fileLsCmd.Page++
			fileLsCmd.Ls(f)
		case "p":
			if fileLsCmd.Page > 1 {
				fileLsCmd.Page--
			}
			fileLsCmd.Ls(f)
		case "rm":
			fileDeleteCmd.Delete(f, op[1])
		case "mkdir":
			dir := fileLsCmd.Dir
			mkdir.Mkdir(f, op[1], dir)
		case "cd":
			cd.Cd(f, op[1])
		case "url":
			fileShareCmd.URL(f, op[1])
		}
	}
}
