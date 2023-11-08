package download

import (
	"fmt"
	"github.com/bitmyth/pdrive-cli/cli/cmd/factory"
	"github.com/bitmyth/pdrive-cli/cli/cmd/file/ls"
	"io"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
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
	link := fmt.Sprintf("%s://%s/v1/files/download?"+val, f.HttpSchema, hostname)

	resp, err := client.Get(link)
	if err != nil {
		return
	}

	//respData, _ := io.ReadAll(resp.Body)

	infoColor := cs.Cyan
	// save file
	cd := resp.Header.Get("Content-Disposition")
	//println("Content-Disposition ", cd)
	var name string
	if cd != "" {
		name = strings.Split(cd, "=")[1]
	} else {
		name = id
	}

	fmt.Fprintln(f.IOStreams.Out, "Downloading ", infoColor(name))

	saveFile, err := os.Create(name)
	if err != nil {
		return
	}

	//_, _ = create.Write(respData)
	gauge := NewSpeedGauge()
	go gauge.Show()
	size, _ := io.Copy(saveFile, io.TeeReader(resp.Body, gauge))
	_ = saveFile.Close()
	gauge.Stop()
	<-gauge.stopped

	fmt.Fprintln(f.IOStreams.Out, "Size:", infoColor(fmt.Sprintf("%d", size)))
}

func NewSpeedGauge() *SpeedGauge {
	return &SpeedGauge{
		stop:    make(chan struct{}),
		stopped: make(chan struct{}, 1),
	}
}

type SpeedGauge struct {
	count   int64
	stop    chan struct{}
	stopped chan struct{}
}

func (s *SpeedGauge) Write(b []byte) (int, error) {
	c := len(b)
	atomic.AddInt64(&s.count, int64(c))
	return c, nil
}

func (s *SpeedGauge) Stop() {
	close(s.stop)
}

func (s *SpeedGauge) Show() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		select {
		case <-s.stop:
			fmt.Print("\r")
			s.stopped <- struct{}{}
			return
		default:
		}
		fmt.Printf("\r %s/s", ls.ByteSize(atomic.LoadInt64(&s.count)).String())
		atomic.StoreInt64(&s.count, 0)
	}
}
