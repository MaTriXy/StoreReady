package playstore

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/MaTriXy/StoreReady/internal/playguidelines"
)

type Finding struct {
	ID        string `json:"id,omitempty"`
	Severity  string `json:"severity"` // "CRITICAL", "WARN", "INFO"
	Guideline string `json:"guideline,omitempty"`
	Title     string `json:"title"`
	Detail    string `json:"detail"`
	Fix       string `json:"fix,omitempty"`
	File      string `json:"file,omitempty"`
}

type ChecklistItem struct {
	Section      string   `json:"section"`
	Title        string   `json:"title"`
	Verification string   `json:"verification"` // automated, hybrid, manual
	Status       string   `json:"status"`       // pass, warning, fail, needs_manual_review
	Notes        string   `json:"notes,omitempty"`
	Sources      []string `json:"sources,omitempty"`
}

type Coverage struct {
	Total          int `json:"total"`
	Automated      int `json:"automated"`
	Hybrid         int `json:"hybrid"`
	Manual         int `json:"manual"`
	Passed         int `json:"passed"`
	Warnings       int `json:"warnings"`
	Failed         int `json:"failed"`
	ManualRequired int `json:"manual_required"`
}

type Summary struct {
	Total    int  `json:"total"`
	Critical int  `json:"critical"`
	Warns    int  `json:"warns"`
	Infos    int  `json:"infos"`
	Passed   bool `json:"passed"`
}

type Result struct {
	ProjectPath   string          `json:"project_path"`
	Findings      []Finding       `json:"findings"`
	Checklist     []ChecklistItem `json:"checklist,omitempty"`
	Coverage      Coverage        `json:"coverage"`
	Summary       Summary         `json:"summary"`
	Elapsed       time.Duration   `json:"elapsed"`
	ApplicationID string          `json:"application_id,omitempty"`
	TargetSDK     string          `json:"target_sdk,omitempty"`
	MinSDK        string          `json:"min_sdk,omitempty"`
	VersionCode   string          `json:"version_code,omitempty"`
	ManifestPath  string          `json:"manifest_path,omitempty"`
}

func Run(projectPath string) (*Result, error) {
	result := &Result{ProjectPath: projectPath}

	if _, err := os.Stat(projectPath); err != nil {
		return nil, fmt.Errorf("cannot access path: %w", err)
	}

	manifestFiles, err := findFiles(projectPath, "AndroidManifest.xml")
	if err != nil {
		return nil, err
	}

	if len(manifestFiles) == 0 {
		result.Findings = append(result.Findings, Finding{
			ID:        "manifest-missing",
			Severity:  "CRITICAL",
			Guideline: "Manifest",
			Title:     "No AndroidManifest.xml found",
			Detail:    "A Play Store Android app needs at least one AndroidManifest.xml.",
			Fix:       "Ensure your Android project is present and includes app/src/main/AndroidManifest.xml.",
		})
		result.Checklist, result.Coverage = buildChecklist(result.Findings)
		result.Summary = computeSummary(result.Findings)
		return result, nil
	}

	for _, manifestPath := range manifestFiles {
		content, readErr := os.ReadFile(manifestPath)
		if readErr != nil {
			continue
		}
		analyzeManifest(string(content), manifestPath, result)
	}

	gradleFiles, err := findGradleFiles(projectPath)
	if err != nil {
		return nil, err
	}
	for _, gradlePath := range gradleFiles {
		content, readErr := os.ReadFile(gradlePath)
		if readErr != nil {
			continue
		}
		analyzeGradle(string(content), gradlePath, result)
	}

	if result.TargetSDK == "" {
		result.Findings = append(result.Findings, Finding{
			ID:        "target-sdk-missing",
			Severity:  "WARN",
			Guideline: "Target SDK",
			Title:     "Target SDK not detected",
			Detail:    "Could not find targetSdk/targetSdkVersion in Gradle config.",
			Fix:       "Set targetSdk in your module build.gradle/build.gradle.kts and keep it updated for Play requirements.",
		})
	}

	if result.ApplicationID == "" {
		result.Findings = append(result.Findings, Finding{
			ID:        "appid-missing",
			Severity:  "WARN",
			Guideline: "Package ID",
			Title:     "Application ID not detected",
			Detail:    "Could not find applicationId in Gradle config or package in AndroidManifest.xml.",
			Fix:       "Set a stable applicationId (for example: com.company.app) in your app module Gradle file.",
		})
	}

	if result.VersionCode == "" {
		result.Findings = append(result.Findings, Finding{
			ID:        "versioncode-missing",
			Severity:  "WARN",
			Guideline: "Versioning",
			Title:     "versionCode not detected",
			Detail:    "Could not find versionCode in Gradle config.",
			Fix:       "Set an integer versionCode in your app module to ensure Play Store updates can be uploaded.",
		})
	}

	result.Findings = dedup(result.Findings)
	result.Checklist, result.Coverage = buildChecklist(result.Findings)
	result.Summary = computeSummary(result.Findings)
	return result, nil
}

