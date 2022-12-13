package login

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	"github.com/bitmyth/pdrive-cli/cli/config"
	"github.com/bitmyth/pdrive-cli/cli/iostreams"
	"github.com/bitmyth/pdrive-cli/cli/survey/prompter"
	"github.com/spf13/cobra"
	"io"
	"net/http"
)

type Options struct {
	IO         *iostreams.IOStreams
	Config     func() (config.Config, error)
	HttpClient func() (*http.Client, error)
	Prompter   prompter.Prompter

	MainExecutable string

	Interactive bool

	Hostname   string
	HttpSchema string
	Scopes     []string
	Token      string
	Web        bool
}

func NewCmdLogin(f *factory.Factory, runF func(*Options) error) *cobra.Command {
	opts := &Options{
		IO:         f.IOStreams,
		Config:     f.Config,
		HttpClient: f.HttpClient,
		Prompter:   f.Prompter,
		HttpSchema: f.HttpSchema,
	}

	var tokenStdin bool

	cmd := &cobra.Command{
		Use:   "login",
		Args:  cobra.ExactArgs(0),
		Short: "Authenticate with a PDrive host",
		RunE: func(cmd *cobra.Command, args []string) error {
			//if opts.Hostname == "" {
			//	cfg, _ := opts.Config()
			//	opts.Hostname, _ = cfg.DefaultHost()
			//}

			if opts.IO.CanPrompt() && opts.Token == "" {
				opts.Interactive = true
			}
			return loginRun(opts)
		},
	}

	cmd.Flags().BoolVar(&tokenStdin, "with-token", false, "Read token from standard input")
	cmd.Flags().BoolVarP(&opts.Web, "web", "w", false, "Open a browser to authenticate")

	return cmd
}

func loginRun(opts *Options) error {
	cfg, _ := opts.Config()

	hostname := opts.Hostname
	if hostname == "" {
		var err error
		hostname, err = opts.Prompter.InputHostName()
		if err != nil {
			return err
		}
		cfg.Set("", "default_host", hostname)
		opts.Hostname = hostname
	}

	existingToken, _ := cfg.AuthToken(hostname)
	if existingToken != "" && opts.Interactive {
		keepGoing, err := opts.Prompter.Confirm(fmt.Sprintf("You're already logged into %s. Do you want to re-authenticate?", hostname), false)
		if err != nil {
			return err
		}
		if !keepGoing {
			return nil
		}

	}

	options := []string{"Login by username password", "Paste an authentication token"}
	var err error
	var authMode int
	authMode, err = opts.Prompter.Select(
		"How would you like to authenticate PDrive CLI?",
		options[0],
		options)
	if err != nil {
		return err
	}
	var authToken string
	if authMode == 0 {
		var username, password string
		username, err = opts.Prompter.InputUserName()
		if err != nil {
			return err
		}
		password, err = opts.Prompter.Password("Input your password")
		if err != nil {
			return err
		}
		client, _ := opts.HttpClient()
		url := fmt.Sprintf("%s://%s/v1/login", opts.HttpSchema, hostname)

		var c = map[string]string{"name": username, "password": password}
		marshal, _ := json.Marshal(c)
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(marshal))
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		var tokenResp map[string]string
		respBody, _ := io.ReadAll(resp.Body)
		println(string(respBody))
		json.Unmarshal(respBody, &tokenResp)
		fmt.Fprintf(opts.IO.Out, tokenResp["token"])
		authToken = tokenResp["token"]

	} else {
		authToken, err = opts.Prompter.AuthToken()
		if err != nil {
			return err
		}
	}
	cfg.Set(hostname, "oauth_token", authToken)
	return cfg.Write()
}
