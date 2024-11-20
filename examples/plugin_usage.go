package main

import (
	"context"
	"log"

	"github.com/yahao333/gort/internal/plugin"
)

func main() {
	ctx := context.Background()

	// Create plugin manager
	pm := plugin.NewPluginManager("./plugins", nil)

	// Initialize plugin manager
	if err := pm.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize plugin manager: %v", err)
	}

	// Load AWS provider plugin
	if err := pm.LoadPlugin(ctx, "aws-provider"); err != nil {
		log.Fatalf("Failed to load AWS provider plugin: %v", err)
	}

	// Get plugin instance
	p, err := pm.GetPlugin("aws-provider")
	if err != nil {
		log.Fatalf("Failed to get plugin: %v", err)
	}

	// Cast to provider plugin
	provider, ok := p.(plugin.ProviderPlugin)
	if !ok {
		log.Fatal("Invalid plugin type")
	}

	// Use the provider
	resource, err := provider.CreateResource(ctx, plugin.ResourceSpec{
		Type: "aws_instance",
		Name: "example",
		Properties: map[string]interface{}{
			"instance_type": "t2.micro",
			"ami":           "ami-12345678",
		},
	})

	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}

	log.Printf("Created resource: %+v", resource)
}
