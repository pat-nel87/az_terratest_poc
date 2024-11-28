resource "azurerm_resource_group" "az_terra_test_rg" {
  name     = var.resource_group_name
  location = var.resource_group_location
}
