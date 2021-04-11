package model

import (
	"image"

	"github.com/fogleman/gg"
	"github.com/zengqiang96/gdraw/internal/color"
	"github.com/zengqiang96/gdraw/internal/core"
	"github.com/zengqiang96/gdraw/internal/imgutil"
	"github.com/zengqiang96/gdraw/internal/logger"
	"go.uber.org/zap"
)

type Model struct {
	Sw, Sh     int
	Scale      float64
	Background color.Color
	Target     *image.RGBA
	Current    *image.RGBA
	Score      float64
	Shapes     []core.Shape
	Colors     []color.Color
	Scores     []float64
	Context    *gg.Context
	Workers    []*core.Worker
}

func NewModel(target image.Image, background color.Color, size, numWorker int) *Model {
	w := target.Bounds().Size().X
	h := target.Bounds().Size().Y
	aspect := float64(w) / float64(h)
	var sw, sh int
	var scale float64
	if aspect >= 1 {
		sw = size
		sh = int(float64(size) / aspect)
		scale = float64(size) / float64(w)
	} else {
		sh = size
		sw = int(float64(size) * aspect)
		scale = float64(size) / float64(h)
	}

	logger.Info("输出的大小", zap.Int("sw", sw), zap.Int("sh", sh))

	model := Model{
		Sw:         sw,
		Sh:         sh,
		Scale:      scale,
		Background: background,
		Target:     imgutil.Image2RGBA(target),
		Current:    imgutil.UniformRGBA(target.Bounds(), background.NRGBA()),
	}
	model.Score = core.DifferenceFull(model.Target, model.Current)
	model.Context = model.newContext()

	for i := 0; i < numWorker; i++ {
		worker := core.NewWorker(model.Target)
		model.Workers = append(model.Workers, worker)
	}
	return &model
}

func (model *Model) newContext() *gg.Context {
	dc := gg.NewContext(model.Sw, model.Sh)
	dc.Scale(model.Scale, model.Scale)
	dc.Translate(0.5, 0.5)
	dc.SetColor(model.Background.NRGBA())
	dc.Clear()
	return dc
}

func (model *Model) Frames(scoreDelta float64) []image.Image {
	var result []image.Image
	dc := model.newContext()
	result = append(result, imgutil.Image2RGBA(dc.Image()))
	previous := 10.0
	for i, shape := range model.Shapes {
		c := model.Colors[i]
		dc.SetRGBA255(c.R, c.G, c.B, c.A)
		shape.Draw(dc, model.Scale)
		dc.Fill()
		score := model.Scores[i]
		delta := previous - score
		if delta >= scoreDelta {
			previous = score
			result = append(result, imgutil.Image2RGBA(dc.Image()))
		}
	}
	return result
}

func (model *Model) Step(shapeType core.ShapeType, alpha, repeat int) int {
	state := model.runWorkers(shapeType, alpha, 1000, 100, 16)
	model.add(state.Shape, state.Alpha)

	for i := 0; i < repeat; i++ {
		state.Worker.Init(model.Current, model.Score)
		a := state.Energy()
		state = core.HillClimb(state, 100).(*core.State)
		b := state.Energy()
		if a == b {
			break
		}
		model.add(state.Shape, state.Alpha)
	}
	counter := 0
	for _, worker := range model.Workers {
		counter += worker.Counter
	}
	return counter
}

func (model *Model) add(shape core.Shape, alpha int) {
	before := imgutil.CopyRGBA(model.Current)
	lines := shape.Rasterize()
	color := core.ComputeColor(model.Target, model.Current, lines, alpha)
	core.DrawLines(model.Current, color, lines)
	score := core.DifferencePartial(model.Target, before, model.Current, model.Score, lines)

	model.Score = score
	model.Shapes = append(model.Shapes, shape)
	model.Colors = append(model.Colors, color)
	model.Scores = append(model.Scores, score)

	model.Context.SetRGBA255(color.R, color.G, color.B, color.A)
	shape.Draw(model.Context, model.Scale)
}

func (model *Model) runWorkers(shapeType core.ShapeType, alpha, n, age, m int) *core.State {
	wn := len(model.Workers)
	ch := make(chan *core.State, wn)

	wm := m / wn
	if m%wn != 0 {
		wm++
	}
	for _, worker := range model.Workers {
		worker.Init(model.Current, model.Score)
		go model.runWorker(worker, shapeType, alpha, n, age, wm, ch)
	}

	var bestEnergy float64
	var bestState *core.State
	for i := 0; i < wn; i++ {
		state := <-ch
		energy := state.Energy()
		if i == 0 || energy < bestEnergy {
			bestEnergy = energy
			bestState = state
		}
	}
	return bestState
}

func (model *Model) runWorker(worker *core.Worker, shapeType core.ShapeType, alpha, n, age, m int, ch chan *core.State) {
	ch <- worker.BestHillClimbState(shapeType, alpha, n, age, m)
}
