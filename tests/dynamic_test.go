package test

import (
	"fmt"
	"os"
	"strings"
	"terraform_testing_poc/tests/utils"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestDynamicModule(t *testing.T) {
	//modulePath := "../modules/az_rg"
	modulePath := os.Getenv("MODULE_PATH")
	if modulePath == "" {
		t.Fatalf("MODULE_PATH environment variable is not set")
	}

	utils.AddProviderFile(t, modulePath)
	// Parse variables
	optionalVars, requiredVars, err := utils.ParseVariables(modulePath)
	if err != nil {
		t.Fatalf("Failed to parse variables.tf: %v", err)
	}
	fmt.Printf("Optional variables: %v\n", optionalVars)
	fmt.Printf("Required variables: %v\n", requiredVars)
	// Generate environment variables for required variables
	envVars := map[string]string{}

	for _, reqVar := range requiredVars {
		envVarName := "TF_VAR_" + reqVar

		// Check if the variable name contains "location"
		if strings.Contains(strings.ToLower(reqVar), "location") {
			envVars[envVarName] = "eastus" // Default value for location variables
		} else {
			envVars[envVarName] = "test-" + reqVar // Example value for other variables
		}
	}
	// Include optional variables with default values
	inputs := map[string]interface{}{}
	for key, value := range optionalVars {
		inputs[key] = value
	}

	// Parse outputs
	outputs, err := utils.ParseOutputs(modulePath)
	if err != nil {
		t.Fatalf("Failed to parse outputs.tf: %v", err)
	}

	// Debugging: Print environment variables
	for key, value := range envVars {
		fmt.Printf("Setting environment variable: %s=%s\n", key, value)
	}

	// Run Terraform tests
	t.Parallel()
	terraformOptions := &terraform.Options{
		TerraformDir: modulePath,
		Vars:         inputs,  // Optional variables
		EnvVars:      envVars, // Required variables as environment variables
	}

	defer terraform.Destroy(t, terraformOptions)
	//terraform.InitAndApply(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	// Validate outputs
	for outputName := range outputs {
		actualValue := terraform.Output(t, terraformOptions, outputName)
		assert.NotEmpty(t, actualValue, "Output "+outputName+" should not be empty")
	}
}

