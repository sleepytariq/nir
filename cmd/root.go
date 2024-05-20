package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version string = "0.1.0"

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	downloadCmd.Flags().BoolP("trusted-only", "t", false, "show only torrents from trusted sources")
	downloadCmd.Flags().StringP("output-directory", "o", "", "path to output directory to save torrent files in")
	downloadCmd.Flags().Int("queue", 3, "maximum number of torrents downloaded concurrently")
	rootCmd.AddCommand(downloadCmd)

	streamCmd.Flags().BoolP("trusted-only", "t", false, "show only torrents from trusted sources")
	streamCmd.Flags().Bool("vlc", false, "use vlc instead of mpv")
	rootCmd.AddCommand(streamCmd)

	rootCmd.AddCommand(clearCmd)
}

var rootCmd = &cobra.Command{
	Use:     "nir",
	Short:   "Download or stream torrents from nyaa.si",
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
