package main // import "github.com/uphy/revealjs-docker/bootstrap"

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/uphy/revealjs-docker/bootstrap/revealjs"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "dir,d",
			Value: ".",
			Usage: "path to the reveal.js installation directory",
		},
	}

	var server *revealjs.RevealJS
	app.Before = func(ctx *cli.Context) error {
		dir := ctx.String("dir")
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return errors.New("`dir` not exist")
		}
		var err error
		server, err = revealjs.NewRevealJS(dir)
		if err != nil {
			return fmt.Errorf("failed to initialize app: %s", err)
		}
		return nil
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name:  "init",
			Usage: "Generate config file and slide files",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "overwrite,o",
				},
			},
			ArgsUsage: fmt.Sprintf("[%s]", strings.Join(revealjs.FilesetNames, "|")),
			Action: func(ctx *cli.Context) error {
				var name string
				if ctx.NArg() == 0 {
					name = revealjs.FilesetNames[0]
				} else {
					name = ctx.Args().First()
				}
				return revealjs.Generate(name, server.DataDirectory(), ctx.Bool("overwrite"))
			},
		},
		cli.Command{
			Name:  "start",
			Usage: "Start reveal.js server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dir,d",
					Value: ".",
				},
			},
			Action: func(ctx *cli.Context) error {
				if err := server.Start(); err != nil {
					return fmt.Errorf("failed to start server: %s", err)
				}

				var signalc chan os.Signal
				signalc = make(chan os.Signal, 1)
				signal.Notify(signalc, os.Interrupt)
				<-signalc
				return nil
			},
		},
		cli.Command{
			Name: "build",
			Action: func(ctx *cli.Context) error {
				return errors.New("not implemented yet")
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println("failed to execute: ", err)
	}
}
