package cmd

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path"

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

		readerCloser, err := getReaderCloser(deta, args[0])
		if err != nil {
			return fmt.Errorf("could not get reader: %w", err)
		}
		defer readerCloser.Close()

		if err := copyContent(deta, readerCloser, args[1]); err != nil {
			return fmt.Errorf("could not copy file: %w", err)
		}

		return nil
	},
}

func getReaderCloser(deta *deta.Deta, src string) (io.ReadCloser, error) {
	if src == "-" {
		return os.Stdin, nil
	}

	if info, err := os.Stat(src); err == nil {
		if info.IsDir() {
			return nil, fmt.Errorf("reading from directories is not supported")
		}

		return os.Open(src)
	}

	srcUrl, err := url.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("could not parse source url: %w", err)
	}

	if srcUrl.Scheme != "deta" {
		return nil, fmt.Errorf("source must be a deta url")
	}

	disk, err := drive.New(deta, srcUrl.Host)
	if err != nil {
		return nil, fmt.Errorf("could not get drive: %w", err)
	}

	reader, err := disk.Get(srcUrl.Path)
	if err != nil {
		return nil, fmt.Errorf("could not get file: %w", err)
	}

	return reader, nil
}

func copyContent(deta *deta.Deta, reader io.Reader, dest string) error {
	if dest == "-" {
		if _, err := io.Copy(os.Stdout, reader); err != nil {
			return fmt.Errorf("could not read file: %w", err)
		}
		return nil
	}

	dirname := path.Dir(dest)
	if _, err := os.Stat(dirname); err == nil {
		writer, err := os.Create(dest)
		if err != nil {
			return fmt.Errorf("could not create file: %w", err)
		}
		defer writer.Close()

		if _, err = io.Copy(writer, reader); err != nil {
			return fmt.Errorf("could not read file: %w", err)
		}

		return nil
	}

	destUrl, err := url.Parse(dest)
	if err != nil {
		return fmt.Errorf("could not parse destination url: %w", err)
	}

	if destUrl.Scheme != "deta" {
		return fmt.Errorf("destination must be a deta url")
	}

	disk, err := drive.New(deta, destUrl.Host)
	if err != nil {
		return fmt.Errorf("could not get drive: %w", err)
	}

	if _, err := disk.Put(&drive.PutInput{
		Name: destUrl.Path,
		Body: reader,
	}); err != nil {
		return fmt.Errorf("could not put file: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(cpCmd)
}
