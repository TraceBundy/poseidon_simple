package poseidon

import (
	"github.com/PlatONnetwork/poseidon/ff"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)
func TestPoseidon(t *testing.T)  {
	b0 := big.NewInt(0)
	b1 := big.NewInt(1)
	b2 := big.NewInt(2)
	Hash := func(inputs []*big.Int) (*big.Int, error) {
		p, _ := NewRate(len(inputs))
		for _, i := range inputs {
			p.Write(i.Bytes())
		}
		return new(big.Int).SetBytes(p.Sum(nil)), nil
	}

	h, err := Hash([]*big.Int{b1})
	assert.Nil(t, err)
	assert.Equal(t,
		"18586133768512220936620570745912940619677854269274689475585506675881198879027",
		h.String())

	h, err = Hash([]*big.Int{b1, b2})
	assert.Nil(t, err)
	assert.Equal(t,
		"7853200120776062878684798364095072458815029376092732009249414926327459813530",
		h.String())


	h, err = Hash([]*big.Int{b1, b2, b0, b0, b0})
	assert.Nil(t, err)
	assert.Equal(t,
		"1018317224307729531995786483840663576608797660851238720571059489595066344487",
		h.String())

	h, err = Hash([]*big.Int{b1, b2, b0, b0, b0, b0})
	assert.Nil(t, err)
	assert.Equal(t,
		"15336558801450556532856248569924170992202208561737609669134139141992924267169",
		h.String())

	b3 := big.NewInt(3)
	b4 := big.NewInt(4)
	h, err = Hash([]*big.Int{b3, b4, b0, b0, b0})
	assert.Nil(t, err)
	assert.Equal(t,
		"5811595552068139067952687508729883632420015185677766880877743348592482390548",
		h.String())
	h, err = Hash([]*big.Int{b3, b4, b0, b0, b0, b0})
	assert.Nil(t, err)
	assert.Equal(t,
		"12263118664590987767234828103155242843640892839966517009184493198782366909018",
		h.String())

	b5 := big.NewInt(5)
	b6 := big.NewInt(6)
	h, err = Hash([]*big.Int{b1, b2, b3, b4, b5, b6})
	assert.Nil(t, err)
	assert.Equal(t,
		"20400040500897583745843009878988256314335038853985262692600694741116813247201",
		h.String())
}

func TestPoseidonError(t *testing.T) {
	_, err := NewRate(0)
	assert.NotNil(t, err)
	_, err = NewRate(6)
	assert.Nil(t, err)
	_, err = NewRate(7)
	assert.NotNil(t, err)
}

func BenchmarkPoseidonHash(b *testing.B) {
	b0 := ff.NewInt(0)
	b1, _ := ff.NewString("12242166908188651009877250812424843524687801523336557272219921456462821518061")
	b2, _ := ff.NewString("12242166908188651009877250812424843524687801523336557272219921456462821518061")

	bigArray4 := [][]byte{b1.Bytes(), b2.Bytes(), b0.Bytes(), b0.Bytes(), b0.Bytes(), b0.Bytes()}
	for i := 0; i < b.N; i++ {
		p, _ := NewRate(len(bigArray4))
		for _, b := range bigArray4 {
			p.Write(b)
		}
		p.Sum(nil)
	}
}