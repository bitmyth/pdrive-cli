package cd

import (
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	fileLsCmd "github.com/bitmyth/pdrive-cli/cli/cmd/file/ls"
)

func Cd(f *factory.Factory, name string) error {
	//cs := f.IOStreams.ColorScheme()
	//out := f.IOStreams.Out
	//infoColor := cs.Cyan
	switch {
	case name == "..":
		if fileLsCmd.Current.Parent != nil {
			println("parent id", fileLsCmd.Current.Parent.ID)
			fileLsCmd.Cd(*fileLsCmd.Current.Parent)
		}
	case name == "-":
		//fileLsCmd.Cd(fileLsCmd.Previous)
	default:
		for _, file := range fileLsCmd.Files {
			if file.Name == name {
				fileLsCmd.Cd(file)
				break
			}
		}

	}

	fileLsCmd.Ls(f)
	return nil
}
