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
	// Tùy chỉnh theo repo của bạn
	replacements = map[string]string{
		// bech32 hrp & biến thể
		"dydxvaloperpub": "vindaxvaloperpub",
		"dydxvalconspub": "vindaxvalconspub",
		"dydxvaloper":    "vindaxvaloper",
		"dydxvalcons":    "vindaxvalcons",
		"dydxpub":        "vindaxpub",
		"dydx1":          "vindax1",

		// denom, chain-id, binary
		// "adv4tnt":        "uVNDX",
		// "dydx-mainnet-1": "vindax-1",
		// "dydxprotocold":  "vindaxd",
		// "dydxd":          "vindaxd",

		// Go module path
		// "github.com/dydxprotocol/v4-chain": "github.com/miexs/vd-chain",
	}

	// Bỏ qua các thư mục không nên chạm
	skipDirs = map[string]bool{
		".git": true, "vendor": true, "node_modules": true, "build": true, "out": true,
		"_backup": true, "_backup_rebrand": true, ".cache": true, ".idea": true, ".vscode": true,
		"third_party": true, "thirdparty": true, "proto/google": true,
	}

	// Bỏ qua chính file tool này (self-exclude)
	skipFiles = map[string]bool{
		"scripts/rebrand.go": true, // chạy từ repo root
		"rebrand.go":         true, // phòng trường hợp chạy trong chính thư mục scripts
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
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		orig := string(data)
		mod := orig
		for old, newv := range replacements {
			mod = strings.ReplaceAll(mod, old, newv)
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
