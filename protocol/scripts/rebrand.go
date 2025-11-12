// scripts/rebrand.go
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	// Customize according to your repo
	replacements = map[string]string{
		// bech32 hrp & variants
		"dydxvaloperpub": "vindaxvaloperpub",
		"dydxvalconspub": "vindaxvalconspub",
		"dydxvaloper":    "vindaxvaloper",
		"dydxvalcons":    "vindaxvalcons",
		"dydxpub":        "vindaxpub",
		"dydx1":          "vindax1",

		// message type URLs
		"/dydxprotocol.": "/vindax.",

		// proto package names and file paths
		"dydxprotocol.": "vindax.",
		"dydxprotocol/": "vindax/",
		// "dydxprotocol":  "vindax", // for gRPC gateway patterns and standalone strings

		// denom, chain-id, binary
		// "avdtn":        "avdtn",
		// "dydx-mainnet-1": "vindax-1",
		// "vindaxd":  "vindaxd",
		// "dydxd":          "vindaxd",

		// Go module path
		// "github.com/dydxprotocol/v4-chain": "https://github.com/danielvindax/vd-chain",
	}

	// Skip directories that should not be touched
	skipDirs = map[string]bool{
		".git": true, "vendor": true, "node_modules": true, "build": true, "out": true,
		"_backup": true, "_backup_rebrand": true, ".cache": true, ".idea": true, ".vscode": true,
		"third_party": true, "thirdparty": true, "proto/google": true,
	}

	// Skip this tool file itself (self-exclude)
	skipFiles = map[string]bool{
		"scripts/rebrand.go": true, // run from repo root
		"rebrand.go":         true, // in case running in the scripts directory itself
	}

	// Skip generated files
	skipGenerated = map[string]bool{
		".pb.go":    true,
		".pb.gw.go": true,
		"go.mod":    true,
		"go.sum":    true,
	}
)

func isTextFile(path string) bool {
	switch {
	case strings.HasSuffix(path, ".go"),
		strings.HasSuffix(path, ".json"),
		strings.HasSuffix(path, ".yaml"),
		strings.HasSuffix(path, ".yml"),
		strings.HasSuffix(path, ".sh"),
		strings.HasSuffix(path, ".md"),
		strings.HasSuffix(path, ".toml"),
		strings.HasSuffix(path, "Makefile"):
		return true
	default:
		return false
	}
}

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	var patched int

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if skipDirs[filepath.Base(path)] {
				return filepath.SkipDir
			}
			return nil
		}
		rel := path
		if r, e := filepath.Rel(".", path); e == nil {
			rel = filepath.ToSlash(r)
		}
		if skipFiles[rel] || skipFiles[filepath.Base(rel)] {
			return nil
		}

		if !isTextFile(path) {
			return nil
		}

		// Skip generated files
		for suffix := range skipGenerated {
			if strings.HasSuffix(path, suffix) {
				return nil
			}
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		orig := string(data)
		mod := orig

		// Apply replacements in order from longest to shortest to avoid conflicts
		// (e.g., "dydxprotocol." should be replaced before "dydxprotocol")
		type replacement struct {
			old string
			new string
		}
		var sortedReplacements []replacement
		for old, newv := range replacements {
			sortedReplacements = append(sortedReplacements, replacement{old, newv})
		}
		// Sort by length descending
		for i := 0; i < len(sortedReplacements); i++ {
			for j := i + 1; j < len(sortedReplacements); j++ {
				if len(sortedReplacements[i].old) < len(sortedReplacements[j].old) {
					sortedReplacements[i], sortedReplacements[j] = sortedReplacements[j], sortedReplacements[i]
				}
			}
		}

		// Apply replacements, but skip import lines in Go files
		if strings.HasSuffix(path, ".go") {
			lines := strings.Split(orig, "\n")
			var result []string
			inImportBlock := false

			for _, line := range lines {
				trimmed := strings.TrimSpace(line)

				// Check for single-line import: import "package" or import _ "package"
				if strings.HasPrefix(trimmed, "import ") {
					// Check if it's a single-line import (no opening paren on same line)
					if !strings.Contains(trimmed, "(") || strings.HasSuffix(trimmed, ")") {
						result = append(result, line) // Keep original, don't replace
						continue
					}
				}

				// Check for start of import block: import (
				if strings.HasPrefix(trimmed, "import (") {
					inImportBlock = true
					result = append(result, line) // Keep original
					continue
				}

				// Check if we're inside import block
				if inImportBlock {
					result = append(result, line) // Keep original, don't replace
					// Check if this line closes the import block
					if strings.TrimSpace(line) == ")" {
						inImportBlock = false
					}
					continue
				}

				// Apply replacements to non-import lines
				modifiedLine := line
				for _, r := range sortedReplacements {
					modifiedLine = strings.ReplaceAll(modifiedLine, r.old, r.new)
				}
				result = append(result, modifiedLine)
			}
			mod = strings.Join(result, "\n")
		} else {
			// For non-Go files, apply replacements normally
			for _, r := range sortedReplacements {
				mod = strings.ReplaceAll(mod, r.old, r.new)
			}
		}
		if mod != orig {
			if err := os.WriteFile(path, []byte(mod), 0o644); err != nil {
				return err
			}
			patched++
			fmt.Println("🪄 patched:", path)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "rebrand error:", err)
		os.Exit(1)
	}
	fmt.Printf("✅ rebrand replace done. files patched: %d\n", patched)
}
