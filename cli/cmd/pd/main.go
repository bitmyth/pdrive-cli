package main

import (
	"fmt"
	surveyCore "github.com/AlecAivazis/survey/v2/core"
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	"github.com/bitmyth/pdrive-cli/cli/cmd/root"
	"github.com/mgutz/ansi"
	"os"
)

type exitCode int

const (
	exitOK     exitCode = 0
	exitError  exitCode = 1
	exitCancel exitCode = 2
	exitAuth   exitCode = 4
)

func main() {
	code := mainRun()
	os.Exit(int(code))
}

func mainRun() exitCode {
	cmdFactory := factory.New()
	stderr := cmdFactory.IOStreams.ErrOut
	if !cmdFactory.IOStreams.ColorEnabled() {
		surveyCore.DisableColor = true
		ansi.DisableColors(true)
	} else {
		// override survey's poor choice of color
		surveyCore.TemplateFuncsWithColor["color"] = func(style string) string {
			switch style {
			case "white":
				if cmdFactory.IOStreams.ColorSupport256() {
					return fmt.Sprintf("\x1b[%d;5;%dm", 38, 242)
				}
				return ansi.ColorCode("default")
			default:
				return ansi.ColorCode(style)
			}
		}
	}

	_, err := cmdFactory.Config()
	if err != nil {
		fmt.Fprintf(stderr, "failed to read configuration:  %s\n", err)
		return exitError
	}

	rootCmd := root.NewCmdRoot(cmdFactory)
	if cmd, err := rootCmd.ExecuteC(); err != nil {
		_ = cmd
		fmt.Fprintf(stderr, "error: %s\n", err)
		return exitError
	}
	return exitOK

}
