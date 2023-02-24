package cmd

import (
	"fmt"
	"net/url"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/drive"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm [url]",
	Short: "remove a file from the drive",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		deta, ok := cmd.Context().Value(detaKey).(*deta.Deta)
		if !ok {
			return fmt.Errorf("could not get drive from context")
		}

		url, err := url.Parse(args[0])
		if err != nil {
			return fmt.Errorf("could not parse url: %w", err)
		}

		if url.Scheme != "deta" {
			return fmt.Errorf("url must be a deta url")
		}

		drive, err := drive.New(deta, url.Host)
		if err != nil {
			return fmt.Errorf("could not get drive: %w", err)
		}

		if _, err := drive.Delete(url.Path); err != nil {
			return fmt.Errorf("could not delete file: %w", err)
		}

		fmt.Printf("deleted %s", url)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
