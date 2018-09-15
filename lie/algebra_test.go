package lie

import (
	"testing"

	"github.com/mjschust/cblocks/util"
)

func TestTypeAConvertWeightToEpc(t *testing.T) {
	cases := []struct {
		alg  TypeA
		wt   Weight
		want []int
	}{
		{TypeA{1}, Weight{0}, []int{0, 0}},
		{TypeA{1}, Weight{1}, []int{1, 0}},
		{TypeA{1}, Weight{2}, []int{2, 0}},
		{TypeA{2}, Weight{0, 0}, []int{0, 0, 0}},
		{TypeA{2}, Weight{1, 0}, []int{1, 0, 0}},
		{TypeA{2}, Weight{0, 1}, []int{1, 1, 0}},
		{TypeA{2}, Weight{1, 1}, []int{2, 1, 0}},
	}

	for _, c := range cases {
		got := make([]int, c.alg.rank+1)
		got = c.alg.convertWeightToEpc(c.wt, got)
		if !equals(got, c.want) {
			t.Errorf("convertWeightToEpc(%v) = %v, want %v", c.wt, got, c.want)
		}
	}
}

func TestTypeAConvertEpCoord(t *testing.T) {
	cases := []struct {
		alg  TypeA
		epc  []int
		want Weight
	}{
		{TypeA{1}, []int{0, 0}, Weight{0}},
		{TypeA{1}, []int{1, 1}, Weight{0}},
		{TypeA{1}, []int{1, 0}, Weight{1}},
		{TypeA{1}, []int{0, 1}, Weight{-1}},
		{TypeA{2}, []int{0, 0, 0}, Weight{0, 0}},
		{TypeA{2}, []int{1, 1, 1}, Weight{0, 0}},
		{TypeA{2}, []int{1, 0, 0}, Weight{1, 0}},
		{TypeA{2}, []int{1, 1, 0}, Weight{0, 1}},
		{TypeA{2}, []int{2, 1, 0}, Weight{1, 1}},
		{TypeA{2}, []int{1, 2, 0}, Weight{-1, 2}},
	}

	for _, c := range cases {
		var got Weight = make([]int, c.alg.rank)
		c.alg.convertEpCoord(c.epc, got)
		if !equals(got, c.want) {
			t.Errorf("convertEpCoord(%v) = %v, want %v", c.epc, got, c.want)
		}
	}
}

func TestTypeADualCoxeter(t *testing.T) {
	cases := []struct {
		alg  Algebra
		want int
	}{
		{TypeA{1}, 2},
		{TypeA{2}, 3},
		{TypeA{3}, 4},
	}

	for _, c := range cases {
		got := c.alg.DualCoxeter()
		if got != c.want {
			t.Errorf("DualCoxeter() == %v, want %v", got, c.want)
		}
	}
}

func TestTypeAPositiveRoots(t *testing.T) {
	cases := []struct {
		alg  Algebra
		want []Root
	}{
		{TypeA{1}, []Root{Root{1}}},
		{TypeA{2}, []Root{Root{1, 0}, Root{1, 1}, Root{0, 1}}},
		{TypeA{3}, []Root{
			Root{1, 0, 0},
			Root{1, 1, 0},
			Root{1, 1, 1},
			Root{0, 1, 0},
			Root{0, 1, 1},
			Root{0, 0, 1},
		}},
	}

	for _, c := range cases {
		got := c.alg.PositiveRoots()
		if len(got) != len(c.want) {
			t.Errorf("len(PositiveRoots()) == %v, want %v", len(got), len(c.want))
		}
		for i := range c.want {
			if !equals(got[i], c.want[i]) {
				t.Errorf("PositiveRoots() == %v, want %v", got, c.want)
			}
		}
	}
}

func TestTypeAWeights(t *testing.T) {
	cases := []struct {
		alg   Algebra
		level int
		want  [][]int
	}{
		{TypeA{1}, 0, [][]int{{0}}},
		{TypeA{1}, 1, [][]int{{0}, {1}}},
		{TypeA{1}, 2, [][]int{{0}, {1}, {2}}},
		{TypeA{2}, 0, [][]int{{0, 0}}},
		{TypeA{2}, 1, [][]int{{0, 0}, {1, 0}, {0, 1}}},
		{TypeA{2}, 2, [][]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}, {2, 0}, {0, 2}}},
		{TypeA{3}, 0, [][]int{
			{0, 0, 0},
		}},
		{TypeA{3}, 1, [][]int{
			{0, 0, 0},
			{1, 0, 0},
			{0, 1, 0},
			{0, 0, 1},
		}},
		{TypeA{3}, 2, [][]int{
			{0, 0, 0},
			{1, 0, 0},
			{0, 1, 0},
			{0, 0, 1},
			{1, 1, 0},
			{0, 1, 1},
			{1, 0, 1},
			{2, 0, 0},
			{0, 2, 0},
			{0, 0, 2},
		}},
	}

	for _, c := range cases {
		wantSet := util.NewVectorMap()
		for _, wt := range c.want {
			wantSet.Put(wt, true)
		}
		got := c.alg.Weights(c.level)
		if len(got) != len(c.want) {
			t.Errorf("len(Weights(%v)) == %v, want %v", c.level, len(got), len(c.want))
		}
		for _, gotWt := range got {
			_, present := wantSet.Remove(gotWt)
			if !present {
				t.Errorf("Weights(%v) should not contain %v", c.level, gotWt)
			}
		}
		if wantSet.Size() != 0 {
			t.Errorf("Weight(%v) is missing %v", c.level, wantSet.Keys())
		}
	}
}

