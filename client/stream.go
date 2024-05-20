package client

import (
	"fmt"
	"net"
	"net/http"
	"nir/finder"
	"nir/nyaa"
	"nir/player"
	"path/filepath"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
)

func (client *Client) StreamTorrent(nyaaTorrent *nyaa.NyaaTorrent, p *player.Player) error {
	torr, err := client.AddMagnet(nyaaTorrent.Magnet)
	if err != nil {
		return err
	}

	fmt.Printf("Fetching %s info...\n", torr.Name())
	<-torr.GotInfo()

	var reader torrent.Reader
	var filename string

	var playableFiles []*torrent.File

	for _, file := range torr.Files() {
		if strings.HasSuffix(file.Path(), ".mkv") || strings.HasSuffix(file.Path(), ".mp4") {
			playableFiles = append(playableFiles, file)
		}
	}

	if len(playableFiles) == 0 {
		return fmt.Errorf("no playable files in torrent")
	}

	if len(playableFiles) == 1 {
		reader = playableFiles[0].NewReader()
		filename = filepath.Base(playableFiles[0].Path())
	} else {
		file, err := finder.FindFiles(playableFiles)
		if err != nil {
			return err
		}
		reader = file.NewReader()
		filename = filepath.Base(file.Path())
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}

	portNum := listener.Addr().(*net.TCPAddr).Port

	http.HandleFunc("GET /stream", func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, filename, time.Time{}, reader)
	})

	go http.Serve(listener, nil)

	if err = p.Play(fmt.Sprintf("http://127.0.0.1:%d/stream", portNum), filename); err != nil {
		return err
	}

	return nil
}
