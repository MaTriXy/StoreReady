package cli

import (
	"fmt"
	"strings"

	"github.com/MaTriXy/StoreReady/internal/playguidelines"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var playGuidelinesCmd = &cobra.Command{
	Use:     "play-guidelines",
	Aliases: []string{"pg"},
	Short:   "Browse and search Google Play policy guidelines",
}

var playGuidelinesSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search Google Play guidelines by keyword",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runPlayGuidelinesSearch,
}

var playGuidelinesShowCmd = &cobra.Command{
	Use:   "show [section]",
	Short: "Show a specific guideline section (e.g. 'GP-3', 'GP-1.2')",
	Args:  cobra.ExactArgs(1),
	RunE:  runPlayGuidelinesShow,
}

var playGuidelinesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all top-level Google Play guideline sections",
	RunE:  runPlayGuidelinesList,
}

func init() {
	playGuidelinesCmd.AddCommand(playGuidelinesSearchCmd)
	playGuidelinesCmd.AddCommand(playGuidelinesShowCmd)
	playGuidelinesCmd.AddCommand(playGuidelinesListCmd)
	rootCmd.AddCommand(playGuidelinesCmd)
}

func runPlayGuidelinesSearch(cmd *cobra.Command, args []string) error {
	query := strings.Join(args, " ")
	db, err := playguidelines.Load()
	if err != nil {
		return fmt.Errorf("failed to load Play guidelines: %w", err)
	}

	results := db.Search(query)
	purple.Printf("\n  Google Play guidelines matching '%s'\n\n", query)

	if len(results) == 0 {
		dim.Println("  No matching guidelines found.")
		return nil
	}

	for _, g := range results {
		bold := color.New(color.Bold)
		bold.Printf("  %s  ", g.Section)
		fmt.Printf("%s [%s]\n", g.Title, strings.ToUpper(g.Verification))
		dim.Printf("  %s\n\n", truncate(g.Content, 120))
	}

	return nil
}

func runPlayGuidelinesShow(cmd *cobra.Command, args []string) error {
	section := args[0]
	db, err := playguidelines.Load()
	if err != nil {
		return fmt.Errorf("failed to load Play guidelines: %w", err)
	}

	g, found := db.Get(section)
	if !found {
		return fmt.Errorf("guideline section '%s' not found", section)
	}

	purple.Printf("\n  Google Play Guideline %s\n", g.Section)
	color.New(color.Bold).Printf("  %s\n\n", g.Title)
	fmt.Printf("  Verification: %s\n\n", strings.ToUpper(g.Verification))
	fmt.Printf("  %s\n", g.Content)

	if len(g.AutomatedChecks) > 0 {
		fmt.Println()
		color.New(color.FgHiBlue).Println("  Automated checks:")
		for _, check := range g.AutomatedChecks {
			fmt.Printf("    • %s\n", check)
		}
	}

	if len(g.ManualChecks) > 0 {
		fmt.Println()
		color.New(color.FgYellow).Println("  Manual checks:")
		for _, check := range g.ManualChecks {
			fmt.Printf("    • %s\n", check)
		}
	}

	if len(g.Sources) > 0 {
		fmt.Println()
		dim.Println("  Sources:")
		for _, source := range g.Sources {
			fmt.Printf("    %s\n", source)
		}
	}

	if len(g.Subsections) > 0 {
		fmt.Println()
		dim.Println("  Subsections:")
		for _, s := range g.Subsections {
			fmt.Printf("    %s  %s\n", s.Section, s.Title)
		}
	}

	fmt.Println()
	return nil
}

func runPlayGuidelinesList(cmd *cobra.Command, args []string) error {
	db, err := playguidelines.Load()
	if err != nil {
		return fmt.Errorf("failed to load Play guidelines: %w", err)
	}

	purple.Println("\n  Google Play Policy Guidelines Matrix")
	for _, g := range db.TopLevel() {
		bold := color.New(color.Bold)
		bold.Printf("  %s  ", g.Section)
		fmt.Printf("%s [%s]\n", g.Title, strings.ToUpper(g.Verification))
		for _, s := range g.Subsections {
			dim.Printf("      %s  %s [%s]\n", s.Section, s.Title, strings.ToUpper(s.Verification))
		}
		fmt.Println()
	}

	return nil
}
