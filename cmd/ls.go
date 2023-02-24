package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/drive"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls [url]",
	Short: "list files in the drive",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		deta, ok := cmd.Context().Value(detaKey).(*deta.Deta)
		if !ok {
			return fmt.Errorf("could not get drive from context")
		}

		recursive, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return fmt.Errorf("could not get recursive flag: %w", err)
		}

		entryUrl, err := url.Parse(args[0])
		if err != nil {
			return fmt.Errorf("could not parse url: %w", err)
		}

		if entryUrl.Scheme != "deta" {
			return fmt.Errorf("url must be a deta url")
		}

		disk, err := drive.New(deta, entryUrl.Host)
		if err != nil {
			return fmt.Errorf("could not get drive: %w", err)
		}

		prefix := entryUrl.Path
		if !strings.HasSuffix(prefix, "/") {
			prefix = fmt.Sprintf("%s/", prefix)
		}

		out, err := disk.List(1000, prefix, "")
		if err != nil {
			return fmt.Errorf("could not list files: %w", err)
		}

		for _, name := range out.Names {
			if !recursive {
				relativePath := strings.TrimPrefix(name, prefix)
				if strings.Contains(relativePath, "/") {
					continue
				}
			}
			fmt.Printf("%s%s\n", entryUrl, name)
		}

		return nil
	},
}

func init() {
	lsCmd.Flags().BoolP("recursive", "r", false, "list files recursively")

	rootCmd.AddCommand(lsCmd)
}
