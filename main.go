package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Metronome deploy Drone plugin"
	app.Usage = "metronome deploy Drone plugin"
	app.Action = run
	app.Flags = []cli.Flag{

		cli.StringFlag{
			Name:   "server",
			Usage:  "metronome server",
			Value:  "http://master.mesos:9000",
			EnvVar: "PLUGIN_SERVER",
		},
		cli.StringFlag{
			Name:   "metronomefile",
			Usage:  "application metronome file",
			EnvVar: "PLUGIN_METRONOMEFILE",
		},
		cli.StringFlag{
			Name:   "job_config",
			Usage:  "application in-line config",
			EnvVar: "PLUGIN_JOB_CONFIG",
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {

	plugin := Plugin{
		Server:        c.String("server"),
		Metronomefile: c.String("metronomefile"),
		JobConfig:     c.String("job_config"),
	}

	return plugin.Exec()
}
