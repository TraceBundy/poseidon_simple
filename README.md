# poseidon


Reference implementation for the Poseidon Hashing algorithm.


#### Reference

[Starkad and Poseidon: New Hash Functions for Zero Knowledge Proof Systems](https://eprint.iacr.org/2019/458.pdf)

This repository has been created so there's a unique library that holds the tools & functions
required to perform Poseidon Hashes.

## Usage

```bigquery
func main()  {
	inputs := []*ff.Scalar{ff.NewInt(0), ff.NewInt(1)}
	p, err := NewRate(len(inputs))
    if err != nil {
        return
    }
	for _, i := range inputs {
		p.Write(i.Bytes())
	}
	hash := p.Sum(nil)
	fmt.Println(hex.EncodeToString(hash))
}
```
```

## Licensing

This code is licensed under the GNU Lesser General Public License v3.0. Please see [LICENSE](https://github.com/PlatONnetwork/PlatON-Go/blob/develop/COPYING) for further info.