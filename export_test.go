package wordimg

func (g *Generator) SetRand(ivalue int, fvalue float64) {
	g.randFunc = func() float64 {
		return fvalue
	}
	g.colorGen.randFunc = func(_ int) int {
		return ivalue
	}
}
