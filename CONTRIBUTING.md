# Contribution Guidelines

#### Introduction

The GoSPN library provides an easy and intuitive framework for building
[Sum-Product Networks](http://spn.cs.washington.edu/) (SPNs), computing
inference, and learning both structure and weights. GoSPN is written in
the Go programming language, and we highly encourage new contributions
to the library to improve the scope and quality of GoSPN. In this
document we provide a guide on how to contribute to GoSPN. Any questions
on GoSPN's development or usage can be directed to
[gospn-dev](https://groups.google.com/forum/#!forum/gospn-dev).

#### Table of Contents

[Project Scope](#project-scope)

[Packages](#packages)

[Contributing](#contributing)
  * [Working Together](#working-together)
  * [Reporting Bugs](#reporting-bugs)
  * [Suggesting Enhancements](#suggesting-enhancements)
  * [First Contribution](#first-contribution)
  * [Code Contribution](#code-contribution)
  * [Code Review](#code-review)
  * [How Can I Help?](#how-can-i-help)
  * [Style](#style)

## Project Scope

GoSPN's objective is to provide a general, easy, flexible and fast
library for SPNs. The library should provide the building blocks for new
learning algorithms as well as being a platform for applying SPNs to
real-life scenarios. Code should be implemented in pure Go. Calls to C
and assembly should be justified. Code should readable above all. Clever
obscure tricks should be avoided if possible. Whenever possible, a
reference to the original article of the algorithm implemented should be
added to the documentation.

## Packages

The GoSPN library is composed of several packages. Features should
respect the repository's package structure and add code to the most
relevant package. If a contribution doesn't fit any of the existing
packages, please start a discussion on the [mailing
list](https://groups.google.com/forum/#!forum/gospn-dev).

* app - applications in real-world use cases (e.g. image classification
  and completion)
* common - data structures and functions commonly used (e.g. queue,
  stack, colors)
* conc - concurrency data structures and algorithms
* data - dataset manipulation, dataset format importing and exporting
* io - input/output for dataset formats, images, HTTP and SPN
  importing/exporting
* learn - SPN derivation and learning algorithms
* score - scoring (e.g. classification score, confusion matrix)
* spn - inference, topology, searching and serialization of SPNs
* sys - garbage collection forcing, logging, global randomization and
  time
* utils - algorithms used in learning (e.g. clustering, G-test,
  chi-square, statistical functions, union-find, log-space computation)

## Contributing

#### Working Together

When contributing or participating, please follow the [Gopher
values](https://golang.org/conduct#values):

* Be friendly and welcoming
* Be patient
* Be thoughtful
* Be respectful
* Be charitable
* Avoid destructive behavior

#### Reporting Bugs

If you encounter a bug, please [open an
issue](https://github.com/RenatoGeh/gospn/issues/new) on the [GitHub
repository](https://github.com/RenatoGeh/gospn). The issue title should
be prepended by the package and subpackage name, e.g. `learn/gens: issue
name`. Please include in your report your environment, such as Go
version, OS and GoSPN version. If possible, write a small readable
sample code that reproduces the bug.

#### Suggesting Enhancements

If the scope of the enhancement is small, please open an issue. If it is
large, such as adding a new package to the repository or refactoring the
interface, please start a discussion on the [mailing
list](https://groups.google.com/forum/#!forum/gospn-dev).

#### First Contribution

Thank you for contributing to GoSPN! Every contributor should be added
to the
[CONTRIBUTORS](https://github.com/RenatoGeh/gospn/blob/dev/CONTRIBUTORS)
and [AUTHORS](https://github.com/RenatoGeh/gospn/blob/dev/AUTHORS) file.
As part of your pull request, please add yourself to them. The GoSPN
code follows the [BSD
license](https://github.com/RenatoGeh/gospn/blob/dev/LICENSE). When
adding code, be aware of the licensing that comes with it. All added
code must not conflict with the BSD license.

#### Code Contribution

Every pull request should be self-contained and only address a single
issue. If you wish to contribute with several features or fixes, please
open a separate pull request for each. This allows for code to be more
easily reviewed, increasing the chances your contribution is merged to
the master branch. Commits should:

- Have commit titles preferably not exceed 50 characters
- Have commit titles prepended by the package and subpackage
- Have commit body message lines never exceed 72 characters
- Always be Signed-off-by the authors
- Be clear on what it does
- Only do one thing

An example of a good commit message:

```
learn: gens: add concurrency to gens.Learn

This patch adds concurrency to the Gens-Domingos learning algorithm.
This is done through a conc.Queue in every cluster and independency
step. It also limits the number of concurrent processes to a given
quantity.

Signed-off-by: Author Name <author@email.com>
```

Every contribution should also be formatted with
[gofmt](https://golang.org/cmd/gofmt/) and preferably pass
[golint](https://github.com/golang/lint). Functions should always be
documented following the standard Go documentation format.

#### Code Review

When reviewing code, please be nice and welcoming to new contributors.
We've all been there, and a friendly review might mean more
contributions in the future, meaning everyone wins when we're awesome to
each other.

* `LGTM` — looks good to me
* `SGTM` — sounds good to me
* `s/foo/bar/` — please replace `foo` with `bar`
* `s/foo/bar/g` — please replace every occurrence of `foo` with `bar`
* `PTAL` - please take a look

At least one reviewer must aprove (by saying LGTM or manually aproving
the PR) before a merge. Any contributor may also ask for a PTAL from
another contributor for an additional review.

#### How Can I Help?

If you are looking to help GoSPN, you may search for [open
issues](https://github.com/RenatoGeh/gospn/issues). If you are new, good
first contributions are labelled `good first issue`. Improving
documentation also counts as contributions, and you are more then
welcome to help us with that! Tests, performance patches, fixes are all
good contributions.

#### Style

We follow the [Go style](https://github.com/golang/go/wiki/CodeReviewComments).

## Acknowledgements

This contribution guide was inspired by
[GoNum's](https://github.com/gonum/gonum/blob/master/CONTRIBUTING.md).
