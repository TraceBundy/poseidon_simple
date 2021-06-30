package ff

import (
	"fmt"
	"math/big"
	"strconv"
)

type Scalar struct {
	*big.Int
}

func NewString(s string) (*Scalar, bool) {
	l := &Scalar{new(big.Int)}
	i, res := new(big.Int).SetString(s, 10)
	if !res {
		return nil, res
	}
	l.SetBytes(i.Bytes())
	return l, res
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

// ExpMod computes z = (x ^ y) % p.
func (s *Scalar) ExpMod(x *Scalar, y *Scalar, p *Scalar) *Scalar {
	s.Int.Exp(x.Int, y.Int, p.Int)
	return s
}

func (s *Scalar) InField(modulus *Scalar)bool  {
	return s.Cmp(modulus.Int) == -1
}

func (s *Scalar) Bytes()[]byte  {
	var buf [32]byte
	b := s.Int.Bytes()
	copy(buf[32-len(b):], b)
	return buf[:]
}

func (s *Scalar) MarshalJSON() ([]byte, error)  {
	if s == nil || s.Int == nil {
		return []byte(fmt.Sprintf(`"%s"`,"null")), nil
	}
	return []byte(fmt.Sprintf(`"%s"`,s.String())), nil
}

func (s *Scalar) UnmarshalJSON(data []byte) error {
	num , err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	if num == "null" {
		return nil
	}
	n, res := new(big.Int).SetString(num, 10)
	if !res {
		return fmt.Errorf("not a decimal number")
	}
	s.Int = n
	return nil
}