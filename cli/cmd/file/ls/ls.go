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
	Page    = 1
)

func Cd(dir File) {
	Page = 1
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

	cs := f.IOStreams.ColorScheme()
	infoColor := cs.Cyan
	hostname, _ := cfg.DefaultHost()

	query := url.Values{}
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
	var files IndexResp
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
			tp.AddField(fmt.Sprintf("%d (%s)", f.Size, ByteSize(f.Size).String()))
		} else {
			tp.AddField(fmt.Sprintf("%d", f.Size))
		}
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

type ByteSize float64

const (
	KB ByteSize = 1 << (10 * (iota + 1))
	MB
	GB
	TB
)

func (b ByteSize) String() string {
	switch {
	case b >= TB:
		return fmt.Sprintf("%.1fT", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.1fG", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.1fM", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.1fK", b/KB)
	}
	return fmt.Sprintf("%.1fB", b)
}
