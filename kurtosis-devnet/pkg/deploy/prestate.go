package deploy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tokamak-network/tokamak-thanos/kurtosis-devnet/pkg/build"
)

type PrestateInfo struct {
	URL    string            `json:"url"`
	Hashes map[string]string `json:"hashes"`
}

type localPrestateHolder struct {
	info       *PrestateInfo
	baseDir    string
	buildDir   string
	dryRun     bool
	builder    *build.PrestateBuilder
	urlBuilder func(path ...string) string
}

func (h *localPrestateHolder) GetPrestateInfo(ctx context.Context) (*PrestateInfo, error) {
	if h.info != nil {
		return h.info, nil
	}

	prestatePath := []string{"proofs", "op-program", "cannon"}
	prestateURL := h.urlBuilder(prestatePath...)

	// Create build directory with the final path structure
	buildDir := filepath.Join(append([]string{h.buildDir}, prestatePath...)...)
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create prestate build directory: %w", err)
	}

	info := &PrestateInfo{
		URL:    prestateURL,
		Hashes: make(map[string]string),
	}

	if h.dryRun {
		// In dry run, populate with placeholder keys to avoid template errors during first pass
		info.Hashes["prestate_mt64"] = "dry_run_placeholder"
		info.Hashes["prestate_interop"] = "dry_run_placeholder"
		h.info = info
		return info, nil
	}

	// Map of known file prefixes to their keys
	fileToKey := map[string]string{
		"prestate-proof-mt64.json":    "prestate_mt64",
		"prestate-proof-interop.json": "prestate_interop",
	}

	// Build all prestate files directly in the target directory
	if err := h.builder.Build(ctx, buildDir); err != nil {
		return nil, fmt.Errorf("failed to build prestates: %w", err)
	}

	// Find and process all prestate files
	matches, err := filepath.Glob(filepath.Join(buildDir, "prestate-proof-mt64.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to find prestate files: %w", err)
	}

	// Create prestates directory for challenger compatibility
	prestatesDir := filepath.Join(buildDir, "prestates")
	if err := os.MkdirAll(prestatesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create prestates directory: %w", err)
	}

	var prestateProofCopied bool

	// Process each file to rename it to its hash
	for _, filePath := range matches {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read prestate %s: %w", filepath.Base(filePath), err)
		}

		var data struct {
			Pre string `json:"pre"`
		}
		if err := json.Unmarshal(content, &data); err != nil {
			return nil, fmt.Errorf("failed to parse prestate %s: %w", filepath.Base(filePath), err)
		}

		// Store hash with its corresponding key
		if key, exists := fileToKey[filepath.Base(filePath)]; exists {
			info.Hashes[key] = data.Pre
		}

		// Copy first prestate file to prestates/prestate-proof.json for challenger compatibility
		if !prestateProofCopied {
			prestateProofPath := filepath.Join(prestatesDir, "prestate-proof.json")
			if err := os.WriteFile(prestateProofPath, content, 0644); err != nil {
				return nil, fmt.Errorf("failed to copy prestate to prestates directory: %w", err)
			}

			log.Printf("Copied %s to prestates/prestate-proof.json for challenger compatibility", filepath.Base(filePath))
			prestateProofCopied = true
		}

		// Rename files to hash-based names
		newFileName := data.Pre + ".json"
		hashedPath := filepath.Join(buildDir, newFileName)
		if err := os.Rename(filePath, hashedPath); err != nil {
			return nil, fmt.Errorf("failed to rename prestate %s: %w", filepath.Base(filePath), err)
		}
		log.Printf("%s available at: %s/%s\n", filepath.Base(filePath), prestateURL, newFileName)

		// Rename the corresponding binary file
		binFilePath := strings.Replace(strings.TrimSuffix(filePath, ".json"), "-proof", "", 1) + ".bin.gz"
		newBinFileName := data.Pre + ".bin.gz"
		binHashedPath := filepath.Join(buildDir, newBinFileName)
		if err := os.Rename(binFilePath, binHashedPath); err != nil {
			return nil, fmt.Errorf("failed to rename prestate %s: %w", filepath.Base(binFilePath), err)
		}
		log.Printf("%s available at: %s/%s\n", filepath.Base(binFilePath), prestateURL, newBinFileName)
	}

	h.info = info

	return info, nil
}
