package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"poc-vault-go-kube/config"
	"poc-vault-go-kube/handlers"
	"poc-vault-go-kube/vault"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Vault client
	vaultClient, err := vault.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Vault client: %v", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "vault-demo",
	})

	// Setup routes
	handlers.SetupRoutes(app, vaultClient)

	// Start server
	log.Fatal(app.Listen(":3000"))
}
