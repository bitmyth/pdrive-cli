package share

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	"io"
	"net/http"
	"strconv"
)

func URL(f *factory.Factory, id string) {
	i, _ := strconv.Atoi(id)
	createShare(f, i)
	//cfg, _ := f.Config()
	//hostname, _ := cfg.DefaultHost()
	//get, _ := cfg.Get(hostname, "oauth_token")
	//
	//url := fmt.Sprintf("%s://%s/v1/files/%s?raw=1&token=%s", f.HttpSchema, hostname, id, get)
	//
	//cs := f.IOStreams.ColorScheme()
	//infoColor := cs.Cyan
	//
	//fmt.Fprintln(f.IOStreams.Out, infoColor(url))
}

type Share struct {
	FileID int `json:"file_id"`
}

func createShare(f *factory.Factory, fileID int) {
	cfg, _ := f.Config()
	cs := f.IOStreams.ColorScheme()

	client, _ := f.HttpClient()

	hostname, _ := cfg.DefaultHost()
	url := fmt.Sprintf("%s://%s/v1/shares", f.HttpSchema, hostname)

	marshal, err := json.Marshal(Share{FileID: fileID})
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(marshal))
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	respData, _ := io.ReadAll(resp.Body)
	type Link struct {
		Link string
	}
	var l Link

	json.Unmarshal(respData, &l)

	url = fmt.Sprintf("%s://%s/%s", f.HttpSchema, hostname, l.Link)
	fmt.Fprintf(f.IOStreams.Out, "url: %s\n", cs.Cyan(url))

	return
}
