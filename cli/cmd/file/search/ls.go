package search

import (
	"encoding/json"
	"fmt"
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	"github.com/bitmyth/pdrive-cli/cli/cmd/file/ls"
	"github.com/bitmyth/pdrive-cli/cli/tableprinter"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

var (
	Files   []ls.File
	Dir     = 0
	Root    = ls.File{ID: 0, Dir: 0}
	Current = Root
	Page    = 1
)

func Search(f *factory.Factory, s string) {
	client, _ := f.HttpClient()

	cfg, _ := f.Config()

	cs := f.IOStreams.ColorScheme()
	infoColor := cs.Cyan
	hostname, _ := cfg.DefaultHost()

	query := url.Values{}
	query.Add("name", s)
	query.Add("size", "20")
	query.Add("page", strconv.Itoa(Page))
	//query.Add("sort", `{"id":"desc"}`)
	query.Add("dir", strconv.Itoa(int(Dir)))

	apiUrl := fmt.Sprintf("%s://%s/v1/files?"+`sort={"id":"desc"}&`+query.Encode(), f.HttpSchema, hostname)

	req, _ := http.NewRequest(http.MethodGet, apiUrl, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(f.IOStreams.Out, cs.Red(err.Error()))
		return
	}

	respData, _ := io.ReadAll(resp.Body)
	var files ls.IndexResp
	err = json.Unmarshal(respData, &files)
	if err != nil {
		println(string(respData))
	}
	Files = files.Files

	tp := tableprinter.New(f.IOStreams)
	tp.HeaderRow("ID", "Name", "Size", "Path", "Date")
	for _, f := range files.Files {
		tp.AddField(fmt.Sprintf("%d", f.ID))
		if f.IsDir {
			tp.AddField(f.Name, tableprinter.WithColor(cs.Yellow))
		} else {
			tp.AddField(f.Name, tableprinter.WithColor(infoColor))
		}
		if f.Size > 1<<10 {
			tp.AddField(fmt.Sprintf("%d (%s)", f.Size, ls.ByteSize(f.Size).String()))
		} else {
			tp.AddField(fmt.Sprintf("%d", f.Size))
		}
		tp.AddField(f.Path)
		tp.AddField(f.CreatedAt.Format("2006-01-02 15:04:05"))
		tp.EndRow()
	}
	tp.Render()
}
