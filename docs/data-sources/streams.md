---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "nginxproxymanager_streams Data Source - nginxproxymanager"
subcategory: ""
description: |-
  Stream data source.
---

# nginxproxymanager_streams (Data Source)

Stream data source.

## Example Usage

```terraform
# Fetch all streams
data "nginxproxymanager_streams" "all" {}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `streams` (Attributes List) The streams. (see [below for nested schema](#nestedatt--streams))

<a id="nestedatt--streams"></a>
### Nested Schema for `streams`

Read-Only:

- `created_on` (String) The date and time the stream was created.
- `enabled` (Boolean) Whether the stream is enabled.
- `forwarding_host` (String) The forwarding host of the stream.
- `forwarding_port` (Number) The forwarding port of the stream.
- `id` (Number) The ID of the stream.
- `incoming_port` (Number) The incoming port of the stream.
- `meta` (Map of String) The meta data associated with the stream.
- `modified_on` (String) The date and time the stream was last modified.
- `owner_user_id` (Number) The ID of the user that owns the stream.
- `tcp_forwarding` (Boolean) Whether TCP forwarding is enabled.
- `udp_forwarding` (Boolean) Whether UDP forwarding is enabled.

