package command

import (
	"github.com/urfave/cli/v2"
)

func App() *cli.App {
	app := cli.NewApp()
	app.Name = "gdraw"

	app.Flags = []cli.Flag{

	}

	
	return nil
}
