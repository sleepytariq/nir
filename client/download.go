package client

import (
	"fmt"
	"io"
	"nir/nyaa"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gosuri/uilive"
)

func (client *Client) DownloadTorrentFiles(torrents *[]*nyaa.NyaaTorrent, outputDirectory string) error {
	for _, torr := range *torrents {
		_, err := client.AddMagnet(torr.Magnet)
		if err != nil {
			return err
		}
	}

	for _, torr := range client.Torrents() {
		fmt.Printf("Fetching %s info...\n", torr.Name())
		<-torr.GotInfo()
	}

	var wg sync.WaitGroup
	done := make(chan bool, 1)

	wg.Add(1)
	go func() { defer wg.Done(); client.ShowProgress(done) }()

	for _, torr := range client.Torrents() {
		client.Sem <- true
		go func() {
			defer func() { <-client.Sem }()
			torr.DownloadAll()
			<-torr.Complete.On()
		}()
	}

	client.WaitAll()
	done <- true
	wg.Wait()

	client.CopyTorrents(outputDirectory)

	fmt.Printf("Successfully downloaded %d torrents\n", len(*torrents))

	return nil
}

func (client *Client) CopyTorrents(dst string) error {
	for _, torr := range client.Torrents() {
		for _, file := range torr.Files() {
			dir := filepath.Dir(file.Path())

			if err := os.MkdirAll(filepath.Join(dst, dir), os.ModePerm); err != nil {
				return err
			}

			if err := func() error {
				fp, err := os.Create(filepath.Join(filepath.Join(dst, dir), filepath.Base(file.Path())))
				if err != nil {
					return err
				}
				fileReader := file.NewReader()
				io.Copy(fp, fileReader)
				return nil
			}(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (client *Client) ShowProgress(done chan bool) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	writer := uilive.New()
	writer.Start()
	defer writer.Stop()

	torrents := client.Torrents()

	for range ticker.C {
		select {
		case <-done:
			fmt.Fprint(writer, "\r")
			return
		default:
			msg := ""
			for _, torr := range torrents {
				if torr.Complete.Bool() {
					msg += fmt.Sprintf("Completed %s\n", torr.Name())
				} else {
					currBytes := humanize.IBytes(uint64(torr.BytesCompleted()))
					totalBytes := humanize.IBytes(uint64(torr.Info().TotalLength()))
					msg += fmt.Sprintf("Downloading %s [%s/%s]\n", torr.Name(), currBytes, totalBytes)
				}
			}
			fmt.Fprint(writer, msg)
		}
	}
}
