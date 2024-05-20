package cmd

import (
	"fmt"
	"nir/client"
	"nir/finder"
	"nir/nyaa"
	"os"

	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:     "download [flags] [query]",
	Short:   "download torrents from nyaa.si",
	Aliases: []string{"d"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		trustedOnly, _ := cmd.Flags().GetBool("trusted-only")
		outputDirectory, _ := cmd.Flags().GetString("output-directory")
		queueSize, _ := cmd.Flags().GetInt("queue")
		query := args[0]

		torrents, err := nyaa.SearchNyaa(query, trustedOnly)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		selectedTorrent, err := finder.FindTorrents(torrents)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		client, err := client.NewClient(queueSize)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer client.Close()

		if err = client.DownloadTorrentFiles(selectedTorrent, outputDirectory); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
