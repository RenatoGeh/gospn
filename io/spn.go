package io

import (
	"github.com/RenatoGeh/gospn/spn"
	"io/ioutil"
	"os"
)

// SaveSPN serializes an SPN and writes it to a file. Suggested extension: ".spn".
func SaveSPN(filename string, S spn.SPN) error {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		return err
	}
	bytes := spn.Marshal(S)
	_, err = f.Write(bytes)
	return err
}

// LoadSPN reads a binary file that contains a serialized SPN.
func LoadSPN(filename string) (spn.SPN, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	S := spn.Unmarshal(data)
	return S, nil
}
