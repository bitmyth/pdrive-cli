package cat

import (
	"encoding/hex"
	"fmt"
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	"github.com/bitmyth/pdrive-cli/cli/secret"
	"io"
	"net/url"
	"time"
)

func Cat(f *factory.Factory, id string) {
	cs := f.IOStreams.ColorScheme()

	defer func(t time.Time) {
		fmt.Fprintf(f.IOStreams.Out, "Time %s\n", cs.Cyan(time.Now().Sub(t).String()))
	}(time.Now())

	client, _ := f.HttpClient()

	cfg, _ := f.Config()
	hostname, _ := cfg.DefaultHost()

	values := url.Values{}
	values.Add("id", id)
	val := values.Encode()
	url := fmt.Sprintf("%s://%s/v1/files/download?"+val, f.HttpSchema, hostname)

	resp, err := client.Get(url)
	if err != nil {
		return
	}

	respData, _ := io.ReadAll(resp.Body)

	infoColor := cs.Cyan

	// save file
	fmt.Fprintln(f.IOStreams.Out, "Size:", infoColor(fmt.Sprintf("%d", len(respData))))
	fmt.Fprintln(f.IOStreams.Out, cs.Green(string(respData)))

	_, err = hex.DecodeString(string(respData))
	if err == nil {
		decrypt, er := secret.NewRSA(f.KeyFile).Decrypt(string(respData))
		if er == nil {
			fmt.Fprintln(f.IOStreams.Out, cs.Cyan("Decrypted"))
			fmt.Fprintln(f.IOStreams.Out, cs.Green(string(decrypt)))
		}
	}
}
