// Package cmd
// cmd/serve.go
// @title zGinv API
// @version 1.0
// @description Централизованная инвентаризация серверов и API для управления
// @contact.name Евгений Гончаров
// @contact.url https://openbld.net
// @license.name MIT
// @host localhost:8080
// @BasePath /api
package cmd

import (
	"fmt"
	//"github.com/gofiber/swagger"
	"log"
	"os"
	"zGinv/api"
	"zGinv/db"

	"github.com/gofiber/fiber/v3"
	"github.com/spf13/cobra"

	"github.com/Flussen/swagger-fiber-v3"

	_ "zGinv/docs" // swaggo gen docs
)

var servePort int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Running HTTP API server",
	Run: func(cmd *cobra.Command, args []string) {
		db.InitDB()

		p := os.Getenv("ZGINV_PORT")
		if p == "" {
			p = fmt.Sprintf("%d", servePort)
		}

		r := fiber.New()
		r.Use(func(c fiber.Ctx) error {
			fmt.Printf("%s %s\n", c.Method(), c.Path())
			return c.Next()
		})

		r.Get("/swagger/*", swagger.HandlerDefault)

		r.Get("/", func(c fiber.Ctx) error {
			return c.SendString("🧩 zGinv API. Try /api/servers")
		})

		routes := r.Group("/api")
		api.RegisterRoutes(routes)

		log.Printf("🚀 Server runs on http://localhost:%s\n", p)
		log.Fatal(r.Listen(":" + p))
	},
}

func init() {
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8080, "Порт HTTP сервера (по умолчанию 8080)")
	rootCmd.AddCommand(serveCmd)
}
