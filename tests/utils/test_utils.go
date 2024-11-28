package utils

import (
	"os"
	"path/filepath"
	"testing"
)

// AddProviderFile dynamically creates a provider.tf file in the specified directory.
func AddProviderFile(t *testing.T, moduleDir string) {
	providerContent := `
provider "azurerm" {
  features {}
}`
	providerFilePath := filepath.Join(moduleDir, "provider.tf")
	if err := os.WriteFile(providerFilePath, []byte(providerContent), 0644); err != nil {
		t.Fatalf("Failed to write provider.tf: %v", err)
	}

	// Clean up the provider.tf file after the test
	t.Cleanup(func() {
		if err := os.Remove(providerFilePath); err != nil {
			t.Errorf("Failed to remove provider.tf: %v", err)
		}
	})
}
