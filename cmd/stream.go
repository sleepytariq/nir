package cmd

import (
	"fmt"
	"nir/client"
	"nir/finder"
	"nir/nyaa"
	"nir/player"
	"os"

	"github.com/spf13/cobra"
)

var streamCmd = &cobra.Command{
	Use:     "stream [flags] [query]",
	Short:   "stream torrents from nyaa.si",
	Aliases: []string{"s"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		trustedOnly, _ := cmd.Flags().GetBool("trusted-only")
		useVLC, _ := cmd.Flags().GetBool("vlc")
		query := args[0]

		torrents, err := nyaa.SearchNyaa(query, trustedOnly)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		selectedTorrent, err := finder.FindTorrent(torrents)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		client, err := client.NewClient(1)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer client.Close()

		var p *player.Player
		if useVLC {
			p = &player.VLC
		} else {
			p = &player.MPV
		}

		err = client.StreamTorrent(selectedTorrent, p)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
