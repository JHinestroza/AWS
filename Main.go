package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"os"
	"path/filepath"
)

type usuario struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Console struct {
	Code string `json:"code"`
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world!")
	})

	app.Post("/interpreter", parser)
	app.Post("/Login", LOGIN)
	app.Get("/SelectFile", SelectFile)

	app.Listen(":8080")
}
func parser(c *fiber.Ctx) error {
	var console Console
	if err := c.BodyParser(&console); err != nil {
		fmt.Println(console.Code)
		return err
	}

	result := Exec(console.Code)

	return c.JSON(fiber.Map{
		"message": result,
	})

}
func LOGIN(c *fiber.Ctx) error {
	var user usuario
	if err := c.BodyParser(&user); err != nil {
		return err
	}

	result := "login -user=" + user.Username + " -pass=" + user.Password

	return c.JSON(fiber.Map{

		"message": result,
	})
}

func SelectFile(c *fiber.Ctx) error {
	var files []string
	root := "Discos/"

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, filepath.Base(path))
		}
		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to list files")
	}

	return c.JSON(files)
}
