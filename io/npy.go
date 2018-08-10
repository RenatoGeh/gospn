package io

import (
	"errors"
	"github.com/RenatoGeh/npyio"
	"os"
)

var (
	validTypes = [39]string{
		"u1", "<u1", "|u1", "uint8",
		"u2", "<u2", "|u2", ">u2", "uint16",
		"u4", "<u4", "|u4", ">u4", "uint32",
		"u8", "<u8", "|u8", ">u8", "uint64",
		"i1", "<i1", "|i1", ">i1", "int8",
		"i2", "<i2", "|i2", ">i2", "int16",
		"i4", "<i4", "|i4", ">i4", "int32",
		"i8", "<i8", "|i8", ">i8", "int64",
	}

	ErrNonIntegerType  = errors.New("gospn: npy data type is non integer.")
	ErrNonDatasetShape = errors.New("gospn: npy data does not have dimension two.")
)

func integerType(t int) int {
	if t < 4 {
		return 0
	}
	return (t + 1) / 5
}

func isValid(t string) int {
	for i, v := range validTypes {
		if t == v {
			return integerType(i)
		}
	}
	return -1
}

// NpyReader is a .npy reader. GoSPN supports only integer data for now.
type NpyReader struct {
	f *os.File      // File stream.
	r *npyio.Reader // Npyio reader.
	p int           // Position in number of elements (not instances!).
	t int           // Actual type.
	s [2]int        // Shape.
}

// NewNpyReader creates a new *NpyReader from .npy file fname.
func NewNpyReader(fname string) (*NpyReader, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	r, err := npyio.NewReader(f)
	if err != nil {
		return nil, err
	}
	s := r.Header.Descr.Shape
	if len(s) != 2 {
		return nil, ErrNonDatasetShape
	}
	t := isValid(r.Header.Descr.Type)
	if t < 0 {
		return nil, ErrNonIntegerType
	}
	return &NpyReader{f, r, 0, t, [2]int{s[0], s[1]}}, nil
}

func read(r *NpyReader, n int) ([]map[int]int, []int, error) {
	k := n * r.s[1]
	u := make([]map[int]int, n)
	l := make([]int, n)
	y := r.s[1]
	switch r.t {
	case 0:
		m := make([]uint8, k)
		err := r.r.Read(&m)
		if err != nil {
			return nil, nil, err
		}
		for i, v := range m {
			p := i / y
			if u[p] == nil {
				u[p] = make(map[int]int)
			}
			j := i % y
			if j == r.s[1]-1 {
				l[p] = int(v)
			} else {
				u[p][j] = int(v)
			}
		}
	case 1:
		m := make([]uint16, k)
		err := r.r.Read(&m)
		if err != nil {
			return nil, nil, err
		}
		for i, v := range m {
			p := i / y
			if u[p] == nil {
				u[p] = make(map[int]int)
			}
			j := i % y
			if j == r.s[1]-1 {
				l[p] = int(v)
			} else {
				u[p][j] = int(v)
			}
		}
	case 2:
		m := make([]uint32, k)
		err := r.r.Read(&m)
		if err != nil {
			return nil, nil, err
		}
		for i, v := range m {
			p := i / y
			if u[p] == nil {
				u[p] = make(map[int]int)
			}
			j := i % y
			if j == r.s[1]-1 {
				l[p] = int(v)
			} else {
				u[p][j] = int(v)
			}
		}
	case 3:
		m := make([]uint64, k)
		err := r.r.Read(&m)
		if err != nil {
			return nil, nil, err
		}
		for i, v := range m {
			p := i / y
			if u[p] == nil {
				u[p] = make(map[int]int)
			}
			j := i % y
			if j == r.s[1]-1 {
				l[p] = int(v)
			} else {
				u[p][j] = int(v)
			}
		}
	case 4:
		m := make([]int8, k)
		err := r.r.Read(&m)
		if err != nil {
			return nil, nil, err
		}
		for i, v := range m {
			p := i / y
			if u[p] == nil {
				u[p] = make(map[int]int)
			}
			j := i % y
			if j == r.s[1]-1 {
				l[p] = int(v)
			} else {
				u[p][j] = int(v)
			}
		}
	case 5:
		m := make([]int16, k)
		err := r.r.Read(&m)
		if err != nil {
			return nil, nil, err
		}
		for i, v := range m {
			p := i / y
			if u[p] == nil {
				u[p] = make(map[int]int)
			}
			j := i % y
			if j == r.s[1]-1 {
				l[p] = int(v)
			} else {
				u[p][j] = int(v)
			}
		}
	case 6:
		m := make([]int32, k)
		err := r.r.Read(&m)
		if err != nil {
			return nil, nil, err
		}
		for i, v := range m {
			p := i / y
			if u[p] == nil {
				u[p] = make(map[int]int)
			}
			j := i % y
			if j == r.s[1]-1 {
				l[p] = int(v)
			} else {
				u[p][j] = int(v)
			}
		}
	case 7:
		m := make([]int64, k)
		err := r.r.Read(&m)
		if err != nil {
			return nil, nil, err
		}
		for i, v := range m {
			p := i / y
			if u[p] == nil {
				u[p] = make(map[int]int)
			}
			j := i % y
			if j == r.s[1]-1 {
				l[p] = int(v)
			} else {
				u[p][j] = int(v)
			}
		}
	}
	return u, l, nil
}

// Read reads n instances from file and returns a dataset and label slice.
func (r *NpyReader) Read(n int) ([]map[int]int, []int, error) {
	return read(r, n)
}

// ReadAll reads all instances from file and returns a dataset and label slice.
func (r *NpyReader) ReadAll() ([]map[int]int, []int, error) {
	return read(r, r.s[0])
}

// Close closes this reader's stream.
func (r *NpyReader) Close() {
	r.f.Close()
}
