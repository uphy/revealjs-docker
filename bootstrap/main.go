package main // import "github.com/uphy/revealjs-docker/bootstrap"

import (
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/uphy/revealjs-docker/bootstrap/revealjs"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		cli.Command{
			Name: "start",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dir,d",
					Value: ".",
				},
			},
			Action: func(ctx *cli.Context) error {
				dir := ctx.String("dir")
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					return errors.New("`dir` not exist")
				}

				app, err := revealjs.NewRevealJS(dir)
				if err != nil {
					return fmt.Errorf("failed to initialize app: %s", err)
				}

				if err := app.Start(); err != nil {
					return fmt.Errorf("failed to start app: %s", err)
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
