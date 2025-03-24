package generator

import (
	"math/rand/v2"
)

type Generator interface {
	Generate() float64
}

type generator struct {
	mean      float64
	deviation float64
}

func NewGenerator() Generator {
	return &generator{
		mean:      rand.Float64()*20 - 10,
		deviation: rand.Float64()*1.2 + 0.3,
	}
}

func (g *generator) Generate() float64 {
	return rand.NormFloat64()*g.deviation + g.mean
}
