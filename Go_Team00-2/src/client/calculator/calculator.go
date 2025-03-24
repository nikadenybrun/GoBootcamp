package calculator

import (
	"math"
	"team00/client/model"
	"team00/config"
)

type Calculator interface {
	Calculate(value float64) model.CalculatorRate
}

type calculator struct {
	k         float64
	mean      float64
	variance  float64
	count     int
	threshold int
}

func NewCalculator(cfg config.Calculator) Calculator {
	return &calculator{
		k:         cfg.K,
		threshold: cfg.Threshold,
	}
}

func (c *calculator) Calculate(value float64) model.CalculatorRate {
	c.count++
	delta := value - c.mean
	c.mean += delta / float64(c.count)
	c.variance += delta * (value - c.mean)

	deviation := c.stddev()

	res := model.CalculatorRate{
		Mean:      c.mean,
		Deviation: deviation,
	}

	if c.getAnomaly(value, deviation) {
		res.Anomaly.Exist = true
		res.Anomaly.Frequency = value
		res.Anomaly.Diff = math.Abs(value - c.mean)
	}

	return res
}

func (c *calculator) getAnomaly(num float64, deviation float64) bool {
	if c.count < c.threshold {
		return false
	}

	anomalyFlag := c.k * deviation
	num = math.Abs(num)
	value := num - c.mean
	if value > anomalyFlag || (-1*value) < (-1*anomalyFlag) {
		return true
	}

	return false
}

func (c *calculator) stddev() float64 {
	if c.count < 2 {
		return 0
	}

	return math.Sqrt(c.variance / float64(c.count-1))
}
