GoSPN
=====

[![Build Status](https://travis-ci.org/RenatoGeh/gospn.svg?branch=stable)](https://travis-ci.org/RenatoGeh/gospn)
[![Go Report Card](https://goreportcard.com/badge/github.com/renatogeh/gospn)](https://goreportcard.com/report/github.com/renatogeh/gospn)
[![GoDoc](https://godoc.org/github.com/RenatoGeh/gospn?status.svg)](https://godoc.org/github.com/RenatoGeh/gospn)

![](./gospnpher.png "My crude attempt at drawing Renee French's Go Gopher.")

My crude (and slightly terrifying) rendition of Renee French's Go [Gopher](https://blog.golang.org/gopher) writing what's on his mind.

An implementation of Sum-Product Networks (SPNs) in Go
------------------------------------------------------

### Overview

Sum-Product Networks (SPNs) are deep probabilistic graphical models
(PGMs) that compactly represent tractable probability distributions.
Exact inference in SPNs is computed in time linear in the number of
edges, an attractive feature that distinguishes SPNs from other PGMs.
However, learning SPNs is a tough task. There have been many advances in
learning the structure and parameters of SPNs in the past few years. One
interesting feature is the fact that we can make use of SPNs' deep
architecture and perform deep learning on these models. Since the number
of hidden layers not only doesn't negatively impact the tractability of
inference of SPNs but also augments the representability of this model,
it is very much desirable to continue research on deep learning of SPNs.

This project aims to provide a simple framework for Sum-Product
Networks. Our objective is to provide inference tools and implement
various learning algorithms present in literature.

### Features

**Completed**

- Soft inference (marginal probabilities)
- Hard inference (MAP)
- Gens-Domingos learning schema
   * *Learning the Structure of Sum-Product Networks*, R. Gens & P.
     Domingos, ICML 2013
   * [pdf](http://spn.cs.washington.edu/papers/slspn.pdf)
- Generative gradient descent
- SPN derivatives
- Support for `.arff` dataset format (discrete variables for now)
- Poon-Domingos deep structure
- Dennis-Ventura clustering structural learning algorithm
    * *Learning the Architecture of Sum-Product Networks Using Clustering
     on Variables*, A. Dennis & D. Ventura, NIPS 25 (2012)
    * [pdf](http://papers.nips.cc/paper/4544-learning-the-architecture-of-sum-product-networks-using-clustering-on-variables.pdf)



**Under development**

- Testing performance on Poon-Domingos' deep structure using hard
  generative gradient descent.
- Testing performance on Dennis-Ventura's.


**To do (high priority)**

- Soft and Hard GD discriminative learning
  * *Discriminative Learning of Sum-Product Networks*, R. Gens & P.
    Domingos, NIPS 25 (2012)
  * [pdf](http://spn.cs.washington.edu/papers/dspn.pdf)

**To do (low priority)**

- Generative EM
- EM clustering
- Mauro-Vergari structural learning algorithm
    * *Simplifying, Regularizing and Strengthening Sum-Product Network
    Structure Learning*, A. Vergari & N. Mauro & F. Esposito, ECMLPKDD
    2015
    * [pdf](http://www.di.uniba.it/~vergari/papers/Simplifying,%20Regularizing%20and%20Strengthening%20Sum-Product%20Network%20Structure%20Learning.pdf)

- Language modelling SPN
  * *Language Modelling with Sum-Product Networks*, Cheng *et al*,
  INTERSPEECH 2014
  * [pdf](http://spn.cs.washington.edu/papers/is14.pdf)
 - Support for `.csv` dataset file format.

### Branches

- `dev` contains the development version of GoSPN.
- `stable` contains a stable version of GoSPN.
- `nlp` contains deprecated NLP model.

### Usage

#### As a Go library

Check the docs: https://godoc.org/github.com/RenatoGeh/gospn

You can find more elaborate examples at: https://github.com/RenatoGeh/benchmarks

To import the learning algorithms:

```
import "github.com/RenatoGeh/gospn/learn
```

To parse an ARFF format dataset and perform learning with the
Gens-Domingos structure learning algorithm:

First import the relevant packages (`learn/gens` for Gens' structural
learning algorithm, `io` for `ParseArff` and `spn` for inference
methods):

```
import (
  "github.com/RenatoGeh/gospn/learn/gens"
  "github.com/RenatoGeh/gospn/io"
  "github.com/RenatoGeh/gospn/spn"
)
```

Extract contents from an ARFF file (for now only discrete variables):

```
name, scope, values, labels := io.ParseArff("filename.arff")
```

Send the relevant information to the learning algorithm:

```
S := gens.Learn(scope, values, -1, 0.0001, 4.0, 4)
```

`S` is the resulting SPN. We can now compute the marginal probabilities
given a `spn.VarSet`:

```
evidence := make(spn.VarSet)
evidence[0] = 1 // Variable 0 = 1
// Summing out variable 1
evidence[2] = 0 // Variable 2 = 0
// Summing out all other variables.
p := S.Value(evidence)
// p is the marginal Pr(evidence), since S is already valid and normalized.
```

The method `S.Value` may repeat calculations if the SPN's graph is not a
tree. To use dynamic programming and avoid recomputations, use
`spn.Storer`:

```
T := spn.NewStorer()
t := T.NewTicket() // Creates a new DP table.
spn.StoreInference(S, evidence, t, T) // Stores inference values from each node to T(t).
p := T.Single(t, S) // Returns the first value inside node S: T(t, S).
```

Finding the approximate MPE works the same way. Let `evidence` be some
evidence, the MPE is given by:

```
args, mpe := S.ArgMax(evidence) // mpe is the probability and args is the argmax valuation.
```

Similarly to `S.Value`, `S.ArgMax` may recompute values if the graph is
not a tree. Use `StoreMAP` if the graph is a general DAG instead.

```
_, args := spn.StoreMAP(S, evidence, t, T)
mpe := T.Single(t, S)
```

#### As a standalone program

**GoSPN is no longer supported as a standalone program. Use it as a
library instead.**

### Dependencies

GoSPN is written in Go. Go is an open source language originally developed
at Google. It's a simple yet powerful and fast language built with
efficiency in mind. Installing Go is easy. Pre-compiled packages are
available for FreeBSD, Linux, Mac OS X and Windows for both 32 and
64-bit processors. For more information see <https://golang.org/doc/install>.

#### GoNum

We have deprecated GNU GSL in favor of GoNum (<https://github.com/gonum/>).
GoNum is written in Go, meaning when installing GoSPN, the Go package
manager should automatically install all dependencies (including GoNum).

In case this does not occur and something like this comes up on the
screen:

```
cannot find package "[...]/gonum/stat" in any of
```

Enter the following commands:

```
go get -u gonum.org/v1/gonum/stat
go get -u gonum.org/v1/gonum/mathext
```

We have deprecated functions that made GoSPN independent of GoNum or GNU
GSL, so we recommend installing GoNum.

The dependencies below can be considered "optional", in that the user
can avoid using them if they wish.

#### graph-tool (optional)

Graph-tool is a Python module for graph manipulation and drawing. Since
the SPNs we'll generate with most learning algorithms may have hundreads
of thousands of nodes and hundreds of layers, we need a fast and
efficient graph drawing tool for displaying our graphs. Since graph-tool
uses C++ metaprogramming extensively, its performance is comparable to a
C++ library.

Graph-tool uses the C++ Boost Library and can be compiled with OpenMP, a
library for parallel programming on multiple cores architecture that may
decrease graph compilation time significantly.

Compiling graph-tool can take up to 80 minutes and 3GB of RAM. If you do
not plan on compiling the graphs GoSPN outputs, it is highly recommended
that you do not install graph-tool.

Subdependencies and installation instructions are listed at
<https://graph-tool.skewed.de/download>.

#### Graphviz

GoSPN also supports graph drawing with Graphviz. See `io/output.go`.

### Compiling and Running GoSPN

To get the source code through Go's `go get` command, run the following
command:

```
$ go get -u github.com/RenatoGeh/gospn
```

This should install GoSPN to your $GOPATH directory. Note that this will
download the repo's default branch (namely `dev`). To get to the
`stable` branch, after downloading the repo as described above, run:

```
git checkout stable
```

~~Compiling the code is easy. First go to the GoSPN source dir.

```
$ cd $GOPATH/src/github.com/RenatoGeh/gospn/
```

To compile and run (the standalone version):

```
$ go run main.go <args>
```

Where `args` is a list of arguments. See Usage for more information.~~
(Standalone version is deprecated)

When using GoSPN as a library, simply run your code as usual.

### Updating GoSPN

To update GoSPN, run:

```
go get -u github.com/RenatoGeh/gospn
```

### Code and Docs Organization

In this section we describe the general layout that we intend to follow
for both code and documentation. For more information on SPNs, look for
the documentation present in this repository under directory `/doc`.

#### Code

Code documentation can be found at <https://godoc.org/github.com/RenatoGeh/gospn>.

#### Documentation

An analysis on our Gens-Domingos implementation can be found at
<https://github.com/RenatoGeh/gospn/blob/master/doc/analysis/analysis.pdf>.

### Datasets

We use the following datasets:

* Custom hand drawn numerical digits dataset (https://github.com/RenatoGeh/datasets/tree/master/digits)
* Olivetti Faces Dataset by AT&T Laboratories Cambridge
* Caltech101: L. Fei-Fei, R. Fergus and P. Perona. *Learning generative visual models
  from few training examples: an incremental Bayesian approach tested on
  101 object categories.* IEEE. CVPR 2004, Workshop on Generative-Model
  Based Vision. 2004

For a list of all available datasets, see:

* https://godoc.org/github.com/RenatoGeh/gospn/data
* https://github.com/RenatoGeh/datasets

### Results

In our [analysis'](https://github.com/RenatoGeh/gospn/blob/master/doc/analysis/analysis.pdf)
experiments section we show some results from the three datasets
enumerated above. We include some graphs and image completions here.
More images completions can be found at
[/results/olivetti_3bit/](https://github.com/RenatoGeh/gospn/tree/master/results/olivetti_3bit).
We also ran experiments on a modified Caltech-101 dataset due to memory
constraints. Read the analysis document (Experiments section) for more information.

#### Image classifications

![Digits dataset correct
classifications](https://raw.githubusercontent.com/RenatoGeh/gospn/dev/doc/analysis/imgs/digits_percs.png)

![Caltech dataset correct
classifications](https://raw.githubusercontent.com/RenatoGeh/gospn/dev/doc/analysis/imgs/caltech_percs.png)

#### Image completions with prior face knowledge

![Olivetti faces dataset C1 39
completions](https://raw.githubusercontent.com/RenatoGeh/gospn/dev/results/olivetti_3bit/r1/face_cmpl_39.png)

![Olivetti faces dataset C1 9
completions](https://raw.githubusercontent.com/RenatoGeh/gospn/dev/results/olivetti_3bit/r1/face_cmpl_9.png)

#### Image completions without prior face knowledge

![Olivetti faces dataset C2 39
completions](https://raw.githubusercontent.com/RenatoGeh/gospn/dev/results/olivetti_3bit/r2/face_cmpl_39.png)

![Olivetti faces dataset C2 9
completions](https://raw.githubusercontent.com/RenatoGeh/gospn/dev/results/olivetti_3bit/r2/face_cmpl_9.png)

### Literature

The following are articles that used GoSPN.

- *Credal Sum-Product Networks*, D. Mauá & F. Cozman & D. Conaty & C.
  Campos, PMLR 2017
    * [pdf](http://proceedings.mlr.press/v62/mau%C3%A117a/mau%C3%A117a.pdf)
- *Approximation Complexity of Maximum A Posteriori Inference in
  Sum-Product Networks*, D. Conaty & D. Mauá & C. Campos, UAI 2017
    * [pdf](https://arxiv.org/pdf/1703.06045.pdf)

### Acknowledgements

This project is part of my undergraduate research project supervised by
Prof. [Denis Deratani Mauá](https://www.ime.usp.br/~ddm/) at the
Institute of Mathematics and Statistics - University of São Paulo. We currently
have financial support from CNPq grant #800585/2016-0.

We would like to greatly thank Diarmaid Conaty and Cassio P. de Campos, both
from Queen's University Belfast, for finding and correcting several
bugs.
