#!/bin/bash

for f in *
do
  filename=$(basename "$f")
  base="${filename%.*}"
  convert "${base}.pgm" "${base}.png"
done
