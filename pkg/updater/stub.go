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
	return nil
}