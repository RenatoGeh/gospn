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

#### godebug

Although Go does work with GNU's GDB, GDB doesn't understand Go well.
An alternative to that is mailgun's godebug. If you do not intend on
using the debugging script `src/debug.sh` then there is no need for
godebug. Otherwise install it by running `go get`:

```
$ go get github.com/mailgun/godebug
```

More information at <https://github.com/mailgun/godebug>.

### Compiling and Running GoSPN

To get the source code through Go's `go get` command, run the following
command:

```
$ go get github.com/RenatoGeh/gospn
```

This should install GoSPN to your $GOPATH directory. Compiling the code
is easy. First go to the GoSPN source dir.

```
$ cd $GOPATH/github.com/RenatoGeh/gospn/src
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

### Code and Docs Organization

In this section we describe the general layout that we intend to follow
for both code and documentation. For more information on SPNs, look for
the documentation present in this repository under directory `/doc`.
