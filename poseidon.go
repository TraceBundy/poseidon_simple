package poseidon

import (
	"errors"
	"fmt"
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

func New() (*Poseidon, error) {
	return NewRate(2)
}

func NewRate(rate int) (*Poseidon, error) {
	if rate == 0 || rate >= len(NRoundsP)-1 {
		return nil, errors.New("invalid rate")
	}
	width := rate + 1
	rf, rp, err := roundNumber(width)
	if err != nil {
		return nil, err
	}
	c := constants.GetConstants(constants.Bn256)
	return NewWithParams(
		&Params{
			Width:     width,
			RF:        rf,
			RP:        rp,
			RoundKeys: c.C[width-2],
			MDSMatrix: c.M[width-2],
			Modulus:   constants.Modulus,
		}), nil
}
func NewWithParams(params *Params) *Poseidon {
	return &Poseidon{
		params: params,
		input:  []*ff.Scalar{ff.NewInt(0)},
	}
}

func roundNumber(width int) (int, int, error) {
	if width > len(NRoundsP)+1 {
		return 0, 0, fmt.Errorf("invalid width")
	}
	return NRoundsF, NRoundsP[width-2], nil
}

func (p *Poseidon) Write(x []byte) (n int, err error) {
	if p.params.Width == len(x) {
		return 0, errors.New("maximum width reached")
	}
	element := ff.NewBytes(x)
	if !element.InField(p.params.Modulus) {
		return 0, errors.New("invalid value not inside Finite Field")
	}
	p.input = append(p.input, element)
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

// Sum will compute the Poseidon digest value. The usage of the bytes parameter is currently not implemented
func (p *Poseidon) Sum(in []byte) []byte {
	p.pad()

	keysOffset := 0

	halfOfFullRound := p.params.RF / 2
	for i := 0; i < halfOfFullRound; i++ {
		p.applyFullRound(&keysOffset)
	}

	for i := 0; i < p.params.RP; i++ {
		p.applyPartialRound(&keysOffset)
	}

	for i := 0; i < halfOfFullRound; i++ {
		p.applyFullRound(&keysOffset)
	}

	return p.input[0].Bytes()
}

// pad will fill the input with zeroed scalars until its length equal the parametrization width
func (p *Poseidon) pad() {
	dif := p.params.Width - len(p.input)
	if dif > 0 {
		pad := makeArray(dif, ff.NewInt(0))
		p.input = append(p.input, pad...)
	}
}

func (p *Poseidon) applyFullRound(keysOffset *int) {
	// Add current round constant to all elements of input
	for i := 0; i < len(p.input); i++ {
		p.input[i] = p.input[i].AddMod(p.input[i], p.params.RoundKeys[*keysOffset], p.params.Modulus)
		*keysOffset++
	}

	// Apply quintic SBox to every element
	for i := 0; i < len(p.input); i++ {
		quinticSBox(p.input[i], p.params.Modulus)
	}

	p.input = mulVec(p.params.MDSMatrix, p.input, p.params.Modulus)
}

func (p *Poseidon) applyPartialRound(keysOffset *int) {
	// Add current round constant to all elements of input
	for i := 0; i < len(p.input); i++ {
		p.input[i] = p.input[i].AddMod(p.input[i], p.params.RoundKeys[*keysOffset], p.params.Modulus)
		*keysOffset++
	}

	// Apply quintic SBox to the first element
	quinticSBox(p.input[0], p.params.Modulus)

	p.input = mulVec(p.params.MDSMatrix, p.input, p.params.Modulus)
}

func makeArray(len int, value *ff.Scalar) []*ff.Scalar {
	result := make([]*ff.Scalar, len)
	for i := 0; i < len; i++ {
		result[i] = value.Clone()
	}
	return result
}

func mulVec(a [][]*ff.Scalar, b []*ff.Scalar, modulus *ff.Scalar) []*ff.Scalar {
	result := makeArray(len(b), ff.NewInt(0))

	for j := 0; j < len(a); j++ {
		line := makeArray(len(b), ff.NewInt(0))
		for k := 0; k < len((a)[j]); k++ {
			line[k].MulMod(a[j][k], b[k], modulus)
		}

		for k := 0; k < len(line); k++ {
			result[j].AddMod(result[j], line[k], modulus)
		}
	}

	return result
}

// quinticSBox will set *a to a^5
func quinticSBox(a *ff.Scalar, modulus *ff.Scalar) {
	a.ExpMod(a, ff.NewInt(5), modulus)
}
