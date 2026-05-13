//go:build !updater

package updater

import (
	"context"

	"github.com/spf13/cobra"
)

type UpdateInfo struct {
	Version   string
	Notes     string
	URL       string
	Published string
}

func CheckForUpdates(ctx context.Context, currentVersion string) (*UpdateInfo, error) {
	return nil, nil
}

func UpdateSelfFromRelease(url, token, hash, binary string) error {
	return nil
}

func GetVersionInfo() (string, string) {
	return "", ""
}

func NewUpdateCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update is disabled (build without -tags updater)",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Update feature is disabled in this build.")
			cmd.Println("To enable, rebuild with: go build -tags updater")
		},
	}
	return cmd
}