func analyzeManifest(content, path string, result *Result) {
	if result.ManifestPath == "" {
		result.ManifestPath = path
	}

	manifestPackageRe := regexp.MustCompile(`<manifest[^>]*\spackage="([^"]+)"`)
	if m := manifestPackageRe.FindStringSubmatch(content); len(m) == 2 && result.ApplicationID == "" {
		result.ApplicationID = m[1]
	}

	if strings.Contains(content, `android:debuggable="true"`) {
		result.Findings = append(result.Findings, Finding{
			ID:        "debuggable-true",
			Severity:  "CRITICAL",
			Guideline: "Security",
			Title:     "Release app appears debuggable",
			Detail:    "android:debuggable=\"true\" can expose internals and is unsafe for production Play builds.",
			Fix:       "Disable debuggable in release builds and verify release manifest merge output.",
			File:      path,
		})
	}

	if strings.Contains(content, `android:usesCleartextTraffic="true"`) {
		result.Findings = append(result.Findings, Finding{
			ID:        "cleartext-true",
			Severity:  "WARN",
			Guideline: "Network Security",
			Title:     "Cleartext network traffic is enabled",
			Detail:    "Allowing HTTP traffic increases MITM risk and may violate security expectations.",
			Fix:       "Use HTTPS by default and remove usesCleartextTraffic unless strictly necessary.",
			File:      path,
		})
	}

	if strings.Contains(content, `android:allowBackup="true"`) {
		result.Findings = append(result.Findings, Finding{
			ID:        "allow-backup-true",
			Severity:  "WARN",
			Guideline: "Data Safety",
			Title:     "Application backups are enabled",
			Detail:    "android:allowBackup=\"true\" may expose user data via backups if not intentionally designed.",
			Fix:       "Set android:allowBackup=\"false\" for sensitive apps, or document backup behavior in your data policy.",
			File:      path,
		})
	}

	permissionRe := regexp.MustCompile(`<uses-permission[^>]*android:name="([^"]+)"`)
	matches := permissionRe.FindAllStringSubmatch(content, -1)
	highRisk := map[string]string{
		"android.permission.READ_SMS":                "SMS/Call Log",
		"android.permission.RECEIVE_SMS":             "SMS/Call Log",
		"android.permission.READ_CALL_LOG":           "SMS/Call Log",
		"android.permission.WRITE_CALL_LOG":          "SMS/Call Log",
		"android.permission.QUERY_ALL_PACKAGES":      "Data Minimization",
		"android.permission.MANAGE_EXTERNAL_STORAGE": "All Files Access",
	}
	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		perm := match[1]
		if guideline, ok := highRisk[perm]; ok {
			result.Findings = append(result.Findings, Finding{
				ID:        "high-risk-permission:" + perm,
				Severity:  "WARN",
				Guideline: guideline,
				Title:     fmt.Sprintf("High-risk permission declared: %s", perm),
				Detail:    "Google Play often requires strong justification and policy declarations for this permission.",
				Fix:       "Verify this permission is strictly required, update Play Console declarations, and provide a clear in-app user benefit.",
				File:      path,
			})
		}
	}

	if result.ApplicationID != "" && strings.HasPrefix(result.ApplicationID, "com.example.") {
		result.Findings = append(result.Findings, Finding{
			ID:        "placeholder-application-id",
			Severity:  "WARN",
			Guideline: "Package ID",
			Title:     "Placeholder application ID detected",
			Detail:    "Package names starting with com.example are typically template values and not suitable for release.",
			Fix:       "Set a stable production package ID (for example: com.company.product).",
			File:      path,
		})
	}
}

func analyzeGradle(content, path string, result *Result) {
	appIDRe := regexp.MustCompile(`(?m)applicationId\s*[= ]\s*["']([^"']+)["']`)
	if m := appIDRe.FindStringSubmatch(content); len(m) == 2 && result.ApplicationID == "" {
		result.ApplicationID = m[1]
	}

	targetSDKRe := regexp.MustCompile(`(?m)targetSdk(?:Version)?\s*[= ]\s*([0-9]+)`)
	if m := targetSDKRe.FindStringSubmatch(content); len(m) == 2 && result.TargetSDK == "" {
		result.TargetSDK = m[1]
	}

	minSDKRe := regexp.MustCompile(`(?m)minSdk(?:Version)?\s*[= ]\s*([0-9]+)`)
	if m := minSDKRe.FindStringSubmatch(content); len(m) == 2 && result.MinSDK == "" {
		result.MinSDK = m[1]
	}

	versionCodeRe := regexp.MustCompile(`(?m)versionCode\s*[= ]\s*([0-9]+)`)
	if m := versionCodeRe.FindStringSubmatch(content); len(m) == 2 && result.VersionCode == "" {
		result.VersionCode = m[1]
	}

	if m := targetSDKRe.FindStringSubmatch(content); len(m) == 2 {
		target, convErr := strconv.Atoi(m[1])
		if convErr == nil && target < 34 {
			result.Findings = append(result.Findings, Finding{
				ID:        "target-sdk-outdated",
				Severity:  "WARN",
				Guideline: "Target SDK",
				Title:     fmt.Sprintf("Target SDK appears old (%d)", target),
				Detail:    "Older target SDK values can block Play uploads as policy requirements advance.",
				Fix:       "Update compileSdk/targetSdk to the latest Play-required API level and retest behavior changes.",
				File:      path,
			})
		}
	}
}

