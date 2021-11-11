package main

import (
	"os"

	nanny "github.com/voje/gonanny/internal/nanny"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "./config.yaml",
				Usage:   "Config file for the app",
				EnvVars: []string{"CONF"},
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		log.Info("App started!")
		log.Infof("Reading config from: %s", c.String("config"))
		conf, err := nanny.ConfigFromFile(c.String("config"))
		if err != nil {
			log.Fatal(err)
		}
		nannyApp := nanny.NewNanny(conf)
		err = nannyApp.Run()
		if err != nil {
			log.Fatal(err)
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
