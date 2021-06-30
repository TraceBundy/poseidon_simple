package constants

import (
	"github.com/PlatONnetwork/poseidon/ff"
)

type Constants struct {
	C [][]*ff.Scalar
	M [][][]*ff.Scalar
}

type Curve int

const (
	Bn256 Curve = iota + 1
)


func GetConstants(curve Curve) *Constants {
	switch curve {
	case Bn256:
		return C
	}
	return nil
}
