package mkdir

import (
	"bytes"
	"fmt"
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

func Mkdir(f *factory.Factory, name string, dir int) error {
	cs := f.IOStreams.ColorScheme()
	out := f.IOStreams.Out
	infoColor := cs.Cyan
	client, _ := f.HttpClient()

	cfg, _ := f.Config()
	hostname, _ := cfg.DefaultHost()

	url := fmt.Sprintf("%s://%s/v1/files", f.HttpSchema, hostname)

	// New multipart writer.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_, err := writer.CreateFormFile("file", name)

	field, _ := writer.CreateFormField("name")
	field.Write([]byte(name))

	field, _ = writer.CreateFormField("dir")
	field.Write([]byte(strconv.Itoa(dir)))

	field, _ = writer.CreateFormField("is_dir")
	field.Write([]byte("true"))

	writer.Close()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body.Bytes()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(out, cs.Red(err.Error()))
		return err
	}
	code := resp.Status
	if resp.StatusCode == http.StatusCreated {
		fmt.Fprintf(out, "%s created dir %s\n", cs.SuccessIcon(), infoColor(name))
	} else {
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Fprintln(out, cs.WarningIcon(), cs.Blue(code), cs.Yellow(string(respBody)))
	}
	return nil

}
