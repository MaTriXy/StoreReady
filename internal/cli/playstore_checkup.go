package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/MaTriXy/StoreReady/internal/playstore"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	playstoreFormat string
	playstoreOutput string
)

var playstoreCheckupCmd = &cobra.Command{
	Use:   "playstore-checkup [path]",
	Short: "Run Google Play pre-submission checks on your Android project",
	Long: `Run local Android checks to catch common Google Play submission risks.

Checks include:
  • AndroidManifest presence and risky flags (debuggable, cleartext, backups)
  • High-risk permission declarations that need policy justification
  • Gradle release metadata (applicationId, targetSdk, versionCode)
  • Basic packaging sanity checks`,
	Aliases: []string{"googleplay-checkup", "android-checkup"},
	Args:    cobra.MaximumNArgs(1),
	RunE:    runPlaystoreCheckup,
}

func init() {
	playstoreCheckupCmd.Flags().StringVar(&playstoreFormat, "format", "terminal", "output format: terminal, json")
	playstoreCheckupCmd.Flags().StringVar(&playstoreOutput, "output", "", "write report to file (stdout if omitted)")
	rootCmd.AddCommand(playstoreCheckupCmd)
}

func runPlaystoreCheckup(cmd *cobra.Command, args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot access path: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path must be a directory: %s", path)
	}

	purple.Println("\n  storeready playstore-checkup — Google Play pre-submission checks.")
	fmt.Printf("  Project: %s\n\n", path)

	start := time.Now()
	result, err := playstore.Run(path)
	if err != nil {
		return fmt.Errorf("playstore-checkup failed: %w", err)
	}
	result.Elapsed = time.Since(start)

	var output *os.File
	if playstoreOutput != "" {
		output, err = os.Create(playstoreOutput)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}

	switch strings.ToLower(playstoreFormat) {
	case "json":
		return writePlaystoreJSON(output, result)
	default:
		return writePlaystoreTerminal(output, result)
	}
}

func writePlaystoreTerminal(w *os.File, result *playstore.Result) error {
	red := color.New(color.FgRed, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen, color.Bold)
	bold := color.New(color.Bold)

	if result.ApplicationID != "" {
		fmt.Fprintf(w, "  App ID:   %s\n", result.ApplicationID)
	}
	if result.TargetSDK != "" {
		fmt.Fprintf(w, "  Target:   SDK %s\n", result.TargetSDK)
	}
	if result.MinSDK != "" {
		fmt.Fprintf(w, "  Min SDK:  %s\n", result.MinSDK)
	}
	if result.VersionCode != "" {
		fmt.Fprintf(w, "  Version:  code %s\n", result.VersionCode)
	}
	if result.ManifestPath != "" {
		fmt.Fprintf(w, "  Manifest: %s\n", result.ManifestPath)
	}
	if result.Coverage.Total > 0 {
		fmt.Fprintf(
			w,
			"  Coverage: %d controls (%d automated, %d hybrid, %d manual)\n",
			result.Coverage.Total,
			result.Coverage.Automated,
			result.Coverage.Hybrid,
			result.Coverage.Manual,
		)
	}
	fmt.Fprintln(w)

	if len(result.Findings) == 0 {
		green.Fprintln(w, "  No issues found!")
		fmt.Fprintln(w)
		printPlaystoreFooter(w, result)
		printPlaystoreChecklist(w, result)
		return nil
	}

	sort.Slice(result.Findings, func(i, j int) bool {
		sevRank := map[string]int{"CRITICAL": 3, "WARN": 2, "INFO": 1}
		ri := sevRank[result.Findings[i].Severity]
		rj := sevRank[result.Findings[j].Severity]
		if ri != rj {
			return ri > rj
		}
		return result.Findings[i].File < result.Findings[j].File
	})

	var criticals, warns, infos []playstore.Finding
	for _, finding := range result.Findings {
		switch finding.Severity {
		case "CRITICAL":
			criticals = append(criticals, finding)
		case "WARN":
			warns = append(warns, finding)
		case "INFO":
			infos = append(infos, finding)
		}
	}

	if len(criticals) > 0 {
		red.Fprintln(w, "  CRITICAL — Blocking release risks")
		fmt.Fprintln(w)
		for _, finding := range criticals {
			printPlaystoreFinding(w, finding, bold)
		}
	}

	if len(warns) > 0 {
		yellow.Fprintln(w, "  WARNING — Policy/security risks")
		fmt.Fprintln(w)
		for _, finding := range warns {
			printPlaystoreFinding(w, finding, bold)
		}
	}

	if len(infos) > 0 {
		dim.Fprintln(w, "  INFO — Best practices")
		fmt.Fprintln(w)
		for _, finding := range infos {
			printPlaystoreFinding(w, finding, bold)
		}
	}

	printPlaystoreFooter(w, result)
	printPlaystoreChecklist(w, result)
	return nil
}

