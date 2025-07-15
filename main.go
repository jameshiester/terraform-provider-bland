package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/jameshiester/terraform-provider-bland/common"
	"github.com/jameshiester/terraform-provider-bland/internal/provider"
)

// Generate the provider document.
//
//go:generate tfplugindocs generate --provider-name powerplatform --rendered-provider-name "Power Platform"
func main() {
	log.Printf("[INFO] Starting the Power Platform Terraform Provider %s %s", common.ProviderVersion, common.Branch)

	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	ctx := context.Background()

	serveOpts := providerserver.ServeOpts{
		Debug:   debug,
		Address: "registry.terraform.io/jameshiester/bland",
	}

	err := providerserver.Serve(ctx, provider.NewBlandProvider(ctx), serveOpts)

	if err != nil {
		log.Fatalf("Error serving provider: %s", err)
	}
}
