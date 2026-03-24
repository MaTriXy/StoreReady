package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type releaseChecklistItem struct {
	ID           string `json:"id"`
	Area         string `json:"area"`
	Title        string `json:"title"`
	Verification string `json:"verification"` // MANUAL | HYBRID
	Why          string `json:"why"`
	Verify       string `json:"verify"`
	Source       string `json:"source,omitempty"`
}

type releaseChecklistSummary struct {
	Total  int `json:"total"`
	Manual int `json:"manual"`
	Hybrid int `json:"hybrid"`
}

type releaseChecklistOutput struct {
	GeneratedAt string                  `json:"generated_at"`
	AppType     string                  `json:"app_type"`
	Items       []releaseChecklistItem  `json:"items"`
	Summary     releaseChecklistSummary `json:"summary"`
}

var (
	releaseChecklistFormat     string
	releaseChecklistOutputPath string
	releaseChecklistAppType    string
)

var releaseChecklistCmd = &cobra.Command{
	Use:   "release-checklist",
	Short: "Structured manual gate list before App Store submit",
	Long: `Print a structured release checklist for items that are not fully automatable.
Run this before 'storeready publish --confirm' to avoid common App Review misses.`,
	RunE: runReleaseChecklist,
}

func init() {
	releaseChecklistCmd.Flags().StringVar(&releaseChecklistFormat, "format", "terminal", "output format: terminal, json")
	releaseChecklistCmd.Flags().StringVar(&releaseChecklistOutputPath, "output", "", "write checklist to file (stdout if omitted)")
	releaseChecklistCmd.Flags().StringVar(&releaseChecklistAppType, "app-type", "all", "app profile: all, subscription, social, kids, health, games, macos, ai, crypto, vpn")
	rootCmd.AddCommand(releaseChecklistCmd)
}

func runReleaseChecklist(cmd *cobra.Command, args []string) error {
	appType := strings.ToLower(strings.TrimSpace(releaseChecklistAppType))
	if appType == "" {
		appType = "all"
	}

	validAppTypes := map[string]bool{
		"all": true, "subscription": true, "social": true, "kids": true, "health": true,
		"games": true, "macos": true, "ai": true, "crypto": true, "vpn": true,
	}
	if !validAppTypes[appType] {
		return fmt.Errorf("invalid --app-type '%s' (use: all, subscription, social, kids, health, games, macos, ai, crypto, vpn)", appType)
	}

	items := append([]releaseChecklistItem{}, baseReleaseChecklistItems()...)
	items = append(items, appTypeSpecificChecklistItems(appType)...)

	out := releaseChecklistOutput{
		GeneratedAt: time.Now().Format(time.RFC3339),
		AppType:     appType,
		Items:       items,
		Summary:     buildReleaseChecklistSummary(items),
	}

	var output *os.File
	var err error
	if releaseChecklistOutputPath != "" {
		output, err = os.Create(releaseChecklistOutputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}

	switch strings.ToLower(releaseChecklistFormat) {
	case "json":
		enc := json.NewEncoder(output)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	default:
		return writeReleaseChecklistTerminal(output, out)
	}
}

func buildReleaseChecklistSummary(items []releaseChecklistItem) releaseChecklistSummary {
	s := releaseChecklistSummary{Total: len(items)}
	for _, item := range items {
		switch item.Verification {
		case "MANUAL":
			s.Manual++
		case "HYBRID":
			s.Hybrid++
		}
	}
	return s
}

func writeReleaseChecklistTerminal(w *os.File, out releaseChecklistOutput) error {
	bold := color.New(color.Bold)
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgHiBlue)

	purple.Println("\n  storeready release-checklist — manual gate review")
	fmt.Fprintf(w, "  App Type: %s\n", out.AppType)
	fmt.Fprintf(w, "  Items:    %d (%d manual, %d hybrid)\n\n", out.Summary.Total, out.Summary.Manual, out.Summary.Hybrid)

	for _, item := range out.Items {
		switch item.Verification {
		case "HYBRID":
			blue.Fprintf(w, "  [HYBRID] ")
		default:
			yellow.Fprintf(w, "  [MANUAL] ")
		}
		bold.Fprintf(w, "%s ", item.ID)
		fmt.Fprintf(w, "%s\n", item.Title)
		dim.Fprintf(w, "           Area: %s\n", item.Area)
		fmt.Fprintf(w, "           Why: %s\n", item.Why)
		fmt.Fprintf(w, "           Verify: %s\n", item.Verify)
		if item.Source != "" {
			dim.Fprintf(w, "           Source: %s\n", item.Source)
		}
		fmt.Fprintln(w)
	}

	dim.Println("  Run after this:")
	fmt.Fprintln(w, "  storeready publish --app-id <ID> --version <X.Y.Z> --build <BUILD_ID> --confirm")
	fmt.Fprintln(w)
	return nil
}

