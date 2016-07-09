#!/bin/bash

FILE=$1

dot -Tpng $FILE > dot.png
neato -Tpng $FILE > neato.png
fdp -Tpng $FILE > fdp.png
sfdp -Tpng $FILE > sfdp.png
circo -Tpng $FILE > circo.png
twopi -Tpng $FILE > twopi.png
patchwork -Tpng $FILE > patchwork.png

