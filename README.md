GoSPN
=====

[![Build Status](https://travis-ci.org/RenatoGeh/gospn.svg?branch=master)](https://travis-ci.org/RenatoGeh/gospn)
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

This project aims to provide a simple implementation of structural
learning of SPNs. We seek to follow the paper *Learning the Structure of
Sum-Product Networks* by Robert Gens and Pedro Domingos (ICML 2013) and
implement our own version of structural learning based on the schema
provided by the article.

Our objective is not only educational - in the sense that we wish to
learn more about the peculiarities of SPNs - but also documentational,
as we also intend on documenting and recording what we have learned in a
simpler, clearer way then how it is currently written in literature.

### Usage

To run GoSPN, we must complete a few steps:

1. Prepare the dataset:
  - Let `ds` be your dataset's name.
  - Create a new directory `/data/ds`, where root is the root of the
    GoSPN package.
  - Each subdirectory inside `/data/ds` represents a different class.
  - For example: if we have three classes, `dog`, `cat` and `mouse`,
    then we might have three subdirectories inside `/data/ds` named
    `dog`, `cat` and `mouse`.
  - Copy your class instances into `/data/ds/classname`.
2. Compile the dataset into a `.data` file:
  - If the dataset is an image, take note of the dimensions and max
    value pixels take.
  - Let `w` and `h` be the width and height of the images, and `m` be
    the max value.
  - Compile the data with `go run main.go -mode=data -width=w -height=h
    -max=m -dataset=ds`
  - This will generate a `.data` file inside `/data/ds/all/`. By default
    it is named `all.data`.
3. Run a job by running GoSPN with the following syntax.


```
Usage:
  go run main.go [-p] [-rseed] [-clusters] [-iterations] [-concurrents]
Arguments:
  p           - is the partition in the interval (0, 1) to be used for
                cross-vdalidation. If ommitted, p defaults to 0.7.
  rseed       - the seed to be used when choosing which instances to be used
                as train and which to be used as test set. If ommitted, rseed
                defaults to -1, which chooses a random seed according to the
                current time.
  clusters    - how many k-clusters to be used during training on instance
                splits. If clusters = -1, then use DBSCAN. Else if
                clusters = -2, then use OPTICS. Else, if clusters > 0,
                then use k-means clustering. By default, clusters is set
                to -1.
  iterations  - how many iterations to be run when running a
                classification job. This allows for better, more general
                and randomized results, as some test/train partitions may
                become degenerated.
  concurrents - GoSPN makes use of Go's native concurrency and is able
                to run on multiple cores in parallel. Argument concurrents
                defines the number of concurrent jobs GoSPN should run
                at most. If concurrents <= 0, then concurrents = nCPU,
                where nCPU is the number of CPUs the running machine has
                available. By default, concurrents = -1.
  dataset     - name of the dataset to be parsed or compiled. Setting
                -mode=data will compile data. Ommitting -mode or setting
                -mode to something different than data will either run
                completion or classification.
  width       - width of the images to be classified or completed.
  height      - height of the images to be classified or completed.
  max         - maximum pixel value the images can take.
  mode        - whether to convert a directory structure into a data
                file (data), run an image completion job (cmpl) or a
                classification job (class).
```

Running `go run main.go -help` shows the help page.

#### For step 3 to run a classification job:

1. Choose a partition value `p` such that `0 < p < 1`. For instance,
   `p=0.8`.
2. Choose an `rseed` value. For instance, `rseed=-1`.
3. Choose a `clusters` value (k-means with `clusters` clusters, DBSCAN
   or OPTICS). For instance, `clusters=-1`.
4. Choose the number of iterations `iterations`. For instance,
   `iterations=5`.
5. Run `go run main.go -p=0.8 -rseed=-1 -clusters=-1 -iterations=5`.

#### For step 3 to run an image completion job:

1. Run `go run main.go -p=0`.

It is recommended that you output `stdout` to some file, since GoSPN can
be quite verbose in certain cases. Append a `| tee out.put` to generate
a file `out.put` with all the information GoSPN prints to the screen
whilst still printing `stdout`. For example:

```
  go run main.go -p=0 | tee out.put
```

#### Independence tests:

GoSPN uses either G-Test or Chi-Square. To use G-Test, use the function

```
GTest(p, q int, data [][]int) float64
```

To use Chi-Square:

```
ChiSquareTest(p, q, int, data [][]int) float64
```

In `/utils/indep/indgraph.go`, search for either functions and
replace with the desired function.

To change the significance level on the independence tests, go to
`/utils/indep/chisq.go` and `/utils/indep/gtest.go` and change
the `sigval` variable.

#### DBSCAN and OPTICS parameters:

To change DBSCAN's and OPTICS's parameters, go to `/learn/gens.go`
and change the function`s

```
DBSCAN(data [][]int, eps float64, mp int)
OPTICS(data [][]int, eps float64, mp int)
```

parameters to suit your dataset.

### Dependencies

GoSPN is built in Go. Go is an open source language originally developed
at Google. It's a simple yet powerful and fast language built with
efficiency in mind. Installing Go is easy. Pre-compiled packages are
available for FreeBSD, Linux, Mac OS X and Windows for both 32 and
64-bit processors. For more information see <https://golang.org/doc/install>.

#### Installing Go on Arch Linux
```
# Choose one of the following:
$ pacman -S go     # for the gc compiler (official)
$ pacman -S gcc-go # for the gccgo compiler (frontend)
```

#### Installing Go on Ubuntu
```
$ sudo add-apt-repository ppa:ubuntu-lxc/lxd-stable
$ sudo apt-get update
$ sudo apt-get install golang
```

#### Installing Go on Mac OS X

Follow instructions at <https://golang.org/doc/install#darwinPackageInstructions>.

#### Installing Go on Windows

Follow instructions at <https://golang.org/doc/install#windows>.

#### $GOPATH

Once Go is installed, be sure to check if your $GOPATH is set correctly.
From now on all Go packages should be installed to $GOPATH.

If you're using Linux, your `.bashrc` or `.zshrc` should have the
following lines:

```
# Go path. Replace $YOURDIR with a directory of your choice.
export GOPATH=$YOURDIR/go
# Optionally add Go's path to your $PATH environment.
export PATH="$PATH:$GOPATH/bin"
```

#### GNU GSL Scientific Library

GoSPN uses GNU GSL to compute the cumulative probability function

```
Pr(X^2 <= chi), X^2(df)
```

For the independence test (`utils/indep/indtest.go`). A builtin
Chi-Square function is already present in `utils/indep/indtest.go`
under the name of `Chisquare`. However `Chisquare` has worse numerical
error when compared to its GSL equivalent `ChiSquare` (see
`utils/indep/chisq.go`).

For information on how to compile and install GNU GSL, see
<https://www.gnu.org/software/gsl/>.

If you do not wish to install GNU GSL, simply rename `ChiSquare` in file
`utils/indep/indtest.go` to `Chisquare`.

GoSPN uses Go's `cgo` to run C code inside Go. File
`utils/indep/chisq.go` contains the wrapper function `ChiSquare`
that calls `gsl_cdf_chisq_P` from `gsl/gsl_cdf.h`.

#### godebug

Although Go does work with GNU's GDB, GDB doesn't understand Go well.
An alternative to that is mailgun's godebug. If you do not intend on
using the debugging script `debug.sh` then there is no need for
godebug. Otherwise install it by running `go get`:

```
$ go get github.com/mailgun/godebug
```

More information at <https://github.com/mailgun/godebug>.

#### graph-tool

Graph-tool is a Python module for graph manipulation and drawing. Since
the SPNs we'll generate with this algorithm may have thousands of nodes
and hundreds of layers, we need a fast and efficient graph drawing tool
for displaying our graphs. Since graph-tool uses C++ metaprogramming
extensively, its performance is comparable to a C++ library.

Graph-tool uses the C++ Boost Library and can be compiled with OpenMP, a
library for parallel programming on multiple cores architecture that may
decrease graph compilation time significantly.

Compiling graph-tool can take up to 80 minutes and 3GB of RAM. If you do
not plan on compiling the graphs GoSPN outputs, it is highly recommended
that you do not install graph-tool.

##### graph-tool dependencies

* C++14 compiler (GCC version 5 or above; or clang 3.4 or above),
* C++ Boost libraries v1.54+,
* Python version 2.7, 3 or above
* expat
* SciPy
* NumPy
* CGAL C++ Geometry library v1.7+

Optional dependencies are listed at <https://graph-tool.skewed.de/download>.

##### Installing graph-tool

After installing all dependencies, compile graph-tool by downloading the
source (<https://graph-tool.skewed.de/download>) and compiling the usual
way:

```
./configure
make
```

Then install the Python module:

```
make install
```

If you use Debian, Ubuntu, Arch Linux, Gentoo or Mac OS X, there are
pre-compiled packages available:

###### Debian and Ubuntu

Read <https://graph-tool.skewed.de/download#debian>.

###### Arch Linux

Graph-tool is available at the AUR. Replace `pacaur` with your favorite
AUR helper.

```
pacaur -S python-graph-tool
```

###### Gentoo

```
emerge graph-tool
```

###### Mac OS X

Read <https://graph-tool.skewed.de/download#note-macos> first.

```
# Macports
port install py-graph-tool
# Homebrew
brew tap homebrew/science
brew install graph-tool
```

##### Compiling graphs

GoSPN outputs all graphs to:

```
$GOPATH/github.com/RenatoGeh/gospn/results/dataset_name/
```

Simply run python on the python source code and it will output a PNG
image of the graph.

```
python graph_name.py
```

If graph-tool was compiled with OpenMP it will make use of all available
CPU cores to compile the graph.

### Compiling and Running GoSPN

To get the source code through Go's `go get` command, run the following
command:

```
$ go get github.com/RenatoGeh/gospn
```

This should install GoSPN to your $GOPATH directory. Compiling the code
is easy. First go to the GoSPN source dir.

```
$ cd $GOPATH/github.com/RenatoGeh/gospn/
```

To compile and run:

```
$ go run main.go <args>
```

Where `args` is a list of arguments. See Usage for more information.

To run GoSPN in debug mode:

```
$ ./debug.sh
```

It is recommended that you redirect the standard output to some output
file:

```
$ go run main.go | tee out.put
```

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

The available documentation present at `/doc` does not only concern the
code nor the algorithms implemented in GoSPN. It also provides an
introduction to SPNs in the form of a tutorial, explaining how knowledge
is represented in SPNs and how to perform exact inference. It obviously
also contains a detailed description on the learning algorithms
implemented in GoSPN.

There are two submodules under `/doc`:

* `/doc/tutorial`: is a tutorial on SPNs. It covers from how to
  represent knowledge to inference and learning in SPNs. It is a
  detailed document on SPNs.
* `/doc/code`: is a detailed documentation on the code. It contains only
  the implementational aspect of GoSPN.

### Datasets

We use the following datasets:

* Olivetti Faces Dataset by AT&T Laboratories Cambridge
* Caltech101: L. Fei-Fei, R. Fergus and P. Perona. *Learning generative visual models
  from few training examples: an incremental Bayesian approach tested on
  101 object categories.* IEEE. CVPR 2004, Workshop on Generative-Model
  Based Vision. 2004
