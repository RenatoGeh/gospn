package learn

import (
	"github.com/RenatoGeh/gospn/common"
	"github.com/RenatoGeh/gospn/learn/parameters"
	"github.com/RenatoGeh/gospn/spn"
	"math"
)

// Discriminative performs discriminative parameter learning, taking parameters from the
// parameters.P object bounded to the SPN S. If no parameters.P is found, uses default parameters.
// See parameters.P for more information.
func Discriminative(S spn.SPN, D spn.Dataset, Y []*Variable) spn.SPN {
	P, e := parameters.Retrieve(S)
	if !e {
		P = parameters.Default()
	}
	b := P.BatchSize > 1
	if parameters.Hardness(P.LearningType) == parameters.Hard {
		if b {
			return DiscriminativeHardBGD(S, P.Eta, P.Epsilon, D, Y, P.Normalize, P.BatchSize)
		}
		return DiscriminativeHardGD(S, P.Eta, P.Epsilon, D, Y, P.Normalize)
	}
	if b {
		return DiscriminativeBGD(S, P.Eta, P.Epsilon, D, Y, P.Normalize, P.BatchSize)
	}
	return DiscriminativeGD(S, P.Eta, P.Epsilon, D, Y, P.Normalize)
}

func pullValues(I spn.VarSet, Y []*Variable, y []int) {
	for i, v := range Y {
		u := v.Varid
		y[i] = I[u]
		delete(I, u)
	}
}

func pushValues(I spn.VarSet, Y []*Variable, y []int) {
	for i, v := range Y {
		I[v.Varid] = y[i]
	}
}

// DiscriminativeGD performs discriminative gradient descent on SPN S given data D. Argument eta is
// the learning rate, eps is the convergence difference in likelihood, D is the dataset and norm
// signals whether to normalize weights at each update.
func DiscriminativeGD(S spn.SPN, eta, eps float64, D spn.Dataset, Y []*Variable, norm bool) spn.SPN {
	st := spn.NewStorer()
	Q := &common.Queue{}
	s, d, w := st.NewTicket(), st.NewTicket(), st.NewTicket()
	z, p, u := st.NewTicket(), st.NewTicket(), st.NewTicket()
	P := S.Parameters()
	y := make([]int, len(Y))
	for i := 0; i < P.Iterations; i++ {
		for _, I := range D {
			spn.StoreInference(S, I, s, st)
			DeriveSPN(S, st, d, s, Q)
			DeriveWeights(S, st, w, d, s, Q)
			pullValues(I, Y, y)
			spn.StoreInference(S, I, z, st)
			DeriveSPN(S, st, p, z, Q)
			DeriveWeights(S, st, u, p, z, Q)
			pushValues(I, Y, y)
			applyDGD(S, s, d, w, z, p, u, st, eta, norm, Q)
			st.ResetTickets(s, d, w, z, p, u)
		}
	}
	return S
}

// DiscriminativeGD performs hard (MPE) discriminative gradient descent on SPN S given data D.
// Argument eta is the learning rate, eps is the convergence difference in likelihood, D is the
// dataset and norm signals whether to normalize weights at each update.
func DiscriminativeHardGD(S spn.SPN, eta, eps float64, D spn.Dataset, Y []*Variable, norm bool) spn.SPN {
	st := spn.NewStorer()
	d, p := st.NewTicket(), st.NewTicket()
	P := S.Parameters()
	y := make([]int, len(Y))
	for i := 0; i < P.Iterations; i++ {
		for _, I := range D {
			DeriveHard(S, st, d, I)
			pullValues(I, Y, y)
			DeriveHard(S, st, p, I)
			pushValues(I, Y, y)
			applyHDGD(S, d, p, st, eta, norm)
			st.ResetTickets(d, p)
		}
	}
	return S
}

// DiscriminativeGD performs discriminative mini-batch gradient descent on SPN S given data D.
// Argument eta is the learning rate, eps is the convergence difference in likelihood, D is the
// dataset, signals whether to normalize weights at each update and b is the size of the
// mini-batch.
func DiscriminativeBGD(S spn.SPN, eta, eps float64, D spn.Dataset, Y []*Variable, norm bool, b int) spn.SPN {
	st := spn.NewStorer()
	Q := &common.Queue{}
	s, d, w := st.NewTicket(), st.NewTicket(), st.NewTicket()
	z, p, u := st.NewTicket(), st.NewTicket(), st.NewTicket()
	l := st.NewTicket()
	P := S.Parameters()
	y := make([]int, len(Y))
	var j int
	for i := 0; i < P.Iterations; i++ {
		for _, I := range D {
			spn.StoreInference(S, I, s, st)
			DeriveSPN(S, st, d, s, Q)
			DeriveWeights(S, st, w, d, s, Q)
			pullValues(I, Y, y)
			spn.StoreInference(S, I, z, st)
			DeriveSPN(S, st, p, z, Q)
			DeriveWeights(S, st, u, p, z, Q)
			pushValues(I, Y, y)
			storeDGD(S, s, d, w, z, p, u, l, st, eta, norm, Q)
			st.ResetTickets(s, d, w, z, p, u)
			j++
			if j%b == 0 {
				applyDGDFrom(S, l, st, eta, norm)
				st.Reset(l)
			}
		}
	}
	if j%b != 0 {
		applyDGDFrom(S, l, st, eta, norm)
		st.Reset(l)
	}
	return S
}

