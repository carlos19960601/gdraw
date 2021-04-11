package command

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/nfnt/resize"
	"github.com/urfave/cli/v2"
	"github.com/zengqiang96/gdraw/internal/color"
	"github.com/zengqiang96/gdraw/internal/core"
	"github.com/zengqiang96/gdraw/internal/imgutil"
	"github.com/zengqiang96/gdraw/internal/logger"
	"github.com/zengqiang96/gdraw/internal/model"
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
		&cli.IntFlag{
			Name:    "alpha",
			Aliases: []string{"a"},
			Usage:   "alpha值",
			Value:   128,
		},
		&cli.IntFlag{
			Name:    "repeat",
			Aliases: []string{"r"},
			Usage:   "",
			Value:   0,
		},
		&cli.IntFlag{
			Name:    "resize",
			Aliases: []string{"rs"},
			Usage:   "输出图片的大小",
			Value:   256,
		},
		&cli.IntFlag{
			Name:    "outsize",
			Aliases: []string{"s"},
			Usage:   "输出图片的大小",
			Value:   1024,
		},
		&cli.IntFlag{
			Name:    "workers",
			Aliases: []string{"w"},
			Usage:   "并行worker的数量(默认使用所有cpu的核心数)",
			Value:   0,
		},
		&cli.IntSliceFlag{
			Name:     "graphn",
			Aliases:  []string{"n"},
			Usage:    "图形数量(可以指定多个图形数量)",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "outputs",
			Aliases:  []string{"o"},
			Usage:    "输出结果路径(可以指定多个)",
			Required: true,
		},
	}

	app.Before = func(c *cli.Context) error {
		if err := logger.Init(); err != nil {
			return err
		}
		return nil
	}

	app.Action = func(c *cli.Context) error {
		logger.Info("加载图片", zap.String("图片路径", c.String("input")))
		img, err := imgutil.LoadImage(c.String("input"))
		if err != nil {
			return err
		}
		gns := c.IntSlice("graphn")
		for _, gn := range gns {
			if gn < 1 {
				return fmt.Errorf("图形数量必须大于0")
			}
		}

		size := uint(c.Int("resize"))
		if size > 0 {
			img = resize.Thumbnail(size, size, img, resize.Bilinear)
		}

		var bg color.Color
		if c.String("background") == "" {
			bg = color.MakeColor(imgutil.AverageImageColor(img))
		} else {
			bg = color.MakeHexColor(c.String("background"))
		}

		workers := c.Int("workers")
		if workers < 1 {
			workers = runtime.NumCPU()
		}

		outputs := c.StringSlice("outputs")

		model := model.NewModel(img, bg, c.Int("outsize"), workers)
		start := time.Now()
		frame := 0
		for index, gn := range gns {
			for i := 0; i < gn; i++ {
				frame++
				t := time.Now()
				n := model.Step(core.ShapeType(c.Int("mode")), c.Int("alpha"), c.Int("repeat"))
				nps := float64(n) / time.Since(t).Seconds()
				elapsed := time.Since(start).Seconds()
				fmt.Printf("%d: t=%.3f, score=%.6f, n=%d, n/s=%f\n", frame, elapsed, model.Score, n, nps)
				for _, output := range outputs {
					ext := strings.ToLower(filepath.Ext(output))
					last := index == len(gns)-1 && i == gn-1
					if last {
						path := output
						switch ext {
						case ".png":
							err = imgutil.SavePNG(path, model.Context.Image())
							if err != nil {
								return err
							}
						}
					}
				}
			}
		}
		return nil
	}

	return app
}
