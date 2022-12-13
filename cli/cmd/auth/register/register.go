package register

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
	Scopes     []string
	Name       string
	Email      string
	HttpSchema string
}

func NewCmdRegister(f *factory.Factory, runF func(*Options) error) *cobra.Command {
	opts := &Options{
		IO:         f.IOStreams,
		Config:     f.Config,
		HttpClient: f.HttpClient,
		Prompter:   f.Prompter,
		HttpSchema: f.HttpSchema,
	}

	cmd := &cobra.Command{
		Use:   "register",
		Args:  cobra.ExactArgs(0),
		Short: "Register with a PDrive host",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := opts.Config()
			opts.Hostname, _ = cfg.DefaultHost()

			if opts.IO.CanPrompt() {
				opts.Interactive = true
			}
			return registerRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "User name")
	cmd.Flags().StringVarP(&opts.Email, "email", "e", "", "User email")

	return cmd
}

func registerRun(opts *Options) error {
	cfg, _ := opts.Config()

	hostname := opts.Hostname
	if hostname == "" {
		var err error
		hostname, err = opts.Prompter.InputHostName()
		if err != nil {
			return err
		}
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

		hostname, err = opts.Prompter.InputHostName()
		if err != nil {
			return err
		}
		opts.Hostname = hostname
	}

	name := opts.Name
	var err error
	if name == "" {
		name, err = opts.Prompter.InputUserName()
		if err != nil {
			return err
		}
	}

	email := opts.Email
	if email == "" {
		email, err = opts.Prompter.InputEmail()
		if err != nil {
			return err
		}
	}
	var password string
	password, err = opts.Prompter.Password("Input password")

	user, err := register(opts, name, email, password)
	if err != nil {
		return err
	}

	authToken := user.Tokens[0].Content
	cfg.Set(hostname, "oauth_token", authToken)
	cfg.Set("", "default_host", hostname)

	cs := opts.IO.ColorScheme()
	fmt.Fprintf(opts.IO.Out, "%s register success! token: %s\n", cs.SuccessIcon(), cs.Cyan(authToken))
	return cfg.Write()
}

type User struct {
	ID       uint    `gorm:"primarykey"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Tokens   []Token `json:"tokens" gorm:"-"`
}

type Token struct {
	Content string `json:"content"`
}

func register(opts *Options, name string, email string, password string) (user User, err error) {
	cfg, _ := opts.Config()
	cs := opts.IO.ColorScheme()
	warnColor := cs.Yellow

	client, _ := opts.HttpClient()

	hostname := opts.Hostname
	url := fmt.Sprintf("%s://%s/v1/register", opts.HttpSchema, hostname)

	marshal, err := json.Marshal(User{Name: name, Email: email, Password: password})
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(marshal))
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	respData, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(respData, &user)
	if err != nil {
		fmt.Fprintf(opts.IO.ErrOut, "error %s\n", warnColor(string(respData)))
		return
	}

	cfg.Set("", "default_host", hostname)
	cfg.Write()
	return
}
