---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "nginxproxymanager_certificate_custom Resource - nginxproxymanager"
subcategory: "SSL Certificates"
description: |-
  This resource can be used to manage a custom certificate.
---

# nginxproxymanager_certificate_custom (Resource)

This resource can be used to manage a custom certificate.


## Example Usage

```terraform
resource "nginxproxymanager_certificate_custom" "certificate" {
  name = "Certificate"

  certificate     = file("certificate.pem")
  certificate_key = file("certificate.key")
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `certificate` (String, Sensitive) The contents of the certificate.
- `certificate_key` (String, Sensitive) The contents of the certificate key.
- `name` (String) The name of the certificate.

### Read-Only

- `created_on` (String) The date and time the certificate was created.
- `domain_names` (Set of String) The domain names associated with the certificate.
- `expires_on` (String) The date and time the certificate expires.
- `id` (Number) The ID of the certificate.
- `modified_on` (String) The date and time the certificate was last modified.
- `owner_user_id` (Number) The Id of the user that owns the certificate.

## Import

Import is supported using the following syntax:

```shell
# Certificates can be imported by specifying the numeric identifier of the certificate.
terraform import nginxproxymanager_certificate_custom.certificate 1
```
