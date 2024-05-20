package cmd

import (
	"fmt"
	"nir/cache"
	"os"

	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:     "clear",
	Short:   "clear cached torrents",
	Aliases: []string{"c"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cachedir, err := cache.GetCacheDir("nir")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err = os.RemoveAll(cachedir); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Removed cached torrents")
	},
}
