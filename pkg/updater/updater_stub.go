//go:build !selfupdate

package updater

import (
	"errors"

	"github.com/spf13/cobra"
)

// UpdateSelfFromRelease is a stub that returns an error when selfupdate is disabled.
func UpdateSelfFromRelease(releaseURL, platform, arch, binaryPath string) error {
	return errors.New("self-update not available: rebuild with -tags selfupdate")
}

// NewUpdateCommand returns a command that indicates self-update is disabled.
func NewUpdateCommand(binaryName string) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Self-update (disabled)",
		Long:  "Self-update is disabled in this build. Rebuild with -tags selfupdate to enable.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("self-update not available: rebuild with -tags selfupdate")
		},
	}
}