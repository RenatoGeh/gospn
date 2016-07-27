#!/bin/bash

dot2tex -tmath --figonly --preproc $1.dot > /tmp/tmpdot.dot
dot2tex --figonly /tmp/tmpdot.dot > $1.tex
