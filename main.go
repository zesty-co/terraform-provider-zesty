package main

import (
	"context"
	"log"

	"github.com/spf13/pflag"
	"github.com/zesty-co/terraform-provider-zesty/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	debug   bool
	version string = "local"
)

func init() {
	pflag.BoolVarP(&debug, "debug", "d", false, "set to true to run the provider with support for debuggers like delve")
	pflag.Parse()
}

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/hashicorp/zesty",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