func computeSummary(findings []Finding) Summary {
	s := Summary{}
	for _, finding := range findings {
		s.Total++
		switch finding.Severity {
		case "CRITICAL":
			s.Critical++
		case "WARN":
			s.Warns++
		case "INFO":
			s.Infos++
		}
	}
	s.Passed = s.Critical == 0
	return s
}

func dedup(findings []Finding) []Finding {
	seen := make(map[string]int)
	sevRank := map[string]int{"CRITICAL": 3, "WARN": 2, "INFO": 1}
	var result []Finding

	for _, finding := range findings {
		key := finding.Title + "|" + finding.File
		if idx, ok := seen[key]; ok {
			if sevRank[finding.Severity] > sevRank[result[idx].Severity] {
				result[idx] = finding
			}
			continue
		}
		seen[key] = len(result)
		result = append(result, finding)
	}
	return result
}

func findFiles(root, name string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && skipDir(d.Name()) {
			return filepath.SkipDir
		}
		if !d.IsDir() && d.Name() == name {
			out = append(out, path)
		}
		return nil
	})
	return out, err
}

func findGradleFiles(root string) ([]string, error) {
	var out []string
	allowed := []string{"build.gradle", "build.gradle.kts"}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && skipDir(d.Name()) {
			return filepath.SkipDir
		}
		if !d.IsDir() && slices.Contains(allowed, d.Name()) {
			out = append(out, path)
		}
		return nil
	})
	return out, err
}

func skipDir(name string) bool {
	switch name {
	case ".git", "node_modules", "Pods", "build", ".next", ".gradle":
		return true
	default:
		return false
	}
}

func buildChecklist(findings []Finding) ([]ChecklistItem, Coverage) {
	db, err := playguidelines.Load()
	if err != nil {
		return nil, Coverage{}
	}

	all := db.Flatten()
	findingsByID := make(map[string][]Finding)
	for _, finding := range findings {
		findingsByID[finding.ID] = append(findingsByID[finding.ID], finding)
	}

	var items []ChecklistItem
	coverage := Coverage{}

	for _, guideline := range all {
		if len(guideline.Subsections) > 0 {
			continue
		}

		item := ChecklistItem{
			Section:      guideline.Section,
			Title:        guideline.Title,
			Verification: guideline.Verification,
			Status:       "needs_manual_review",
			Sources:      guideline.Sources,
		}

		switch guideline.Verification {
		case "automated":
			coverage.Automated++
		case "hybrid":
			coverage.Hybrid++
		default:
			coverage.Manual++
		}

		var matched []Finding
		for _, check := range guideline.AutomatedChecks {
			matched = append(matched, matchingFindings(check, findingsByID)...)
		}

		highest := highestSeverity(matched)
		switch guideline.Verification {
		case "automated":
			switch highest {
			case "CRITICAL":
				item.Status = "fail"
				item.Notes = fmt.Sprintf("%d automated finding(s), includes critical issues.", len(matched))
			case "WARN":
				item.Status = "warning"
				item.Notes = fmt.Sprintf("%d automated warning(s) found.", len(matched))
			default:
				item.Status = "pass"
				item.Notes = "Automated checks passed."
			}
		case "hybrid":
			switch highest {
			case "CRITICAL":
				item.Status = "fail"
				item.Notes = "Automated checks found critical issues and manual review is still required."
			case "WARN":
				item.Status = "warning"
				item.Notes = "Automated checks found warnings. Manual policy declarations are required."
			default:
				item.Status = "needs_manual_review"
				item.Notes = "Automated checks passed. Manual Play Console review is still required."
			}
		default:
			item.Status = "needs_manual_review"
			item.Notes = "Manual policy review required."
		}

		switch item.Status {
		case "pass":
			coverage.Passed++
		case "warning":
			coverage.Warnings++
		case "fail":
			coverage.Failed++
		default:
			coverage.ManualRequired++
		}

		items = append(items, item)
	}
	coverage.Total = len(items)

	return items, coverage
}

func matchingFindings(checkPattern string, findingsByID map[string][]Finding) []Finding {
	if strings.HasSuffix(checkPattern, "*") {
		prefix := strings.TrimSuffix(checkPattern, "*")
		var out []Finding
		for findingID, matched := range findingsByID {
			if strings.HasPrefix(findingID, prefix) {
				out = append(out, matched...)
			}
		}
		return out
	}
	return findingsByID[checkPattern]
}

func highestSeverity(findings []Finding) string {
	highest := ""
	for _, finding := range findings {
		if sevRank(finding.Severity) > sevRank(highest) {
			highest = finding.Severity
		}
	}
	return highest
}

func sevRank(severity string) int {
	switch severity {
	case "CRITICAL":
		return 3
	case "WARN":
		return 2
	case "INFO":
		return 1
	default:
		return 0
	}
}
