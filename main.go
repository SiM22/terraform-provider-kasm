package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	kasmProvider "terraform-provider-kasm/internal/provider"
)

func main() {
	err := providerserver.Serve(context.Background(), func() provider.Provider {
		return kasmProvider.New()
	}, providerserver.ServeOpts{
		Address: "registry.terraform.io/hashicorp/kasm",
	})
	if err != nil {
		log.Fatal(err)
	}
}
