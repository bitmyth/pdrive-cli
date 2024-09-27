package factory

import (
	"github.com/bitmyth/pdrive-cli/cli"
	"github.com/bitmyth/pdrive-cli/cli/config"
	"github.com/bitmyth/pdrive-cli/cli/iostreams"
	"github.com/bitmyth/pdrive-cli/cli/survey/prompter"
	"github.com/cli/go-gh/pkg/browser"
	"net/http"
	"os"
)

type Factory struct {
	IOStreams *iostreams.IOStreams
	Browser   browser.Browser

	HttpClient func() (*http.Client, error)
	Config     func() (config.Config, error)
	Prompter   prompter.Prompter
	HttpSchema string
	KeyFile    string
}

func New() *Factory {
	f := &Factory{
		Config: configFunc(),
	}
	f.IOStreams = ioStreams(f) // Depends on Config
	f.HttpClient = httpClientFunc(f)
	f.Prompter = newPrompter(f)  // Depends on Config and IOStreams
	f.HttpSchema = httpSchema(f) // Depends on Config

	return f
}

func configFunc() func() (config.Config, error) {
	var cachedConfig config.Config
	var configError error
	return func() (config.Config, error) {
		if cachedConfig != nil || configError != nil {
			return cachedConfig, configError
		}
		cachedConfig, configError = config.NewConfig()
		return cachedConfig, configError
	}
}

func ioStreams(f *Factory) *iostreams.IOStreams {
	io := iostreams.System()
	cfg, err := f.Config()
	if err != nil {
		return io
	}

	if _, ghPromptDisabled := os.LookupEnv("GH_PROMPT_DISABLED"); ghPromptDisabled {
		io.SetNeverPrompt(true)
	} else if prompt, _ := cfg.GetOrDefault("", "prompt"); prompt == "disabled" {
		io.SetNeverPrompt(true)
	}

	// Pager precedence
	// 1. GH_PAGER
	// 2. pager from config
	// 3. PAGER
	if ghPager, ghPagerExists := os.LookupEnv("GH_PAGER"); ghPagerExists {
		io.SetPager(ghPager)
	} else if pager, _ := cfg.Get("", "pager"); pager != "" {
		io.SetPager(pager)
	}

	return io
}

func httpClientFunc(f *Factory) func() (*http.Client, error) {
	return func() (*http.Client, error) {
		io := f.IOStreams
		cfg, err := f.Config()
		if err != nil {
			return nil, err
		}
		opts := cli.HTTPClientOptions{
			Config:      cfg,
			Log:         io.ErrOut,
			LogColorize: io.ColorEnabled(),
		}
		client, err := cli.NewHTTPClient(opts)
		if err != nil {
			return nil, err
		}
		//-client.Transport = api.ExtractHeader("X-GitHub-SSO", &ssoHeader)(client.Transport)
		return client, nil
	}
}

func newPrompter(f *Factory) prompter.Prompter {
	editor := ""
	io := f.IOStreams
	return prompter.New(editor, io.In, io.Out, io.ErrOut)
}

func httpSchema(f *Factory) string {
	cfg, err := f.Config()
	if err != nil {
		return ""
	}
	s, _ := cfg.GetOrDefault("", "env")
	if s == "dev" {
		return "http"
	}

	if s, _ = cfg.DefaultScheme(); s != "" {
		return s
	}

	return "https"
}
