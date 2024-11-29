package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
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

func ParseVariables(modulePath string) (map[string]interface{}, []string, error) {
	optionalVars := make(map[string]interface{})
	requiredVars := []string{}
	parser := hclparse.NewParser()

	// Read the variables.tf file
	variablesFile := filepath.Join(modulePath, "variables.tf")
	content, err := os.ReadFile(variablesFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read variables.tf: %w", err)
	}

	// Parse the file as raw HCL
	hclFile, diags := parser.ParseHCL(content, variablesFile)
	if diags.HasErrors() {
		return nil, nil, fmt.Errorf("failed to parse variables.tf: %v", diags)
	}

	// Iterate over all blocks in the file
	for _, block := range hclFile.Body.(*hclsyntax.Body).Blocks {
		if block.Type != "variable" {
			continue // Ignore non-variable blocks
		}

		// Ensure the block has exactly one label (variable name)
		if len(block.Labels) != 1 {
			return nil, nil, fmt.Errorf("variable block should have exactly one label (name), found: %v", block.Labels)
		}
		varName := block.Labels[0]

		// Check if the block has a "default" attribute
		defaultAttr, hasDefault := block.Body.Attributes["default"]
		if hasDefault {
			val, diags := defaultAttr.Expr.Value(nil)
			if diags.HasErrors() {
				return nil, nil, fmt.Errorf("failed to evaluate default value for variable %s: %v", varName, diags)
			}
			optionalVars[varName] = val.AsString()
		} else {
			// No "default" means it's a required variable
			requiredVars = append(requiredVars, varName)
		}
	}

	return optionalVars, requiredVars, nil
}

func ParseOutputs(modulePath string) (map[string]string, error) {
	outputs := make(map[string]string)
	parser := hclparse.NewParser()

	outputsFile := filepath.Join(modulePath, "outputs.tf")
	content, err := os.ReadFile(outputsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read outputs.tf: %w", err)
	}

	hclFile, diags := parser.ParseHCL(content, outputsFile)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse outputs.tf: %v", diags)
	}

	// Decode the entire body as raw content
	body := hclFile.Body
	contentBody, _, diags := body.PartialContent(&hcl.BodySchema{})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to decode outputs.tf: %v", diags)
	}

	for _, block := range contentBody.Blocks {
		if block.Type != "output" {
			continue // Ignore non-output blocks
		}

		if len(block.Labels) < 1 {
			return nil, fmt.Errorf("output block missing a name")
		}

		outputName := block.Labels[0]

		// Extract the value attribute
		attrs, diags := block.Body.JustAttributes()
		if diags.HasErrors() {
			return nil, fmt.Errorf("failed to read attributes for output %s: %v", outputName, diags)
		}

		if attr, ok := attrs["value"]; ok {
			// Evaluate the value of the expression
			val, diags := attr.Expr.Value(nil)
			if diags.HasErrors() {
				return nil, fmt.Errorf("failed to evaluate value expression for output %s: %v", outputName, diags)
			}

			// Store the evaluated value as a string
			outputs[outputName] = val.AsString()
		} else {
			return nil, fmt.Errorf("output %s does not have a value attribute", outputName)
		}
	}

	return outputs, nil
}
