package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	appVersion string
	verbose    bool
)

var purple = color.New(color.FgHiMagenta)
var dim = color.New(color.Faint)

var rootCmd = &cobra.Command{
	Use:   "storeready",
	Short: "Pre-submission compliance scanner for Apple App Store and Google Play",
	Long: fmt.Sprintf(`%s

StoreReady scans your app against Apple App Store and Google Play submission
requirements before you submit, catching rejection risks so you ship with confidence.

Get started:
  storeready appstore-checkup .    Apple local preflight checks (offline)
  storeready playstore-checkup .   Google Play local checks (offline)
  storeready play-guidelines list  Browse Google Play policy matrix
  storeready preflight . --ipa X   Apple preflight with IPA binary analysis
  storeready scan --app-id ID      Apple App Store Connect metadata (needs API key)
  storeready guidelines search     Browse Apple's review guidelines`,
		purple.Sprint("storeready — know before you submit.")),
}

func SetVersion(v string) {
	appVersion = v
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(guidelinesCmd)
	rootCmd.AddCommand(versionCmd)
}
