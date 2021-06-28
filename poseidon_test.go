package poseidon_simple

import (
	"encoding/hex"
	"fmt"
	"github.com/PlatONnetwork/poseidon/ff"
	"testing"
)

func TestPoseidon(t *testing.T)  {
	p := New()
	p.Write(ff.NewInt(1).Bytes())
	//p.Write(ff.NewInt(2).Bytes())

	fmt.Println(hex.EncodeToString(p.Sum(nil)))
}
//1988bd3180b84efed5c17d5aab80edc5571fc0cb8974fd83cd816517bcd82f7e