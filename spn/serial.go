package spn

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

func init() {
	gob.Register(&Sum{})
	gob.Register(&Product{})
	gob.Register(&Gaussian{})
	gob.Register(&Multinomial{})
	gob.Register([][]uint32{})
}

// RegisterGobType registers a new SPN type to be marshalled. If you want to create a new SPN type
// and be able to serialize it, you must implement interfaces GobEncoder and GobDecoder as well as
// call this function with a pointer to a concrete type (e.g. RegisterGobType(&NewSPNType{})). All
// basic SPN nodes are already registered.
func RegisterGobType(t interface{}) {
	gob.Register(t)
}

func encodeSPN(enc *gob.Encoder, S SPN) {
	if err := enc.Encode(&S); err != nil {
		panic(err)
	}
}

func decodeSPN(dec *gob.Decoder) SPN {
	var S SPN
	if err := dec.Decode(&S); err != nil {
		panic(err)
	}
	return S
}

func Marshal(S SPN) []byte {
	var net bytes.Buffer
	enc := gob.NewEncoder(&net)

	M := make(map[SPN]uint32)
	var list [][]uint32
	var n uint32
	TopSortTarjanFunc(S, nil, func(Z SPN) bool {
		list = append(list, []uint32{})
		M[Z] = n
		// Invariant: because of topological order, any child of Z has already been visited.
		for _, c := range Z.Ch() {
			p := M[c]
			list[n] = append(list[n], p)
		}
		enc.Encode(n)
		encodeSPN(enc, Z)
		n++
		return true
	})
	enc.Encode(list)

	nb := make([]byte, 4)
	binary.LittleEndian.PutUint32(nb, n)

	return append(nb, net.Bytes()...)
}

func Unmarshal(buffer []byte) SPN {
	net := bytes.NewBuffer(buffer[4:])
	dec := gob.NewDecoder(net)

	n := binary.LittleEndian.Uint32(buffer[:4])
	M := make([]SPN, n)
	for i := uint32(0); i < n; i++ {
		var k uint32
		dec.Decode(&k)
		Z := decodeSPN(dec)
		M[k] = Z
	}

	var list [][]uint32
	dec.Decode(&list)
	for i, ch := range list {
		S := M[i]
		for _, c := range ch {
			Z := M[c]
			S.AddChild(Z)
		}
	}

	// Invariant: topological sort guarantees last guy is root
	return M[n-1]
}
