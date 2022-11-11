package main

import (
	"any-days.com/celebs/db"
	"any-days.com/celebs/logger"
	"any-days.com/celebs/model"
	"any-days.com/celebs/tmdb"
	"github.com/evolidev/evoli"
	"github.com/evolidev/evoli/framework/console"
	"github.com/evolidev/evoli/framework/use"
)

var appLog = logger.AppLog

func main() {
	//err := godotenv.Load(".env")
	//if err != nil {
	//	appLog.Fatal("Failed to load env file: %s", err)
	//}

	db.Db()
	c := evoli.NewApplication()

	mainConsole(c)
	c.Start()
}

func mainConsole(c *evoli.Application) {
	cli := c.Cli

	cli.AddCommand("web:serve {--port} {--secure=false}", "Serve webserver", webServer)
	cli.AddCommand("info", "Serve webserver", info)
	cli.AddCommand("db:migrate", "Migrate DB", migrate)

	cli.AddCommand("people:fetch {--page} {--limit}", "Fetch people from TMDB", fetchPeople)
	cli.AddCommand("people:fuse", "Fuse people", fusePeople)
}

func info(c *console.ParsedCommand) {

}

func webServer(c *console.ParsedCommand) {

}

func migrate(c *console.ParsedCommand) {
	start := use.TimeRecord()
	appLog.Debug("Migrate...")

	model.Migrate()

	appLog.Debug("Took %s", start.ElapsedColored())
}

func fetchPeople(c *console.ParsedCommand) {
	start := use.TimeRecord()
	appLog.Debug("Fetch people...")

	page := c.GetOptionWithDefault("page", "1").Integer()
	limit := c.GetOptionWithDefault("limit", "1000").Integer()

	tmdb.FetchPeople(page, limit)

	appLog.Debug("Took %s", start.ElapsedColored())
}

func fusePeople(c *console.ParsedCommand) {
	start := use.TimeRecord()
	appLog.Debug("Fuse people...")

	tmdb.FusePeople()

	appLog.Debug("Took %s", start.ElapsedColored())
}
