package poseidon_simple

import (
	"errors"
	"github.com/PlatONnetwork/poseidon/constants"
	"github.com/PlatONnetwork/poseidon/ff"
)

const NRoundsF = 8

var NRoundsP = []int{56, 57, 56, 60, 60, 63, 64, 63}

type Poseidon struct {
	params *Params
	input  []*ff.Scalar
}

type Params struct {
	Width     int
	RF        int
	RP        int
	RoundKeys []*ff.Scalar
	MDSMatrix [][]*ff.Scalar
	Modulus   *ff.Scalar
}

func New() *Poseidon {
	return NewRate(2)
}
func NewRate(rate int) *Poseidon {
	width := rate + 1
	rf, rp := roundNumber(width)
	c := constants.GetConstants(constants.Bn256)
	return &Poseidon{
		params: &Params{
			Width:     width,
			RF:        rf,
			RP:        rp,
			RoundKeys: c.C[width-2],
			MDSMatrix: c.M[width-2],
			Modulus:   constants.Modulus,
		},
		input: []*ff.Scalar{ff.NewInt(0)},
	}
}

func roundNumber(width int) (int, int) {
	return NRoundsF, NRoundsP[width-2]
}

func (p *Poseidon) Write(x []byte) (n int, err error) {
	if p.params.Width == len(x) {
		return 0, errors.New("maximum width reached")
	}
	p.input = append(p.input, ff.NewBytes(x))
	return len(x), nil
}

func (p *Poseidon) Reset() {
	p.input = []*ff.Scalar{ff.NewInt(0)}
}

func (p Poseidon) Size() int {
	return p.params.Width
}

func (p Poseidon) BlockSize() int {
	panic("implement me")
}

func (p *Poseidon) Sum(in []byte) []byte {
	p.Pad()
	newState := make([]*ff.Scalar, p.params.Width)
	for i := 0; i < p.params.Width; i++ {
		newState[i] = ff.NewInt(0)
	}

	// ARK --> SBox --> M, https://eprint.iacr.org/2019/458.pdf pag.5
	for i := 0; i < p.params.RF + p.params.RP; i++ {
		ark(p.input, p.params.RoundKeys, i*p.params.Width, p.params.Modulus)
		sbox(p.params.RF, p.params.RP, p.input, i, p.params.Modulus)
		mix(p.input, newState, p.params.MDSMatrix, p.params.Modulus)
		p.input, newState = newState, p.input
	}
	r := p.input[0]
	return r.Bytes()
}

// ark computes Add-Round Key, from the paper https://eprint.iacr.org/2019/458.pdf
func ark(state []*ff.Scalar, c []*ff.Scalar, it int, modulus *ff.Scalar) {
	for i := 0; i < len(state); i++ {
		state[i].AddMod(state[i], c[it+i], modulus)
	}
}

// exp5 performs x^5 mod p
// https://eprint.iacr.org/2019/458.pdf page 8
func exp5(a *ff.Scalar, modulus *ff.Scalar) {
	a.ExpMod(a, ff.NewInt(5), modulus)
}

// sbox https://eprint.iacr.org/2019/458.pdf page 6
func sbox(nRoundsF, nRoundsP int, state []*ff.Scalar, i int, modulus *ff.Scalar) {
	if (i < nRoundsF/2) || (i >= nRoundsF/2+nRoundsP) {
		for j := 0; j < len(state); j++ {
			exp5(state[j], modulus)
		}
	} else {
		exp5(state[0], modulus)
	}
}

// mix returns [[matrix]] * [vector]
func mix(state []*ff.Scalar, newState []*ff.Scalar, m [][]*ff.Scalar, modulus *ff.Scalar) {
	mul := ff.NewInt(0)
	for i := 0; i < len(state); i++ {
		newState[i].SetUint64(0)
		for j := 0; j < len(state); j++ {
			mul.MulMod(m[i][j], state[j], modulus)
			newState[i].AddMod(newState[i], mul, modulus)
		}
	}
}
//// Sum will compute the Poseidon digest value. The usage of the bytes parameter is currently not implemented
//func (p *Poseidon) Sum(in []byte) []byte {
//	p.Pad()
//
//	keysOffset := 0
//
//	halfOfFullRound := p.params.RF / 2
//	for i := 0; i < halfOfFullRound; i++ {
//		p.applyFullRound(&keysOffset)
//	}
//
//	for i := 0; i < p.params.RP; i++ {
//		p.applyPartialRound(&keysOffset)
//	}
//
//	for i := 0; i < halfOfFullRound; i++ {
//		p.applyFullRound(&keysOffset)
//	}
//
//	return p.input[1].Bytes()
//}

//func (p *Poseidon) applyFullRound(keysOffset *int) {
//	// Add current round constant to all elements of input
//	for i := 0; i < len(p.input); i++ {
//		p.input[i] = p.input[i].AddMod(p.input[i], p.params.RoundKeys[*keysOffset], p.params.Modulus)
//		*keysOffset++
//	}
//
//	// Apply quintic SBox to every element
//	for i := 0; i < len(p.input); i++ {
//		QuinticSbox(p.input[i], p.params.Modulus)
//	}
//
//	p.input = mulVec(p.params.MDSMatrix, p.input, p.params.Modulus)
//}
//
//func (p *Poseidon) applyPartialRound(keysOffset *int) {
//	// Add current round constant to all elements of input
//	for i := 0; i < len(p.input); i++ {
//		p.input[i] = p.input[i].AddMod(p.input[i], p.params.RoundKeys[*keysOffset], p.params.Modulus)
//		*keysOffset++
//	}
//
//	// Apply quintic SBox to the first element
//	QuinticSbox(p.input[0], p.params.Modulus)
//
//	p.input = mulVec(p.params.MDSMatrix, p.input, p.params.Modulus)
//}
//
//func mulVec(a [][]*ff.Scalar, b []*ff.Scalar, modulus *ff.Scalar) []*ff.Scalar {
//	result := make([]*ff.Scalar, len(b))
//
//	for j := 0; j < len(a); j++ {
//		line := make([]*ff.Scalar, len(b))
//
//		for k := 0; k < len((a)[j]); k++ {
//			line[k] = ff.NewInt(0)
//			line[k].MulMod(a[j][k], b[k], modulus)
//		}
//
//		for k := 0; k < len(line); k++ {
//			result[j] = ff.NewInt(0)
//			result[j].AddMod(result[j], line[k], modulus)
//		}
//	}
//
//	return result
//}
//
//// QuinticSbox will set *a to a^5
//func QuinticSbox(a *ff.Scalar, modulus *ff.Scalar) {
//	//a.SetBytes(ff.ExpMod(a, ff.NewInt(5), modulus).Bytes())
//	c := a.Clone()
//	for k := 0; k < 4; k++ {
//		a.MulMod(a, c, modulus)
//	}
//}
//
// Pad will fill the input with zeroed scalars until its length equal the parametrization width
func (p *Poseidon) Pad() {
	dif := p.params.Width - len(p.input)
	if dif > 0 {
		pad := make([]*ff.Scalar, dif)
		for i := 0; i < dif; i++ {
			pad[i] = ff.NewInt(0)
		}
		p.input = append(p.input, pad...)
	}
}
