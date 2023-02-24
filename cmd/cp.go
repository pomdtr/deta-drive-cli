package cmd

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/drive"
	"github.com/spf13/cobra"
)

var cpCmd = &cobra.Command{
	Use:   "cp [source] [destination]",
	Short: "copy a file from one drive to another",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		deta, ok := cmd.Context().Value(detaKey).(*deta.Deta)
		if !ok {
			return fmt.Errorf("could not get drive from context")
		}

		srcUrl, err := url.Parse(args[0])
		if err != nil {
			return fmt.Errorf("could not parse source url: %w", err)
		}

		destUrl, err := url.Parse(args[1])
		if err != nil {
			return fmt.Errorf("could not parse destination url: %w", err)
		}

		if srcUrl.Scheme != "deta" && destUrl.Scheme != "deta" {
			return fmt.Errorf("source or destination must be a deta url")
		}

		if srcUrl.Scheme != "deta" {
			disk, err := drive.New(deta, destUrl.Host)
			if err != nil {
				return fmt.Errorf("could not get drive: %w", err)
			}

			file, err := os.Open(srcUrl.Path)
			if err != nil {
				return fmt.Errorf("could not read file: %w", err)
			}

			if _, err = disk.Put(&drive.PutInput{
				Name: destUrl.Path,
				Body: file,
			}); err != nil {
				return fmt.Errorf("could not put file: %w", err)
			}

			fmt.Println("file uploaded")
			return nil
		}

		if destUrl.Scheme != "deta" {
			disk, err := drive.New(deta, srcUrl.Host)
			if err != nil {
				return fmt.Errorf("could not get drive: %w", err)
			}

			file, err := os.Create(destUrl.Path)
			if err != nil {
				return fmt.Errorf("could not create file: %w", err)
			}

			reader, err := disk.Get(srcUrl.Path)
			if err != nil {
				return fmt.Errorf("could not get file: %w", err)
			}

			if _, err := io.Copy(file, reader); err != nil {
				return fmt.Errorf("could not copy file: %w", err)
			}

			return nil
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cpCmd)
}