func baseReleaseChecklistItems() []releaseChecklistItem {
	return []releaseChecklistItem{
		{
			ID: "RC-001", Area: "Submission State", Title: "Version state is submittable",
			Verification: "HYBRID",
			Why:          "Review submission fails when version/build state is invalid.",
			Verify:       "In ASC, confirm version is PREPARE_FOR_SUBMISSION (or equivalent) and build is VALID.",
			Source:       "Apple ASC submission workflow",
		},
		{
			ID: "RC-002", Area: "Export Compliance", Title: "Encryption declaration is complete",
			Verification: "MANUAL",
			Why:          "Missing export-compliance answers blocks review.",
			Verify:       "Confirm ITSAppUsesNonExemptEncryption and export compliance declaration match app behavior.",
			Source:       "Guideline 5 / ASC export compliance",
		},
		{
			ID: "RC-003", Area: "Content Rights", Title: "Content rights declaration is set",
			Verification: "MANUAL",
			Why:          "Required before submit; missing declaration blocks processing.",
			Verify:       "Set DOES_NOT_USE_THIRD_PARTY_CONTENT or USES_THIRD_PARTY_CONTENT in ASC.",
			Source:       "ASC app information requirements",
		},
		{
			ID: "RC-004", Area: "Privacy", Title: "Privacy policy URL is valid for every locale",
			Verification: "HYBRID",
			Why:          "Missing/invalid privacy policy metadata commonly causes metadata rejection.",
			Verify:       "Check app-info localization URLs resolve publicly and match in-app privacy claims.",
			Source:       "Guideline 5.1.1",
		},
		{
			ID: "RC-005", Area: "Metadata", Title: "Store listing claims match in-app behavior",
			Verification: "MANUAL",
			Why:          "Overpromising or misleading metadata causes 2.3 rejections.",
			Verify:       "Review description, subtitle, keywords, and screenshots for truthful feature representation.",
			Source:       "Guideline 2.3",
		},
		{
			ID: "RC-006", Area: "Screenshots", Title: "All required screenshot sets are uploaded per locale",
			Verification: "HYBRID",
			Why:          "Missing storefront/device screenshot sets can block submission.",
			Verify:       "Confirm required device families and localizations are complete in ASC media manager.",
			Source:       "ASC media requirements",
		},
		{
			ID: "RC-007", Area: "Review Notes", Title: "Review notes include account/test credentials",
			Verification: "MANUAL",
			Why:          "Lack of reviewer access is a frequent avoidable rejection.",
			Verify:       "Provide test account, OTP bypass instructions, and any hardware/backend prerequisites.",
			Source:       "App Review operational best practice",
		},
		{
			ID: "RC-008", Area: "Network/Backend", Title: "Production backend endpoints are live",
			Verification: "MANUAL",
			Why:          "Broken endpoints during review lead to 2.1 completeness rejection.",
			Verify:       "Smoke-test sign-in, core flow, paywall and support URLs from a clean device/session.",
			Source:       "Guideline 2.1",
		},
		{
			ID: "RC-009", Area: "Account Management", Title: "Account deletion flow is accessible (if account creation exists)",
			Verification: "MANUAL",
			Why:          "Account apps must support deletion and communicate data deletion path.",
			Verify:       "Confirm in-app delete option, data handling behavior, and help documentation.",
			Source:       "Guideline 5.1.1",
		},
		{
			ID: "RC-010", Area: "Payments", Title: "Digital goods use IAP and restore flow exists",
			Verification: "MANUAL",
			Why:          "Using external payments for digital goods is a common rejection trigger.",
			Verify:       "Validate purchase/restore flows and remove external checkout references for digital content.",
			Source:       "Guideline 3.1.1",
		},
		{
			ID: "RC-011", Area: "Login", Title: "Sign in with Apple parity requirement is satisfied",
			Verification: "MANUAL",
			Why:          "If third-party social login exists, SIWA parity may be required.",
			Verify:       "Ensure SIWA is available where policy requires and does not request prohibited extra fields.",
			Source:       "Guideline 4.8",
		},
		{
			ID: "RC-012", Area: "Legal", Title: "Terms and policy links are present where required",
			Verification: "MANUAL",
			Why:          "Subscription/product pages often fail review due to missing legal links.",
			Verify:       "Ensure ToS/EULA/Privacy links exist in listing and in-app purchase context.",
			Source:       "Guideline 3.1.2 / 5.1.1",
		},
	}
}

