/*
 *
 * @Author evsio0n
 * @Date 2024/10/31 下午5:10
 * @Email <i@evsio0n.com>
 *
 */

package main

import (
	"github.com/Evsio0n/rsce-go/rsceUtil"
	log2 "github.com/evsio0n/log"
	"github.com/urfave/cli"
	"os"
	"sort"
)

var log *log2.Logger

const (
	name        = "rsce-go"
	author      = "evsio0n <i@evsio0n.com>"
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
