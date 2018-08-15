package io

import (
	"errors"
	"github.com/sbinet/npyio"
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
	return &NpyReader{f, r, t, [2]int{s[0], s[1]}}, nil
}

func new_counts(n int, k int) []int {
	c := make([]int, k)
	q, p := n/k, n%k
	for i := 0; i < k; i++ {
		if i < p {
			c[i] = q + 1
		} else {
			c[i] = q
		}
	}
	return c
}

// readBalanced reads .npy file from an *NpyReader, pulling n entries total. It pulls one instance
// at a time, attempting to return a balanced dataset. That is, there should be roughly the same
// number of entries for each label in the resulting dataset. Argument c is the number of classes
// in the dataset.
func readBalanced(r *NpyReader, n, c int) ([]map[int]int, []int, error) {
	D := make([]map[int]int, n)
	L := make([]int, n)
	y := r.s[1]
	C := new_counts(n, c)
	for i := 0; i < n; {
		switch r.t {
		case 0:
			m := make([]uint8, y)
			err := r.r.Read(&m)
			if err != nil {
				return nil, nil, err
			}
			if l := m[y-1]; C[l] > 0 {
				D[i] = make(map[int]int)
				for j := 0; j < y-1; j++ {
					D[i][j] = int(m[j])
				}
				L[i] = int(l)
				C[l]--
				i++
			}
		case 1:
			m := make([]uint16, y)
			err := r.r.Read(&m)
			if err != nil {
				return nil, nil, err
			}
			if l := m[y-1]; C[l] > 0 {
				D[i] = make(map[int]int)
				for j := 0; j < y-1; j++ {
					D[i][j] = int(m[j])
				}
				L[i] = int(l)
				C[l]--
				i++
			}
		case 2:
			m := make([]uint32, y)
			err := r.r.Read(&m)
			if err != nil {
				return nil, nil, err
			}
			if l := m[y-1]; C[l] > 0 {
				D[i] = make(map[int]int)
				for j := 0; j < y-1; j++ {
					D[i][j] = int(m[j])
				}
				L[i] = int(l)
				C[l]--
				i++
			}
		case 3:
			m := make([]uint64, y)
			err := r.r.Read(&m)
			if err != nil {
				return nil, nil, err
			}
			if l := m[y-1]; C[l] > 0 {
				D[i] = make(map[int]int)
				for j := 0; j < y-1; j++ {
					D[i][j] = int(m[j])
				}
				L[i] = int(l)
				C[l]--
				i++
			}
		case 4:
			m := make([]int8, y)
			err := r.r.Read(&m)
			if err != nil {
				return nil, nil, err
			}
			if l := m[y-1]; C[l] > 0 {
				D[i] = make(map[int]int)
				for j := 0; j < y-1; j++ {
					D[i][j] = int(m[j])
				}
				L[i] = int(l)
				C[l]--
				i++
			}
		case 5:
			m := make([]int16, y)
			err := r.r.Read(&m)
			if err != nil {
				return nil, nil, err
			}
			if l := m[y-1]; C[l] > 0 {
				D[i] = make(map[int]int)
				for j := 0; j < y-1; j++ {
					D[i][j] = int(m[j])
				}
				L[i] = int(l)
				C[l]--
				i++
			}
		case 6:
			m := make([]int32, y)
			err := r.r.Read(&m)
			if err != nil {
				return nil, nil, err
			}
			if l := m[y-1]; C[l] > 0 {
				D[i] = make(map[int]int)
				for j := 0; j < y-1; j++ {
					D[i][j] = int(m[j])
				}
				L[i] = int(l)
				C[l]--
				i++
			}
		case 7:
			m := make([]int64, y)
			err := r.r.Read(&m)
			if err != nil {
				return nil, nil, err
			}
			if l := m[y-1]; C[l] > 0 {
				D[i] = make(map[int]int)
				for j := 0; j < y-1; j++ {
					D[i][j] = int(m[j])
				}
				L[i] = int(l)
				C[l]--
				i++
			}
		}
	}

	return D, L, nil
}

// read reads .npy file from an *NpyReader, pulling n entries.
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
			if j == y-1 {
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
			if j == y-1 {
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
			if j == y-1 {
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
			if j == y-1 {
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
			if j == y-1 {
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
			if j == y-1 {
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
			if j == y-1 {
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
			if j == y-1 {
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

// ReadBalanced returns a balanced dataset and label slice totalling n instances. As argument, it
// takes the number of classes c.
func (r *NpyReader) ReadBalanced(n int, c int) ([]map[int]int, []int, error) {
	return readBalanced(r, n, c)
}

// Reset resets the file pointer so it points to the beginning of data.
func (r *NpyReader) Reset() error {
	r.f.Seek(0, 0)
	nr, err := npyio.NewReader(r.f)
	r.r = nr
	return err
}

// Close closes this reader's stream.
func (r *NpyReader) Close() {
	r.f.Close()
}
