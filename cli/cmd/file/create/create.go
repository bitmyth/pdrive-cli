package create

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	"github.com/bitmyth/pdrive-cli/cli/config"
	"github.com/bitmyth/pdrive-cli/cli/iostreams"
	"github.com/spf13/cobra"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

type Options struct {
	HttpClient func() (*http.Client, error)
	IO         *iostreams.IOStreams
	Config     config.Config
	Org        string
	File       string
	Dir        string
	HttpSchema string
	Exclude    []string
}

func NewCmdCreate(f *factory.Factory) *cobra.Command {
	opts := &Options{
		IO:         f.IOStreams,
		HttpSchema: f.HttpSchema,
	}
	opts.HttpClient = f.HttpClient
	opts.IO = f.IOStreams
	opts.Config, _ = f.Config()
	cmd := &cobra.Command{
		Use: "create --name | -n",
		RunE: func(cmd *cobra.Command, args []string) error {

			info := FileInfo{
				FileName: "test",
				Content:  ReadStdIn(),
				Dir:      "",
				FileSize: 0,
			}

			return createFile(opts, info)
		},
	}
	cmd.Flags().StringVarP(&opts.File, "file", "f", "", "absolute file path")

	return cmd
}

func ReadStdIn() string {
	scanner := bufio.NewScanner(os.Stdin)
	var buf bytes.Buffer
	for {
		stopped := scanner.Scan()
		if !stopped {
			if err := scanner.Err(); err != nil {
				if err != io.EOF {
					fmt.Fprintln(os.Stderr, err)
					println("END")
				}
				break
			}
			break
		}

		buf.WriteString("\n")
		buf.Write(scanner.Bytes())
	}
	//fmt.Println(buf.String())
	return buf.String()
}

func createFile(opts *Options, info FileInfo) error {
	cs := opts.IO.ColorScheme()
	infoColor := cs.Cyan
	client, _ := opts.HttpClient()

	cfg := opts.Config
	hostname, _ := cfg.DefaultHost()

	url := fmt.Sprintf("%s://%s/v1/files", opts.HttpSchema, hostname)

	// New multipart writer.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("file", info.Name())
	if err != nil {
		_, err = io.Copy(fw, strings.NewReader(info.Content))
		if err != nil {
			return err
		}
		return err
	}

	field, _ := writer.CreateFormField("name")
	field.Write([]byte(info.Name()))

	field, _ = writer.CreateFormField("dir")
	field.Write([]byte(info.Dir))

	field, _ = writer.CreateFormField("is_dir")
	field.Write([]byte("false"))

	writer.Close()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body.Bytes()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	fmt.Fprintln(opts.IO.Out, fmt.Sprintf("Uploading %s", infoColor(info.Name())))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(opts.IO.Out, cs.Red(err.Error()))
		return err
	}
	code := resp.Status
	if resp.StatusCode == http.StatusCreated {
		fmt.Fprintf(opts.IO.Out, "%s uploaded\n", cs.SuccessIcon())
	} else {
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Fprintln(opts.IO.Out, cs.WarningIcon(), cs.Blue(code), cs.Yellow(string(respBody)))
	}
	return nil

}

type FileInfo struct {
	FileName string
	Content  string
	Dir      string
	FileSize int64
}

func (f FileInfo) Name() string {
	return f.FileName
}

func (f FileInfo) Size() int64 {
	return f.FileSize
}

func (f FileInfo) Mode() fs.FileMode {
	return fs.ModePerm
}

func (f FileInfo) ModTime() time.Time {
	return time.Now()
}

func (f FileInfo) IsDir() bool {
	return false
}

func (f FileInfo) Sys() any {
	return nil
}