// DiscriminativeGD performs hard (MPE) discriminative gradient descent on SPN S given data D.
// Argument eta is the learning rate, eps is the convergence difference in likelihood, D is the
// dataset, norm signals whether to normalize weights at each update and b is the size of the
// mini-batch.
func DiscriminativeHardBGD(S spn.SPN, eta, eps float64, D spn.Dataset, Y []*Variable, norm bool, b int) spn.SPN {
	st := spn.NewStorer()
	d, p := st.NewTicket(), st.NewTicket()
	P := S.Parameters()
	y := make([]int, len(Y))
	var j int
	for i := 0; i < P.Iterations; i++ {
		for _, I := range D {
			DeriveHard(S, st, d, I)
			pullValues(I, Y, y)
			DeriveHard(S, st, p, I)
			pushValues(I, Y, y)
			applyHDGD(S, d, p, st, eta, norm)
			j++
			if j%b == 0 {
				applyHDGD(S, d, p, st, eta, norm)
				st.ResetTickets(d, p)
			}
		}
	}
	if j%b != 0 {
		applyHDGD(S, d, p, st, eta, norm)
		st.ResetTickets(d, p)
	}
	return S
}

// applyDGD computes the discriminative derivative
//  (1/S(Y|X))*dS/dW(Y|X)-(1/S(1|X))*dS/dW(1|X)
// and applies to each weight.
// Argument S is the SPN to be derived, integers s, d, w, z, p and u are the tickets for S(Y|X),
// dSn/dSj(Y|X), dS/dW(Y|X), S(1|X), dSn/dSj(1|W), dS/dW(1|W) respectively on *spn.Storer st.
func applyDGD(S spn.SPN, s, d, w, z, p, u int, st *spn.Storer, eta float64, norm bool, Q *common.Queue) {
	Q.Reset()
	Q.Enqueue(S)
	V := make(map[spn.SPN]bool)
	V[S] = true
	dwt, _ := st.Table(w)
	dut, _ := st.Table(u)
	ist, _ := st.Table(s)
	izt, _ := st.Table(z)
	P := S.Parameters()
	for !Q.Empty() {
		n := Q.Dequeue().(spn.SPN)
		ch := n.Ch()
		if t := n.Type(); t == "sum" {
			sum := n.(*spn.Sum)
			W := sum.Weights()
			dW, _ := dwt.Value(n)
			dU, _ := dut.Value(n)
			iS, _ := ist.Value(n)
			iZ, _ := izt.Value(n)
			for i := range W {
				delta := eta * (math.Exp(dW[i]-iS[i]) - math.Exp(dU[i]-iZ[i]))
				W[i] += delta - 2*P.Lambda*W[i]
			}
			if norm {
				Normalize(W)
			}
		}
		for _, cs := range ch {
			if cs.Type() != "leaf" && V[cs] {
				Q.Enqueue(cs)
				V[cs] = true
			}
		}
	}
}

// storeDGD computes the discriminative gradient, but instead of applying to the weights directly,
// the function stores the update into ticket l, summing previous iteration values.
func storeDGD(S spn.SPN, s, d, w, z, p, u, l int, st *spn.Storer, eta float64, norm bool, Q *common.Queue) {
	Q.Enqueue(S)
	V := make(map[spn.SPN]bool)
	V[S] = true
	dwt, _ := st.Table(w)
	dut, _ := st.Table(u)
	ist, _ := st.Table(s)
	izt, _ := st.Table(z)
	lt, _ := st.Table(l)
	for !Q.Empty() {
		n := Q.Dequeue().(spn.SPN)
		ch := n.Ch()
		if t := n.Type(); t == "sum" {
			sum := n.(*spn.Sum)
			nw := len(sum.Weights())
			dW, _ := dwt.Value(n)
			dU, _ := dut.Value(n)
			iS, _ := ist.Value(n)
			iZ, _ := izt.Value(n)
			for i := 0; i < nw; i++ {
				delta := math.Exp(dW[i]-iS[i]) - math.Exp(dU[i]-iZ[i])
				v, _ := lt.Entry(n, i)
				lt.Store(n, i, v+delta)
			}
		}
		for _, cs := range ch {
			if cs.Type() != "leaf" && V[cs] {
				Q.Enqueue(cs)
				V[cs] = true
			}
		}
	}
}

func applyDGDFrom(S spn.SPN, l int, st *spn.Storer, eta float64, norm bool) {
	T, _ := st.Table(l)
	P := S.Parameters()
	for s, dW := range T {
		if s.Type() == "sum" {
			W := s.(*spn.Sum).Weights()
			for i, d := range dW {
				delta := eta * d
				W[i] += delta - 2*P.Lambda*W[i]
			}
			if norm {
				Normalize(W)
			}
		}
	}
}

func applyHDGD(S spn.SPN, d, p int, st *spn.Storer, eta float64, norm bool) {
	dt, _ := st.Table(d)
	pt, _ := st.Table(p)
	C := make(map[spn.SPN]map[int]float64)
	for s, cnts := range dt {
		if _, e := C[s]; !e {
			C[s] = make(map[int]float64)
		}
		for i, c := range cnts {
			C[s][i] = c
		}
	}
	for s, cnts := range pt {
		if _, e := C[s]; !e {
			C[s] = make(map[int]float64)
		}
		for i, c := range cnts {
			C[s][i] = C[s][i] - c
		}
	}
	for s, cnts := range C {
		// DeriveHard guarantees s is sum.
		sum := s.(*spn.Sum)
		W := sum.Weights()
		for i, delta := range cnts {
			w := W[i]
			W[i] += eta * (delta / w)
		}
		if norm {
			Normalize(W)
		}
	}
}
