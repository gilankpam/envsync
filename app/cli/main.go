package main

import (
	"fmt"
	"os"

	"github.com/bukalapak/envsync"
	"github.com/urfave/cli"
)

func main() {
	var source string
	var target string
	syncer := &envsync.Syncer{}

	app := cli.NewApp()
	app.Name = "envsync"
	app.Usage = "synchronize sample env and actual env file"
	app.UsageText = "envsync -s [sample env] -t [actual env]"
	app.Version = envsync.VERSION
	app.Copyright = "Bukalapak™ © 2018"
	app.Authors = []cli.Author{
		{
			Name: "PT. Bukalapak.com",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "source, s",
			Usage:       "set sample env",
			Value:       "env.sample",
			Destination: &source,
		},
		cli.StringFlag{
			Name:        "target, t",
			Usage:       "set actual env",
			Value:       ".env",
			Destination: &target,
		},
	}
	app.Action = func(c *cli.Context) error {
		newEnv, err := syncer.Sync(source, target)
		if err != nil {
			fmt.Println(err.Error())
		}
		addedLine := len(newEnv)
		if addedLine > 0 {
			fmt.Printf("Added %d new env variable\nRun \"tail -n %d %s\" to view them", addedLine, addedLine, target)
		} else {
			fmt.Printf("%s already uptodate", source)
		}
		return err
	}
	app.Run(os.Args)
}
