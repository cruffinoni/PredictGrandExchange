package adf

type Adf struct {
	values []float64
	diff   float64
	alpha  float64
}

func (a *Adf) getNods() int {
	return len(a.values)
}

func (a *Adf) getLagValues(idx int) float64 {
	if idx <= 0 {
		panic("underflow value")
	}
	if idx > len(a.values) {
		panic("overflow value")
	}
	return a.values[idx-1]
}

func (a *Adf) prepareData() {
	for _, v := range a.values {
		a.diff -= v
	}
}

func New() *Adf {
	return &Adf{}
}
