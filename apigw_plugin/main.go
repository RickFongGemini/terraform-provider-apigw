package main

import (
    "apigw_plugin/terraform-provider-apigw/apigw"

    "github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
    plugin.Serve(&plugin.ServeOpts{
        ProviderFunc: apigw.Provider})
}

