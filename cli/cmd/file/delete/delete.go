package delete

import (
	"fmt"
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	"io"
	"net/http"
)

func Delete(f *factory.Factory, id string) {
	client, _ := f.HttpClient()

	cfg, _ := f.Config()
	hostname, _ := cfg.DefaultHost()

	url := fmt.Sprintf("%s://%s/v1/files/%s", f.HttpSchema, hostname, id)

	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(f.IOStreams.Out, err.Error())
		return
	}

	respData, _ := io.ReadAll(resp.Body)

	cs := f.IOStreams.ColorScheme()
	infoColor := cs.Cyan

	fmt.Fprintln(f.IOStreams.Out, infoColor("deleted"), string(respData))
}