func printPlaystoreFinding(w *os.File, finding playstore.Finding, bold *color.Color) {
	green := color.New(color.FgGreen)
	switch finding.Severity {
	case "CRITICAL":
		color.New(color.FgRed, color.Bold).Fprintf(w, "  [CRITICAL] ")
	case "WARN":
		color.New(color.FgYellow).Fprintf(w, "  [WARN]     ")
	default:
		dim.Fprintf(w, "  [INFO]     ")
	}

	if finding.Guideline != "" {
		bold.Fprintf(w, "[%s] ", finding.Guideline)
	}
	bold.Fprintln(w, finding.Title)

	if finding.File != "" {
		dim.Fprintf(w, "             %s\n", finding.File)
	}

	fmt.Fprintf(w, "             %s\n", finding.Detail)
	if finding.Fix != "" {
		green.Fprintf(w, "             Fix: ")
		fmt.Fprintln(w, finding.Fix)
	}
	fmt.Fprintln(w)
}

func printPlaystoreFooter(w *os.File, result *playstore.Result) {
	red := color.New(color.FgRed, color.Bold)
	green := color.New(color.FgGreen, color.Bold)

	fmt.Fprintln(w)
	dim.Fprintln(w, "  ─────────────────────────────────────────────")
	fmt.Fprintln(w)

	if result.Summary.Passed {
		green.Fprint(w, "  READY")
		fmt.Fprint(w, " — no critical issues found")
	} else {
		red.Fprint(w, "  NOT READY")
		fmt.Fprintf(w, " — %d critical issue(s) must be fixed", result.Summary.Critical)
	}
	fmt.Fprintln(w)

	if result.Summary.Total > 0 {
		fmt.Fprintf(w, "  %d findings: ", result.Summary.Total)
		if result.Summary.Critical > 0 {
			red.Fprintf(w, "%d critical  ", result.Summary.Critical)
		}
		if result.Summary.Warns > 0 {
			color.New(color.FgYellow).Fprintf(w, "%d warn  ", result.Summary.Warns)
		}
		if result.Summary.Infos > 0 {
			dim.Fprintf(w, "%d info", result.Summary.Infos)
		}
		fmt.Fprintln(w)
	}

	dim.Fprintf(w, "  completed in %s\n", result.Elapsed.Round(time.Millisecond))
	fmt.Fprintln(w)
}

func printPlaystoreChecklist(w *os.File, result *playstore.Result) {
	if len(result.Checklist) == 0 {
		return
	}

	var failing []playstore.ChecklistItem
	var warning []playstore.ChecklistItem
	var manual []playstore.ChecklistItem
	for _, item := range result.Checklist {
		switch item.Status {
		case "fail":
			failing = append(failing, item)
		case "warning":
			warning = append(warning, item)
		case "needs_manual_review":
			manual = append(manual, item)
		}
	}

	if len(failing) == 0 && len(warning) == 0 && len(manual) == 0 {
		return
	}

	color.New(color.Bold).Fprintln(w, "  Policy Checklist Review")
	fmt.Fprintln(w)

	for _, item := range failing {
		color.New(color.FgRed, color.Bold).Fprintf(w, "  [FAIL] ")
		fmt.Fprintf(w, "%s %s\n", item.Section, item.Title)
		if item.Notes != "" {
			dim.Fprintf(w, "         %s\n", item.Notes)
		}
	}

	for _, item := range warning {
		color.New(color.FgYellow).Fprintf(w, "  [WARN] ")
		fmt.Fprintf(w, "%s %s\n", item.Section, item.Title)
		if item.Notes != "" {
			dim.Fprintf(w, "         %s\n", item.Notes)
		}
	}

	for _, item := range manual {
		dim.Fprintf(w, "  [MANUAL] %s %s\n", item.Section, item.Title)
	}

	fmt.Fprintln(w)
	dim.Fprintln(w, "  Tip: run `storeready play-guidelines show <SECTION>` for the exact review checklist and source policy links.")
	fmt.Fprintln(w)
}

func writePlaystoreJSON(w *os.File, result *playstore.Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}
