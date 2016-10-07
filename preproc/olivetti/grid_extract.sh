#!/bin/bash

stats=`file $1`
tokens=($stats)

width=${tokens[6]}
height=${tokens[8]}

nx_faces=$2
ny_faces=$3

grid_w=0
grid_h=0

echo "Image information:"

echo "Dimensions $width x $height"

let "grid_w=$width/$nx_faces"
let "grid_h=$height/$ny_faces"

let sx="$grid_w-1"
let sy="$grid_h-1"
echo "Constructing grid of $grid_w x $grid_h..."

dir_count=-1
k=0

for i in `seq 1 $nx_faces`;
do
  for j in `seq 1 $ny_faces`;
  do
    let t="$k % ($nx_faces/2)"
    if ! (( t )); then
      let dir_count="$dir_count + 1"
      echo "Creating new directory ${dir_count}..."
      mkdir "../../data/olivetti/$dir_count"
      mkdir "../../data/olivetti_simple/$dir_count"
    fi

    let dx="($j-1)*$grid_w"
    let dy="($i-1)*$grid_h"
    echo "Cropping a ${sx}x${sy} subimage at ${dx}x${dy}..."
    convert -extract ${sx}x${sy}+${dx}+${dy} $1 -compress none ../../data/olivetti/${dir_count}/grid_${i}x${j}.pgm
    if (( j == 5 || j == 15 )); then
      convert -extract ${sx}x${sy}+${dx}+${dy} $1 -compress none ../../data/olivetti_simple/${dir_count}/grid_${i}x${j}.pgm
    fi
    let k="$k+1"
  done
done
