package download

import (
	"fmt"
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	"io"
	"net/url"
	"os"
	"strings"
	"time"
)

func Download(f *factory.Factory, id string) {
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
	cd := resp.Header.Get("Content-Disposition")
	//println("Content-Disposition ", cd)
	name := strings.Split(cd, "=")[1]
	fmt.Fprintln(f.IOStreams.Out, "Downloading ", infoColor(name))

	create, err := os.Create(name)
	if err != nil {
		return
	}
	_, _ = create.Write(respData)
	_ = create.Close()

	fmt.Fprintln(f.IOStreams.Out, "Size:", infoColor(fmt.Sprintf("%d", len(respData))))
}
