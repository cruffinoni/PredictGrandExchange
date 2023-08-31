package main

type adf struct {
	values []float64
	diff   float64
	alpha  float64
}

func (a *adf) getLagValues(idx int) float64 {
	if idx <= 0 {
		panic("underflow value")
	}
	if idx > len(a.values) {
		panic("overflow value")
	}
	return a.values[idx-1]
}

func (a *adf) prepareData() {
	for _, v := range a.values {
		a.diff -= v
	}
}
