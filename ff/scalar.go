package ff

import (
	"math/big"
)

type Scalar struct {
	*big.Int
}

func NewString(s string) *Scalar {
	l := &Scalar{new(big.Int)}
	i, _ := new(big.Int).SetString(s, 10)
	l.SetBytes(i.Bytes())
	return l
}
func NewBytes(x []byte) *Scalar {
	l := &Scalar{new(big.Int)}
	l.SetBytes(x)
	return l
}

func NewInt(x uint64) *Scalar {
	l := &Scalar{new(big.Int)}
	l.SetUint64(x)
	return l
}

func NewBigInt(x *big.Int) *Scalar {
	l := &Scalar{new(big.Int).Set(x)}
	return l
}


func (s Scalar) Clone() *Scalar {
	return &Scalar{
		new(big.Int).Set(s.Int),
	}
}
func (s *Scalar) Rsh(r *Scalar, i uint) *Scalar {
	s.Int.Rsh(r.Int, i)
	return s
}

func (s *Scalar) Cmp(c *Scalar) int {
	return s.Int.Cmp(c.Int)
}
func (s *Scalar) ModInverse(n *Scalar, p *Scalar) *Scalar {
	s.Int.ModInverse(n.Int, p.Int)
	return s
}

// AddMod computes z = (x + y) % p.
func (s *Scalar) AddMod(x *Scalar, y *Scalar, p *Scalar) *Scalar {
	s.Add(x.Int, y.Int)
	s.Mod(s.Int, p.Int)
	return s
}

// SubMod computes z = (x - y) % p.
func (s *Scalar) SubMod(x *Scalar, y *Scalar, p *Scalar) *Scalar {
	s.Sub(x.Int, y.Int)
	s.Mod(s.Int, p.Int)
	return s
}

// MulMod computes z = (x * y) % p.
func (s *Scalar) MulMod(x *Scalar, y *Scalar, p *Scalar) *Scalar {
	n := &Scalar{new(big.Int).Set(x.Int)}
	s.Int = big.NewInt(0)

	for i := 0; i < y.BitLen(); i++ {
		if y.Bit(i) == 1 {
			s.AddMod(s, n, p)
		}
		n.AddMod(n, n, p)
	}

	return s
}
func (s *Scalar)ExpMod(x *Scalar, y *Scalar, p *Scalar) *Scalar  {
	s.Int.Exp(s.Int, y.Int, p.Int)
	return s
}

//// invMod computes z = (1/x) % p.
//func invMod(x *big.Int, p *big.Int) (z *big.Int) {
//	z = new(big.Int).ModInverse(x, p)
//	return z
//}

// ExpMod computes z = (x^e) % p.
func ExpMod(x *Scalar, y *Scalar, p *Scalar) *Scalar {
	z := new(big.Int).Exp(x.Int, y.Int, p.Int)
	return &Scalar{
		z,
	}
}

// SqrtMod computes z = sqrt(x) % p.
func SqrtMod(x *Scalar, p *Scalar) (z *Scalar) {
	/* assert that p % 4 == 3 */
	if new(big.Int).Mod(p.Int, big.NewInt(4)).Cmp(big.NewInt(3)) != 0 {
		panic("p is not equal to 3 mod 4!")
	}

	/* z = sqrt(x) % p = x^((p+1)/4) % p */

	/* e = (p+1)/4 */
	e := new(big.Int).Add(p.Int, big.NewInt(1))
	e = e.Rsh(e, 2)

	z = ExpMod(x, &Scalar{e}, p)
	return z
}

func (s *Scalar) ModSqrt(x *Scalar, p *Scalar) *Scalar {
	s.Int.ModSqrt(x.Int, p.Int)
	return s
}
