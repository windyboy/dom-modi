package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli/altsrc"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()

	flags := []cli.Flag{
		/*
		   // OK
		   altsrc.NewIntFlag(cli.IntFlag{Name: "test", Value: -1}),
		   altsrc.NewStringFlag(cli.StringFlag{Name: "flag", Value: "default"}),
		*/

		// NG #1
		altsrc.NewIntFlag(cli.IntFlag{Name: "test, t", Value: -1}),
		altsrc.NewStringFlag(cli.StringFlag{Name: "flag", Value: "default"}),
		altsrc.NewBoolFlag(cli.BoolFlag{Name: "change"}),

		/*
		   // NG #2
		   altsrc.NewIntFlag(cli.IntFlag{Name: "test", Value: -1}),
		   altsrc.NewStringFlag(cli.StringFlag{Name: "flag, f", Value: "default"}),

		   // NG #3
		   altsrc.NewIntFlag(cli.IntFlag{Name: "test, t", Value: -1}),
		   altsrc.NewStringFlag(cli.StringFlag{Name: "flag, f", Value: "default"}),
		*/

		cli.StringFlag{Name: "conf, c", Value: "cfg.yml"},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println("yaml ist rad")
		fmt.Printf("config: %s\n", c.String("conf"))
		fmt.Printf("test: %d\n", c.Int("test"))
		fmt.Printf("flag: %s\n", c.String("flag"))
		fmt.Printf("change: %t\n", c.Bool("change"))

		args := c.Args()
		fmt.Println(args)
		var domXml, err := ioutil.ReadFile(args.Get(0))
		if err != nil {
			fmt.Printf("read file error ", err)
			return err
		}
		fmt.Println(string(domXml))
		return nil
	}

	app.Before = altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("conf"))
	app.Flags = flags

	app.Run(os.Args)

}
