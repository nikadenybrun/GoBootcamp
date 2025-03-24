package model

type CalculatorRate struct {
	Mean      float64
	Deviation float64
	Anomaly   struct {
		Exist     bool
		Frequency float64
		Diff      float64
	}
}
