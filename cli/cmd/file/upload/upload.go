package upload

import (
	"bytes"
	"fmt"
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	"github.com/bitmyth/pdrive-cli/cli/config"
	"github.com/bitmyth/pdrive-cli/cli/iostreams"
	"github.com/spf13/cobra"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
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

func NewCmdUpload(f *factory.Factory) *cobra.Command {
	opts := &Options{
		IO:         f.IOStreams,
		HttpSchema: f.HttpSchema,
	}
	opts.HttpClient = f.HttpClient
	opts.IO = f.IOStreams
	opts.Config, _ = f.Config()
	cmd := &cobra.Command{
		Use: "upload --file | -f",
		RunE: func(cmd *cobra.Command, args []string) error {
			return uploadRun(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.File, "file", "f", "", "absolute file path")

	return cmd
}

func RunUploadFile(f *factory.Factory, filePath string, dir int) {
	opts := &Options{
		IO:         f.IOStreams,
		HttpSchema: f.HttpSchema,
		File:       filePath,
		Dir:        strconv.Itoa(dir),
	}
	opts.HttpClient = f.HttpClient
	opts.IO = f.IOStreams
	opts.Config, _ = f.Config()

	uploadRun(opts)
}

func uploadRun(opts *Options) error {
	cs := opts.IO.ColorScheme()
	errColor := cs.Red

	file, err := os.Open(opts.File)
	if err != nil {
		fmt.Fprintln(opts.IO.Out, errColor(err.Error()))
		return err
	}
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if stat.IsDir() {
		err = uploadDir(opts)
		if err != nil {
			fmt.Fprintln(opts.IO.Out, errColor(err.Error()))
			return err
		}
	} else {
		err = uploadFile(opts, FileInfo{FileInfo: stat, Path: opts.File, Dir: opts.Dir})
		if err != nil {
			return err
		}
	}

	return nil
}

func uploadFile(opts *Options, info FileInfo) error {
	cs := opts.IO.ColorScheme()
	infoColor := cs.Cyan
	client, _ := opts.HttpClient()

	cfg := opts.Config
	hostname, _ := cfg.DefaultHost()

	url := fmt.Sprintf("%s://%s/v1/files", opts.HttpSchema, hostname)
	file, err := os.Open(info.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	// New multipart writer.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return err
	}

	if !info.FileInfo.IsDir() {
		_, err = io.Copy(fw, file)
		if err != nil {
			return err
		}
	}

	field, _ := writer.CreateFormField("name")
	field.Write([]byte(file.Name()))

	field, _ = writer.CreateFormField("dir")
	field.Write([]byte(info.Dir))

	field, _ = writer.CreateFormField("is_dir")
	if info.FileInfo.IsDir() {
		field.Write([]byte("true"))
	} else {
		field.Write([]byte("false"))
	}

	writer.Close()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body.Bytes()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	fmt.Fprintln(opts.IO.Out, fmt.Sprintf("Uploading %s", infoColor(file.Name())))
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

type NotifiableReader struct {
	s        []byte
	i        int64 // current reading index
	prevRune int   // index of previous rune; or < 0
}

func NewNotifiableReader(b []byte) *NotifiableReader { return &NotifiableReader{b, 0, -1} }

// Read implements the io.Reader interface.
func (r *NotifiableReader) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	r.prevRune = -1
	n = copy(b, r.s[r.i:])
	r.i += int64(n)
	return
}

func (r *NotifiableReader) WriteTo(w io.Writer) (n int64, err error) {
	r.prevRune = -1
	size := int64(len(r.s))
	if r.i >= size {
		return 0, nil
	}

	for r.i < int64(len(r.s)) {
		fmt.Printf("\033[2K write [%d-%d) of %d %.2f%% \033[0G", r.i, r.i+1000, size, 100*float64(r.i)/float64(size))

		b := r.s[r.i : r.i+1000]
		m, err := w.Write(b)
		if m > len(b) {
			panic("bytes.Reader.WriteTo: invalid Write count")
		}
		r.i += int64(m)
		n = int64(m)
		if m != len(b) && err == nil {
			err = io.ErrShortWrite
		}
	}

	return
}
