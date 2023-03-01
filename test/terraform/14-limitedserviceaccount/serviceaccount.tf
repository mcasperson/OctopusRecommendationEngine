resource "octopusdeploy_user" "bob" {
  display_name  = "Bob Smith"
  email_address = "bob.smith@example.com"
  is_active     = true
  is_service    = true
  username      = "bob"
}

output "service_account_id" {
  value = octopusdeploy_user.bob.id
}
