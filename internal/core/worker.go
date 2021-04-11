package core

import (
	"image"
	"math/rand"
	"time"

	"github.com/golang/freetype/raster"
)

type Worker struct {
	W, H       int
	Target     *image.RGBA
	Current    *image.RGBA
	Buffer     *image.RGBA
	Rasterizer *raster.Rasterizer
	Lines      []Scanline
	Rnd        *rand.Rand
	Score      float64
	Counter    int
}

func NewWorker(target *image.RGBA) *Worker {
	w := target.Bounds().Size().X
	h := target.Bounds().Size().Y
	worker := Worker{
		W:          w,
		H:          h,
		Target:     target,
		Buffer:     image.NewRGBA(target.Bounds()),
		Rasterizer: raster.NewRasterizer(w, h),
		Lines:      make([]Scanline, 0, 4096),
		Rnd:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return &worker
}

func (worker *Worker) Init(current *image.RGBA, score float64) {
	worker.Current = current
	worker.Score = score
	worker.Counter = 0
}

func (worker *Worker) Energy(shape Shape, alpha int) float64 {
	worker.Counter++
	lines := shape.Rasterize()
	color := ComputeColor(worker.Target, worker.Current, lines, alpha)
	CopyLines(worker.Buffer, worker.Current, lines)
	DrawLines(worker.Buffer, color, lines)
	return DifferencePartial(worker.Target, worker.Current, worker.Buffer, worker.Score, lines)
}

func (worker *Worker) BestHillClimbState(shapeType ShapeType, alpha, n, age, m int) *State {
	var bestEnergy float64
	var bestState *State

	for i := 0; i < m; i++ {
		state := worker.BestRandomState(shapeType, alpha, n)
		// before := state.Energy()
		state = HillClimb(state, age).(*State)
		energy := state.Energy()
		// fmt.Printf("%dx random: %f -> %dx hill climb: %f\n", n, before, age, energy)
		if i == 0 || energy < bestEnergy {
			bestEnergy = energy
			bestState = state
		}
	}

	return bestState
}

func (worker *Worker) BestRandomState(shapeType ShapeType, a, n int) *State {
	var bestEnergy float64
	var bestState *State
	for i := 0; i < n; i++ {
		state := worker.RandomState(shapeType, a)
		energy := state.Energy()
		if i == 0 || energy < bestEnergy {
			bestEnergy = energy
			bestState = state
		}
	}
	return bestState
}

func (worker *Worker) RandomState(shapeType ShapeType, alpha int) *State {
	switch shapeType {
	default:
		return worker.RandomState(ShapeType(worker.Rnd.Intn(8)+1), alpha)
	case ShapeTypeTriangle:
		return NewState(worker, NewRandomTriangle(worker), alpha)
	case ShapeTypeRectangle:
		return NewState(worker, NewRandomRectangle(worker), alpha)
	case ShapeTypeEllipse:
		return NewState(worker, NewRandomEllipse(worker), alpha)
	case ShapeTypeCircle:
		return NewState(worker, NewRandomCircle(worker), alpha)
	case ShapeTypeRotatedRectangle:
		return NewState(worker, NewRandomRotatedRectangle(worker), alpha)
	case ShapeTypeRotatedEllipse:
		return NewState(worker, NewRandomRotatedEllipse(worker), alpha)
	case ShapeTypePolygon:
		return NewState(worker, NewRandomPolygon(worker, 4, false), alpha)
	}
}
