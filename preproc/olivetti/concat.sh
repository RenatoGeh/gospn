#!/bin/bash

k=0
out="/tmp/olivetti.png"

if ! [[ -z $1 ]]; then
  for i in `seq 1 20`; do
    cp "grid_${i}x1.png" $out
    for j in `seq 1 10`; do
      convert $out grid_${i}x${j}.png +append $out
    done
    cp $out row_${k}.png
    let k=k+1
    cp "grid_${i}x11.png" $out
    for j in `seq 11 20`; do
      convert $out grid_${i}x${j}.png +append $out
    done
    cp $out row_${k}.png
    let k=k+1
  done
fi

cp "row_$(( RANDOM % 40 )).png" $out
for i in `seq 1 4`; do
  s=$(( RANDOM % 40 ))
  convert $out row_${s}.png -append $out
done
cp $out olivetti_sample.png
