package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/jspenc72/terraform-provider-bitbucket/bitbucket"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: bitbucket.Provider})
}