func appTypeSpecificChecklistItems(appType string) []releaseChecklistItem {
	switch appType {
	case "subscription":
		return []releaseChecklistItem{
			{
				ID: "RC-SUB-001", Area: "Subscriptions", Title: "Pricing copy is not misleading",
				Verification: "MANUAL",
				Why:          "Misleading monthly-equivalent emphasis can trigger subscription metadata rejection.",
				Verify:       "Ensure billing period, trial terms, renewal terms, and effective price are clear and balanced.",
				Source:       "Guideline 3.1.2",
			},
		}
	case "social":
		return []releaseChecklistItem{
			{
				ID: "RC-SOC-001", Area: "UGC", Title: "User-generated content moderation controls are active",
				Verification: "MANUAL",
				Why:          "UGC apps require reporting/blocking/moderation to reduce abuse risk.",
				Verify:       "Verify report/block flows, moderation escalation, and abuse contact path.",
				Source:       "Guideline 1.2",
			},
		}
	case "kids":
		return []releaseChecklistItem{
			{
				ID: "RC-KID-001", Area: "Kids", Title: "Kids category privacy/ad requirements are met",
				Verification: "MANUAL",
				Why:          "Kids apps have stricter data collection and external-linking limits.",
				Verify:       "Confirm age-gating, no behavioral ads, and compliant third-party SDK usage.",
				Source:       "Guideline 1.3",
			},
		}
	case "health":
		return []releaseChecklistItem{
			{
				ID: "RC-HLT-001", Area: "Health", Title: "Medical/health claims are properly qualified",
				Verification: "MANUAL",
				Why:          "Unsubstantiated medical claims are high-risk during review.",
				Verify:       "Validate disclaimers, evidence language, and escalation path for critical conditions.",
				Source:       "Guideline 1.4",
			},
		}
	case "games":
		return []releaseChecklistItem{
			{
				ID: "RC-GME-001", Area: "Games", Title: "Loot box and randomization disclosures are present",
				Verification: "MANUAL",
				Why:          "Missing random-item odds disclosure may fail policy expectations.",
				Verify:       "Expose odds/conditions where required by platform and regional policy.",
				Source:       "Guideline 3.1.1 and regional policy",
			},
		}
	case "macos":
		return []releaseChecklistItem{
			{
				ID: "RC-MAC-001", Area: "macOS", Title: "Temporary exception entitlements are justified or removed",
				Verification: "MANUAL",
				Why:          "Unused/overbroad entitlements trigger rejection questions.",
				Verify:       "Audit entitlements and provide reviewer justification for each exception.",
				Source:       "Guideline 2.4.5(i)",
			},
		}
	case "ai":
		return []releaseChecklistItem{
			{
				ID: "RC-AI-001", Area: "AI", Title: "AI output safety and abuse handling are documented",
				Verification: "MANUAL",
				Why:          "Generative output without safeguards raises review and trust issues.",
				Verify:       "Confirm moderation, safety boundaries, and abuse reporting mechanisms are visible.",
				Source:       "Guideline 1.1 / 1.2",
			},
		}
	case "crypto":
		return []releaseChecklistItem{
			{
				ID: "RC-CRY-001", Area: "Crypto/Finance", Title: "Regulated functionality and regions are compliant",
				Verification: "MANUAL",
				Why:          "Finance/crypto features are subject to jurisdiction and licensing scrutiny.",
				Verify:       "Validate supported regions, disclosures, and licensing references for offered features.",
				Source:       "Guideline 3.1.5 and local law",
			},
		}
	case "vpn":
		return []releaseChecklistItem{
			{
				ID: "RC-VPN-001", Area: "VPN", Title: "VPN purpose and data handling are explicit",
				Verification: "MANUAL",
				Why:          "Networking/VPN apps receive heightened privacy and utility review.",
				Verify:       "Ensure logging policy, data retention, and core user benefit are clearly disclosed.",
				Source:       "Guideline 5.1.1 / Network Extension expectations",
			},
		}
	default:
		return nil
	}
}
