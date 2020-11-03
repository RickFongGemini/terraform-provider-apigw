---
layout: "apigw"
page_title: "APIGW: apigw_s3_key"
sidebar_current: "docs-apigw-s3_key"
description: |-
  s3_key resource in the Terraform provider apigw.
---

# scaffolding_resource

s3_key resource in the Terraform provider apigw.

## Example Usage

```hcl
resource "apigw_s3_key" "example" {
    name = "foo"
    platform = data.apigw_project.exampleProject.platform
    project = data.apigw_project.exampleProject.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - Private s3 key name.

* `platform` - Ceph plaform name.

* `project` - Ceph project ID.

