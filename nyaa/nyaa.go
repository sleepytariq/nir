package nyaa

import (
	"fmt"
	"net/url"

	"github.com/mmcdole/gofeed"
)

type NyaaTorrent struct {
	Name      string
	Magnet    string
	Seeders   string
	Leechers  string
	Downloads string
	Size      string
	Trusted   string
}

func SearchNyaa(query string, trustedOnly bool) (*[]NyaaTorrent, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be an empty string")
	}

	query = url.QueryEscape(query)

	var trustedOnlyValue int

	if trustedOnly {
		trustedOnlyValue = 2
	} else {
		trustedOnlyValue = 0
	}

	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(fmt.Sprintf("https://nyaa.si?page=rss&q=%s&f=%d&m", query, trustedOnlyValue))

	if err != nil {
		return nil, fmt.Errorf("failed to query nyaa")
	}

	if len(feed.Items) == 0 {
		return nil, fmt.Errorf("no results were found")
	}

	var torrents []NyaaTorrent

	for _, item := range feed.Items {
		title := item.Title
		magent := item.Link
		seeders := item.Extensions["nyaa"]["seeders"][0].Value
		leechers := item.Extensions["nyaa"]["leechers"][0].Value
		downloads := item.Extensions["nyaa"]["downloads"][0].Value
		size := item.Extensions["nyaa"]["size"][0].Value
		trusted := item.Extensions["nyaa"]["trusted"][0].Value
		torrents = append(torrents, NyaaTorrent{
			Name:      title,
			Magnet:    magent,
			Seeders:   seeders,
			Leechers:  leechers,
			Downloads: downloads,
			Size:      size,
			Trusted:   trusted,
		})
	}

	return &torrents, nil
}
