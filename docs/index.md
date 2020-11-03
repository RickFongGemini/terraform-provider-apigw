---
layout: "apigw"
page_title: "Provider: APIGW"
sidebar_current: "docs-apigw-index"
description: |-
  Terraform provider apigw.
---

# APIGW Provider

Test publishing an apigw provider.

This is a test repo.

## Example Usage

```hcl
provider "apigw" {
    apikey = "<APIKEY>"
    apigw_url = "<APIGW_URL>"
}

# Example resource configuration
resource "apigw_resourcename" "example" {
  # ...
}
```