func TestTypeAKillingForm(t *testing.T) {
	cases := []struct {
		alg      Algebra
		wt1, wt2 Weight
		want     float64
	}{
		{TypeA{1}, Weight{0}, Weight{0}, 0},
		{TypeA{1}, Weight{1}, Weight{0}, 0},
		{TypeA{1}, Weight{0}, Weight{1}, 0},
		{TypeA{1}, Weight{1}, Weight{1}, 0.5},
		{TypeA{1}, Weight{2}, Weight{1}, 1},
		{TypeA{1}, Weight{1}, Weight{2}, 1},
		{TypeA{1}, Weight{2}, Weight{2}, 2},
		{TypeA{2}, Weight{0, 0}, Weight{0, 0}, 0},
		{TypeA{2}, Weight{1, 0}, Weight{0, 0}, 0},
		{TypeA{2}, Weight{0, 0}, Weight{1, 0}, 0},
		{TypeA{2}, Weight{1, 0}, Weight{1, 0}, 0.6666666666666666},
		{TypeA{2}, Weight{0, 1}, Weight{1, 0}, 0.33333333333333333},
		{TypeA{2}, Weight{1, 0}, Weight{0, 1}, 0.33333333333333333},
		{TypeA{2}, Weight{0, 1}, Weight{0, 1}, 0.6666666666666666},
	}

	for _, c := range cases {
		got := c.alg.KillingForm(c.wt1, c.wt2)
		if got != c.want {
			t.Errorf("KillingForm(%v, %v) == %v, want %v", c.wt1, c.wt2, got, c.want)
		}
	}
}

func TestTypeAIntKillingForm(t *testing.T) {
	cases := []struct {
		alg      Algebra
		wt1, wt2 Weight
		want     int
	}{
		{TypeA{1}, Weight{0}, Weight{0}, 0},
		{TypeA{1}, Weight{1}, Weight{0}, 0},
		{TypeA{1}, Weight{0}, Weight{1}, 0},
		{TypeA{1}, Weight{1}, Weight{1}, 1},
		{TypeA{1}, Weight{2}, Weight{1}, 2},
		{TypeA{1}, Weight{1}, Weight{2}, 2},
		{TypeA{1}, Weight{2}, Weight{2}, 4},
		{TypeA{2}, Weight{0, 0}, Weight{0, 0}, 0},
		{TypeA{2}, Weight{1, 0}, Weight{0, 0}, 0},
		{TypeA{2}, Weight{0, 0}, Weight{1, 0}, 0},
		{TypeA{2}, Weight{1, 0}, Weight{1, 0}, 2},
		{TypeA{2}, Weight{0, 1}, Weight{1, 0}, 1},
		{TypeA{2}, Weight{1, 0}, Weight{0, 1}, 1},
		{TypeA{2}, Weight{0, 1}, Weight{0, 1}, 2},
	}

	for _, c := range cases {
		got := c.alg.IntKillingForm(c.wt1, c.wt2)
		if got != c.want {
			t.Errorf("KillingForm(%v, %v) == %v, want %v", c.wt1, c.wt2, got, c.want)
		}
	}
}

func TestTypeAKillingFactor(t *testing.T) {
	cases := []struct {
		alg  Algebra
		want int
	}{
		{TypeA{1}, 2},
		{TypeA{2}, 3},
		{TypeA{3}, 4},
	}

	for _, c := range cases {
		got := c.alg.KillingFactor()
		if got != c.want {
			t.Errorf("DualCoxeter() == %v, want %v", got, c.want)
		}
	}
}

