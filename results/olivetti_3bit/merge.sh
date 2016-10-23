#!/bin/bash

nfiles=`ls -1 *.png | wc -l`

i=0
k=0
while [ $i -lt $nfiles ]; do
  j=$i
  list=""
  n=0
  let n=i+5
  while [ $j -lt $n ]; do
    list="$list cmpl_${j}.png"
    let j=j+1
  done
  convert $list -append lface_cmpl_${k}.png
  let n=i+10
  list=""
  while [ $j -lt $n ]; do
    list="$list cmpl_${j}.png"
    let j=j+1
  done
  convert $list -append rface_cmpl_${k}.png
  convert lface_cmpl_${k}.png -gravity east -background white -splice 1x0 blface_cmpl_${k}.png
  convert blface_cmpl_${k}.png rface_cmpl_${k}.png +append face_cmpl_${k}.png
  convert face_cmpl_${k}.png -bordercolor white -border 1x1 face_cmpl_${k}.png
  rm lface_cmpl_${k}.png rface_cmpl_${k}.png blface_cmpl_${k}.png
  let i=i+10
  let k=k+1
done
