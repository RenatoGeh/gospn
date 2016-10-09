GoSPN
=====

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

For the independence test (`src/utils/indep/indtest.go`). A builtin
Chi-Square function is already present in `src/utils/indep/indtest.go`
under the name of `Chisquare`. However `Chisquare` has worse numerical
error when compared to its GSL equivalent `ChiSquare` (see
`src/utils/indep/chisq.go`).

For information on how to compile and install GNU GSL, see
<https://www.gnu.org/software/gsl/>.

If you do not wish to install GNU GSL, simply rename `ChiSquare` in file
`src/utils/indep/indtest.go` to `Chisquare`.

GoSPN uses Go's `cgo` to run C code inside Go. File
`src/utils/indep/chisq.go` contains the wrapper function `ChiSquare`
that calls `gsl_cdf_chisq_P` from `gsl/gsl_cdf.h`.

#### godebug

Although Go does work with GNU's GDB, GDB doesn't understand Go well.
An alternative to that is mailgun's godebug. If you do not intend on
using the debugging script `src/debug.sh` then there is no need for
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
$GOPATH/src/github.com/RenatoGeh/gospn/results/dataset_name/
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
$ cd $GOPATH/src/github.com/RenatoGeh/gospn/src
```

To compile and run:

```
$ go run main.go
```

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

Source code for GoSPN is organized in Go packages. Each Go package is a
component of GoSPN's main package:

* `common`: contains the usual "common" data structures, such as pairs,
  queues and stacks.
* `io`: contains Input/Output code. Namely evidence/data parsing and
  graph drawing.
* `learn`: contains learning algorithms.
* `spn`: code that encapsulates the structure of SPNs (nodes and
  MAP, evidence inference).
* `utils`: algorithms that deserve their own package (e.g. clustering
  and independency tests)

For more information on each source file, generate a `godoc` doc page.

```
# Creates a server with all the documentation at localhost:6060.
$ godoc -http=:6060
# Replace chromium with your favorite browser.
$ chromium 127.0.0.1:6060/pkg/github.com/RenatoGeh/gospn
```

If `godoc` is not installed, install it. For Arch Linux:

```
$ pacman -S go-tools
```

Linux distributions that contain a package manager should have similar
package names (e.g. golang-tools for Debian/Ubuntu). If that doesn't
work:

```
$ sudo -E go get golang.org/x/tools/cmd/godoc
```

For more information on `godoc` <http://godoc.org/golang.org/x/tools/cmd/godoc>.

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
