/*
 *
 * @Author evsio0n
 * @Date 2022/3/27 下午3:53
 * @Email <bbq2001820@gmail.com>
 *
 */

package main

import (
	"RSCEUtil"
	"github.com/evsio0n/log"
	"github.com/urfave/cli"
	"os"
	"sort"
)

const (
	name        = "rsce-go"
	author      = "evsio0n <evsio0n@nexet.hk>"
	version     = "0.0.1"
	description = "rsce-go is a tool for unpack or pack RSCE, aka rock chip resource image. \r\n" +
		"\t using -u to unpack, using -p to pack.\r\n" +
		"\t -u [filepath]  unpack rockchip rsce image\r\n" +
		"\t -p [filepath] [filepath] [filepath]  pack all resource to rockchip rsce image\r\n"
)

func main() {
	cmd := &cli.App{
		Name:        name,
		Author:      author,
		Version:     version,
		Description: description,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "unpack,u",
				Usage: "unpack rsce image",
			},
			cli.StringSliceFlag{
				Name:  "pack,p",
				Usage: "pack rsce image",
			},
		},
		Action: Handle,
	}
	sort.Sort(cli.FlagsByName(cmd.Flags))
	err := cmd.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
	}
}
func Handle(c *cli.Context) {

	if c.String("unpack") != "" {
		RSCEUtil.UnPackRSCE(c.String("unpack"))
	}

	if len(c.StringSlice("pack")) > 0 {
		RSCEUtil.GenerateRSCE(c.StringSlice("pack"), "./boot-second")
	}
}
