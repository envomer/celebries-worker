package main

import (
	"any-days.com/celebs/db"
	"any-days.com/celebs/logger"
	"github.com/evolidev/evoli"
	"github.com/evolidev/evoli/framework/console"
	"github.com/joho/godotenv"
)

var appLog = logger.AppLog

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		appLog.Fatal("Failed to load env file: %s", err)
	}

	db.Db()
	c := evoli.NewApplication()

	mainConsole(c)
	c.Start()
}

func mainConsole(c *evoli.Application) {
	cli := c.Cli

	cli.AddCommand("web:serve {--port} {--secure=false}", "Serve webserver", webServer)
	cli.AddCommand("info", "Serve webserver", info)
}

func info(c *console.ParsedCommand) {

}

func webServer(c *console.ParsedCommand) {

}
