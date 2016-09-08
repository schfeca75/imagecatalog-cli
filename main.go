package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	ic "github.com/hortonworks/imagecatalog-cli/cli"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "imagecatalog"
	app.Usage = "Cli for cloudbreak imagecatalog handling"
	app.Version = ic.Version + "-" + ic.BuildTime
	app.Author = "Hortonworks\n\nLICENSE:" + `
	 Apache License
	 Version 2.0, January 2004
	 http://www.apache.org/licenses/`

	app.Flags = []cli.Flag{
		ic.FlDebug,
	}

	app.Before = func(c *cli.Context) error {
		log.SetOutput(os.Stderr)
		log.SetLevel(log.ErrorLevel)
		formatter := &ic.LogFormatter{}
		log.SetFormatter(formatter)
		if c.Bool(ic.FlDebug.Name) {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:        "addversion",
			Description: fmt.Sprintf("add image information to the imagecatalog for the cloudbreak / ambari / hdp version triplet"),
			Usage:       "add image information to the imagecatalog for the cloudbreak / ambari / hdp version triplet",
			Flags:       []cli.Flag{ic.FlImageCatalog, ic.FlOutputImageCatalog, ic.FlCloudbreakVersion, ic.FlAmbariVersion, ic.FlHdpVersion},
			Action:      ic.AddVersion,
		},
	}

	app.Run(os.Args)
}
