package finder

import (
	"fmt"
	"nir/nyaa"

	"github.com/anacrolix/torrent"
	"github.com/ktr0731/go-fuzzyfinder"
)

func FindTorrents(torrents *[]nyaa.NyaaTorrent) (*[]*nyaa.NyaaTorrent, error) {
	idx, err := fuzzyfinder.FindMulti(*torrents, func(i int) string {
		return (*torrents)[i].Name
	}, fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
		if i == -1 {
			return ""
		}
		return fmt.Sprintf(
			"Seeders: %s\nLeechers: %s\nDownloads: %s\nSize: %s\nTrusted: %s",
			(*torrents)[i].Seeders,
			(*torrents)[i].Leechers,
			(*torrents)[i].Downloads,
			(*torrents)[i].Size,
			(*torrents)[i].Trusted,
		)
	}),
		fuzzyfinder.WithPromptString("Select Torrents> "),
	)

	if err != nil {
		return nil, fmt.Errorf("no torrents were selected")
	}

	var selectedTorrents []*nyaa.NyaaTorrent

	for _, val := range idx {
		selectedTorrents = append(selectedTorrents, &(*torrents)[val])
	}

	return &selectedTorrents, nil
}

func FindTorrent(torrents *[]nyaa.NyaaTorrent) (*nyaa.NyaaTorrent, error) {
	idx, err := fuzzyfinder.Find(*torrents, func(i int) string {
		return (*torrents)[i].Name
	}, fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
		if i == -1 {
			return ""
		}
		return fmt.Sprintf(
			"Seeders: %s\nLeechers: %s\nDownloads: %s\nSize: %s\nTrusted: %s",
			(*torrents)[i].Seeders,
			(*torrents)[i].Leechers,
			(*torrents)[i].Downloads,
			(*torrents)[i].Size,
			(*torrents)[i].Trusted,
		)
	}),
		fuzzyfinder.WithPromptString("Select a Torrent> "),
	)

	if err != nil {
		return nil, fmt.Errorf("no torrent was selected")
	}

	return &(*torrents)[idx], nil
}

func FindFiles(files []*torrent.File) (*torrent.File, error) {
	idx, err := fuzzyfinder.Find(files, func(i int) string {
		return files[i].Path()
	},
		fuzzyfinder.WithPromptString("Select a File> "),
		fuzzyfinder.WithHeader("Torrent contains multiple files"),
	)

	if err != nil {
		return nil, fmt.Errorf("no file was selected")
	}

	return files[idx], nil
}
