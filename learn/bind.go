package learn

import "github.com/RenatoGeh/gospn/spn"

type LearnFunc func(sc map[int]*Variable, data spn.Dataset) spn.SPN
