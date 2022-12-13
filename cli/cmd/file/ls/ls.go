package ls

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	"github.com/bitmyth/pdrive-cli/cli/tableprinter"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	Files   []File
	Dir     = 0
	Root    = File{ID: 0, Dir: 0}
	Current = Root
)

func Cd(dir File) {
	if dir.Parent == nil {
		cpy := Current
		dir.Parent = &cpy
	}
	Current = dir

	Dir = int(dir.ID)
}

func Ls(f *factory.Factory) {
	client, _ := f.HttpClient()

	cfg, _ := f.Config()
	hostname, _ := cfg.DefaultHost()

	query := url.Values{}
	query.Add("size", "100")
	//query.Add("sort", `{"id":"desc"}`)
	query.Add("dir", strconv.Itoa(int(Dir)))

	apiUrl := fmt.Sprintf("%s://%s/v1/files?"+`sort={"id":"desc"}&`+query.Encode(), f.HttpSchema, hostname)

	req, _ := http.NewRequest(http.MethodGet, apiUrl, nil)
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	respData, _ := io.ReadAll(resp.Body)
	var files IndexResp
	err = json.Unmarshal(respData, &files)
	if err != nil {
		println(string(respData))
	}
	Files = files.Files

	cs := f.IOStreams.ColorScheme()
	infoColor := cs.Cyan

	tp := tableprinter.New(f.IOStreams)
	tp.HeaderRow("ID", "Name", "Size", "Path", "Date")
	for _, f := range files.Files {
		tp.AddField(fmt.Sprintf("%d", f.ID))
		if f.IsDir {
			tp.AddField(f.Name, tableprinter.WithColor(cs.Yellow))
		} else {
			tp.AddField(f.Name, tableprinter.WithColor(infoColor))
		}
		tp.AddField(fmt.Sprintf("%d", f.Size))
		tp.AddField(f.Path)
		tp.AddField(f.CreatedAt.Format("2006-01-02 15:04:05"))
		tp.EndRow()
	}
	tp.Render()
}

type IndexResp struct {
	Files []File `json:"files,omitempty"`
	Count int64  `json:"count,omitempty"`
}

type DeletedAt sql.NullTime

type File struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt DeletedAt `gorm:"index"`
	Name      string    `json:"name"`
	UserID    uint      `json:"user_id"`
	Mime      string    `json:"mime"`
	Size      int64     `json:"size"`
	Path      string    `json:"path"`
	Dir       uint      `json:"dir"`
	IsDir     bool      `json:"is_dir"`
	Parent    *File     `json:"-"`
}
