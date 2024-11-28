provider "azurerm" {
  features {}
}

module "az_rg" {
  source              = "./modules/az_rg"
  resource_group_name = "az-terra-test-rg"
  resource_group_location = "eastus"
}