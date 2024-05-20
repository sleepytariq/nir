package client

import (
	"fmt"
	"nir/cache"

	"github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
)

type Client struct {
	*torrent.Client
	Sem chan bool
}

func NewClient(queueSize int) (*Client, error) {
	cfg := torrent.NewDefaultClientConfig()
	cfg.Logger.SetHandlers(log.DiscardHandler)

	cacheDir, err := cache.GetCacheDir("nir")
	if err != nil {
		return nil, err
	}

	cfg.DataDir = cacheDir

	torrentClient, err := torrent.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create client")
	}

	sem := make(chan bool, queueSize)

	return &Client{torrentClient, sem}, nil
}
