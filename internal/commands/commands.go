package commands

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/spf13/cobra"
)

func RegisterCommands(cmd *cobra.Command, app App) {
	handler := CreateImportHandler(app)

	cmd.AddCommand(&cobra.Command{
		Use:   "import",
		Short: "Import users",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] != "l" {
				return fmt.Errorf("Invalid argument: %v\n", args[0])
			}
			_, err := url.ParseRequestURI(args[1])
			if err != nil {
				return err
			}
			start := time.Now()
			err = handler.ImportUserByLink(context.Background(), args[1])
			fmt.Printf("Import took %v\n", time.Since(start))
			return err
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:     "ImportPosts",
		Short:   "Import posts",
		Long:    "Import posts to user. If multiple users provide then posts will be assigned randomly to provided users",
		Example: "ImportPosts https://posts.com/json 019ad582-424c-7316-b98c-eee5bf0f8d08 019ad582-42b3-7340-a5b8-6025ce042c01",
		Args:    cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			_, err := url.ParseRequestURI(args[0])
			if err != nil {
				return err
			}
			start := time.Now()
			err = handler.ImportPosts(context.Background(), args[0], args[1:])
			fmt.Printf("Import took %v\n", time.Since(start))
			return err
		},
	})
}
