package handlers

import (
	"github.com/gofiber/fiber/v2"
	"poc-vault-go-kube/vault"
)

func SetupRoutes(app *fiber.App, vaultClient *vault.Client) {
	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	// Get cluster secret
	app.Get("/cluster-secret/:path", func(c *fiber.Ctx) error {
		path := c.Params("path")
		secret, err := vaultClient.GetClusterSecret(path)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"secret": secret,
		})
	})

	// Get application secret
	app.Get("/app-secret/:path", func(c *fiber.Ctx) error {
		path := c.Params("path")
		secret, err := vaultClient.GetAppSecret(path)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"secret": secret,
		})
	})
}
