package command

import (
	"github.com/urfave/cli/v2"
	"github.com/zengqiang96/gdraw/internal/logger"
	"github.com/zengqiang96/gdraw/internal/util"
	"go.uber.org/zap"
)

const usage = `
.___                            
____     __| _/ _______  _____    __  _  __
/ ___\   / __ |  \_  __ \ \__  \   \ \/ \/ /
/ /_/  > / /_/ |   |  | \/  / __ \_  \     / 
\___  /  \____ |   |__|    (____  /   \/\_/  
/_____/        \/                \/           

使用几何图形绘画
`

func App() *cli.App {
	app := cli.NewApp()
	app.Name = "gdraw"
	app.Description = "使用几何图形绘画"
	app.Usage = usage

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "input",
			Aliases:  []string{"i"},
			Usage:    "原始图片路径",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "background",
			Aliases: []string{"b"},
			Usage:   "背景颜色(hex)",
		},
		&cli.IntFlag{
			Name:    "mode",
			Aliases: []string{"m"},
			Usage:   "0=combo 1=triangle 2=rect 3=ellipse 4=circle 5=rotatedrect 6=beziers 7=rotatedellipse 8=polygon",
			Value:   1,
		},
	}

	app.Before = func(c *cli.Context) error {
		if err := logger.Init(); err != nil {
			return err
		}
		return nil
	}

	app.Action = func(c *cli.Context) error {
		zap.S().Info("=====")

		_, err := util.LoadImage(c.String("input"))
		if err != nil {
			return err
		}
		return nil
	}

	return app
}
