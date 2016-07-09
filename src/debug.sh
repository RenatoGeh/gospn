#!/bin/bash

# Uses godebug as the interactive debugger.
# Cheatsheet:
#
# To add a breakpoint, add:
#
#   _ = "breakpoint"
#
# To the break line.
#
#   h(elp)     | show help
#   n(ext)     | next line
#   s(tep)     | next step
#   c(ontinue) | continue until next breakpoint
#   l(ist)     | print code around the current line
#   p(rint)    | print var or expression
#   q(uit)     | exit

pkgs=""

# Get directories.
for f in *; do
  if [[ -d $f ]]; then
    pkgs="github.com/RenatoGeh/gospn/src/$f,$pkgs"
  fi
done

# Remove last character ','.
pkgs=${pkgs:0:${#t}-1}

godebug run -instrument=$pkgs main.go

