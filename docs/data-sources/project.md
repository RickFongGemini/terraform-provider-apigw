---
layout: "apigw"
page_title: "APIGW: apigw_project"
sidebar_current: "docs-apigw-project"
description: |-
  Project data source in the Terraform provider apigw.
---

# apigw_project

Project data source in the Terraform provider apigw.

## Example Usage

```hcl
data "apigw_project" "example" {
  name = "foo"
}
```

## Attributes Reference

* `name` - Project name.
