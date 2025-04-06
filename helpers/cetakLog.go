package helpers

import (
    "github.com/gofiber/fiber/v2"
    "log"
)

func Log(c *fiber.Ctx, jenis string, message string) {
	var color string
	if jenis == "error" {
		color = "\033[31m"
	} else if jenis == "success" {
		color = "\033[32m"
	} else if jenis == "warning" {
		color = "\033[33m"
	} else if jenis == "info" {
		color = "\033[34m"
	}else{
		jenis = "debug"
		color = "\033[0m"
	}

	bg := "\033[48;2;10;10;10m"
	reset := "\033[0m"

	log.Println(
		color + "[Info]" + reset + " : " + color + bg + message + 
		reset + "\n" + "                   Dari IP : " + color + bg + c.IP() + reset + 
		reset + "\n" + "           	  Dari URL : " + color + bg + c.Get("Referer") + reset + 
		reset + "\n" + "            Request ke URL : " + color + bg + c.OriginalURL() + reset + "\n",
	)
}