func TestTypeALevel(t *testing.T) {
	cases := []struct {
		alg  Algebra
		wt   Weight
		want int
	}{
		{TypeA{1}, Weight{0}, 0},
		{TypeA{1}, Weight{1}, 1},
		{TypeA{1}, Weight{2}, 2},
		{TypeA{2}, Weight{0, 0}, 0},
		{TypeA{2}, Weight{1, 0}, 1},
		{TypeA{2}, Weight{0, 1}, 1},
		{TypeA{2}, Weight{1, 1}, 2},
	}

	for _, c := range cases {
		got := c.alg.Level(c.wt)
		if got != c.want {
			t.Errorf("Level(%v) == %v, want %v", c.wt, got, c.want)
		}
	}
}

func TestTypeADual(t *testing.T) {
	cases := []struct {
		alg      Algebra
		wt, want Weight
	}{
		{TypeA{1}, Weight{0}, Weight{0}},
		{TypeA{1}, Weight{1}, Weight{1}},
		{TypeA{1}, Weight{2}, Weight{2}},
		{TypeA{2}, Weight{0, 0}, Weight{0, 0}},
		{TypeA{2}, Weight{1, 0}, Weight{0, 1}},
		{TypeA{2}, Weight{0, 1}, Weight{1, 0}},
		{TypeA{2}, Weight{1, 1}, Weight{1, 1}},
	}

	for _, c := range cases {
		got := c.alg.NewWeight()
		got.Dual(c.alg, c.wt)
		if !equals(got, c.want) {
			t.Errorf("Dual(%v) = %v, want %v", c.wt, got, c.want)
		}
	}
}

func TestTypeAReflectIntoChamber(t *testing.T) {
	cases := []struct {
		alg      Algebra
		wt, want Weight
		parity   int
	}{
		{TypeA{1}, Weight{0}, Weight{0}, 1},
		{TypeA{1}, Weight{1}, Weight{1}, 1},
		{TypeA{1}, Weight{-1}, Weight{1}, -1},
		{TypeA{2}, Weight{0, 0}, Weight{0, 0}, 1},
		{TypeA{2}, Weight{1, 0}, Weight{1, 0}, 1},
		{TypeA{2}, Weight{0, 1}, Weight{0, 1}, 1},
		{TypeA{2}, Weight{-1, 0}, Weight{0, 1}, 1},
		{TypeA{2}, Weight{0, -1}, Weight{1, 0}, 1},
		{TypeA{2}, Weight{-1, -1}, Weight{1, 1}, -1},
	}

	for _, c := range cases {
		got := c.alg.NewWeight()
		parity := got.ReflectToChamber(c.alg, c.wt)
		if !equals(got, c.want) || parity != c.parity {
			t.Errorf("ReflectToChamber(%v) = %v, %v, want %v, %v",
				c.wt, got, parity, c.want, c.parity)
		}
	}
}

func TestTypeAOrbitIterator(t *testing.T) {
	cases := []struct {
		alg   Algebra
		wt    Weight
		orbit []Weight
	}{
		{TypeA{1}, Weight{0}, []Weight{Weight{0}}},
		{TypeA{1}, Weight{1}, []Weight{Weight{1}, Weight{-1}}},
		{TypeA{1}, Weight{2}, []Weight{Weight{2}, Weight{-2}}},
		{TypeA{2}, Weight{0, 0}, []Weight{Weight{0, 0}}},
		{TypeA{2}, Weight{1, 0}, []Weight{
			Weight{1, 0},
			Weight{-1, 1},
			Weight{0, -1}}},
		{TypeA{2}, Weight{0, 1}, []Weight{
			Weight{0, 1},
			Weight{1, -1},
			Weight{-1, 0}}},
		{TypeA{2}, Weight{1, 1}, []Weight{
			Weight{1, 1},
			Weight{-1, 2},
			Weight{2, -1},
			Weight{1, -2},
			Weight{-2, 1},
			Weight{-1, -1}}},
	}

	for _, c := range cases {
		orbitSet := weightSetFromList(c.orbit)
		orbitIter := c.alg.NewOrbitIterator(c.wt)
		orbitSize := 0
		for orbitIter.HasNext() {
			nextWt := c.alg.NewWeight()
			nextWt = orbitIter.Next(nextWt)
			_, present := orbitSet.Get(nextWt)
			if !present {
				t.Errorf("OrbitIterator(%v) does not contain %v", c.wt, nextWt)
			}
			orbitSize++
		}
		if orbitSize != len(c.orbit) {
			t.Errorf("OrbitIterator(%v) is missing orbit elements", c.wt)
		}
	}
}

func weightSetFromList(wts []Weight) util.VectorMap {
	vmap := util.NewVectorMap()
	for _, wt := range wts {
		vmap.Put(wt, true)
	}

	return vmap
}

func equals(v1, v2 []int) bool {
	if (v1 == nil) != (v2 == nil) {
		return false
	}

	if len(v1) != len(v2) {
		return false
	}

	for i := range v1 {
		if v1[i] != v2[i] {
			return false
		}
	}

	return true
}
