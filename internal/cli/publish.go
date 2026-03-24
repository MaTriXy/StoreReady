package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	ascapi "github.com/MaTriXy/StoreReady/internal/asc"
	"github.com/MaTriXy/StoreReady/internal/checks"
	"github.com/MaTriXy/StoreReady/internal/config"
	"github.com/MaTriXy/StoreReady/internal/preflight"
	"github.com/spf13/cobra"
)

var (
	publishAppID       string
	publishVersion     string
	publishBuild       string
	publishMetadataDir string
	publishPath        string
	publishIPA         string
	publishScanTier    int
	publishSkipLocal   bool
	publishSkipASCScan bool
	publishConfirm     bool
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Run StoreReady gates, then publish with ASC CLI",
	Long: `End-to-end lane for App Store release:
  1) Run local StoreReady preflight gates
  2) Run ASC metadata gates via StoreReady scan
  3) Execute 'asc release run' to continue submit flow

By default this command runs ASC in dry-run mode.
Use --confirm to execute the real submission flow.`,
	RunE: runPublish,
}

func init() {
	publishCmd.Flags().StringVar(&publishAppID, "app-id", "", "App Store Connect app ID (required)")
	publishCmd.Flags().StringVar(&publishVersion, "version", "", "App Store version (required, e.g. 1.2.3)")
	publishCmd.Flags().StringVar(&publishBuild, "build", "", "build ID for release attach/submit")
	publishCmd.Flags().StringVar(&publishMetadataDir, "metadata-dir", "", "path to ASC metadata version directory")
	publishCmd.Flags().StringVar(&publishPath, "path", ".", "project path for local preflight checks")
	publishCmd.Flags().StringVar(&publishIPA, "ipa", "", "optional .ipa path for binary preflight checks")
	publishCmd.Flags().IntVar(&publishScanTier, "scan-tier", 2, "max StoreReady ASC scan tier to run (1-4)")
	publishCmd.Flags().BoolVar(&publishSkipLocal, "skip-local-checks", false, "skip local StoreReady preflight checks")
	publishCmd.Flags().BoolVar(&publishSkipASCScan, "skip-asc-scan", false, "skip StoreReady ASC API checks before asc publish")
	publishCmd.Flags().BoolVar(&publishConfirm, "confirm", false, "run asc release for real (default is dry-run)")
	_ = publishCmd.MarkFlagRequired("app-id")
	_ = publishCmd.MarkFlagRequired("version")

	rootCmd.AddCommand(publishCmd)
}

func runPublish(cmd *cobra.Command, args []string) error {
	if publishScanTier < 1 || publishScanTier > 4 {
		return fmt.Errorf("--scan-tier must be between 1 and 4")
	}

	if _, err := exec.LookPath("asc"); err != nil {
		return fmt.Errorf("asc CLI not found in PATH. Install via: brew install asc")
	}

	purple.Println("\n  storeready publish — gate + release lane")
	fmt.Printf("  App ID:   %s\n", publishAppID)
	fmt.Printf("  Version:  %s\n", publishVersion)
	if publishBuild != "" {
		fmt.Printf("  Build:    %s\n", publishBuild)
	}
	if publishConfirm {
		fmt.Println("  Mode:     submit")
	} else {
		fmt.Println("  Mode:     dry-run")
	}
	fmt.Println()

	if !publishSkipLocal {
		dim.Println("  Running local StoreReady preflight gates...")
		start := time.Now()
		localResult, err := preflight.Run(publishPath, publishIPA, verbose)
		if err != nil {
			return fmt.Errorf("local preflight failed: %w", err)
		}
		if localResult.Summary.Critical > 0 {
			return fmt.Errorf("local preflight blocked release: %d critical finding(s). Fix them first or use --skip-local-checks", localResult.Summary.Critical)
		}
		dim.Printf("  ✓ Local preflight passed (%d findings, %s)\n", localResult.Summary.Total, time.Since(start).Round(time.Millisecond))
	}

	if !publishSkipASCScan {
		dim.Println("  Running StoreReady ASC metadata gates...")
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("ASC API credentials missing for scan gate: %w\nrun 'storeready auth setup' or pass --skip-asc-scan", err)
		}

		client, err := ascapi.NewClient(cfg.KeyID, cfg.IssuerID, cfg.PrivateKeyPath)
		if err != nil {
			return fmt.Errorf("failed to create ASC API client: %w", err)
		}

		start := time.Now()
		runner := checks.NewRunner(client, verbose)
		results, err := runner.Run(cmd.Context(), publishAppID, publishBuild, publishScanTier)
		if err != nil {
			return fmt.Errorf("ASC metadata gate failed: %w", err)
		}
		if results.Summary.Blocks > 0 {
			return fmt.Errorf("ASC scan blocked release: %d blocking finding(s). Fix them first or use --skip-asc-scan", results.Summary.Blocks)
		}
		dim.Printf("  ✓ ASC scan passed (%d findings, %s)\n", results.Summary.Total, time.Since(start).Round(time.Millisecond))
	}

	ascArgs := []string{"release", "run", "--app", publishAppID, "--version", publishVersion}
	if publishBuild != "" {
		ascArgs = append(ascArgs, "--build", publishBuild)
	}
	if publishMetadataDir != "" {
		ascArgs = append(ascArgs, "--metadata-dir", publishMetadataDir)
	}
	if publishConfirm {
		ascArgs = append(ascArgs, "--confirm")
	} else {
		ascArgs = append(ascArgs, "--dry-run")
	}

	dim.Printf("  Executing: asc %s\n\n", strings.Join(ascArgs, " "))

	ascCmd := exec.CommandContext(cmd.Context(), "asc", ascArgs...)
	ascCmd.Stdout = os.Stdout
	ascCmd.Stderr = os.Stderr
	ascCmd.Stdin = os.Stdin

	if err := ascCmd.Run(); err != nil {
		return fmt.Errorf("asc release failed: %w", err)
	}

	fmt.Println()
	purple.Println("  ✓ publish lane completed")
	fmt.Println()
	return nil
}
