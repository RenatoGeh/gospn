package language

import (
	"github.com/RenatoGeh/gospn/spn"
)

// Language is a language modelling SPN structure based on the article
// 	Language Modelling with Sum-Product Networks
// 	Wei-Chen Cheng, Stanley Kok, Hoai Vu Pham, Hai Leong Chieu, Kian Ming A. Chai
// 	INTERSPEECH 2014
// We shall refer to this article via the codename LMSPN.
// vfile is the vocabulary filename.
func Language(vfile string) {

}

// Structure returns the SPN structure as described in LMSPN.
func Structure(K, D, N int) spn.SPN {
	return nil
}
