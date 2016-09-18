#!/bin/bash

nums=( zero one two three four five six seven eight nine )
k=0
len_nums=${#nums[@]}

while [ $k -lt $len_nums ]; do
  p1="$1/$k"
  p2="$2/$k"
  echo "Getting ready to merge dirs $p1 and $p2..."

  n="$(ls -1 $p1 | wc -l)"
  m="$(ls -1 $p2 | wc -l)"

  echo "Dir $p1 has $n files."
  echo "Dir $p2 has $m files."

  i=0
  j=0

  mkdir "$k"

  echo "Copying files from $p1..."
  while [ $i -lt $n ]; do
    origin="$1/$k/${nums[$k]}$(($i+1)).pbm"
    dest="$k/${nums[$k]}$(($j+1)).pbm"
    echo "Copying $origin to $dest..."
    cp $origin $dest
    let j=j+1
    let i=i+1
  done
  let i=0
  while [ $i -lt $m ]; do
    origin="$2/$k/${nums[$k]}$(($i+1)).pbm"
    dest="$k/${nums[$k]}$(($j+1)).pbm"
    echo "Copying $origin to $dest..."
    cp $origin $dest
    let j=j+1
    let i=i+1
  done

  let k=k+1
done
