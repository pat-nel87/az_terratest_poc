package test

import (
	"terraform_testing_poc/tests/utils"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformResourceGroup(t *testing.T) {
	t.Parallel()

	utils.AddProviderFile(t, "../modules/az_rg") // Add provider.tf file

	terraformOptions := &terraform.Options{
		TerraformDir: "../modules/az_rg", // Path to your module

		Vars: map[string]interface{}{
			"resource_group_name":     "TestResourceGroup",
			"resource_group_location": "EastUS",
		},
	}

	// Ensure terraform destroy runs at the end of the test
	defer terraform.Destroy(t, terraformOptions)

	// Run terraform init and apply
	terraform.InitAndApply(t, terraformOptions)

	// Validate outputs
	resourceGroupId := terraform.Output(t, terraformOptions, "resource_group_id")
	assert.Contains(t, resourceGroupId, "/resourceGroups/TestResourceGroup")
}
