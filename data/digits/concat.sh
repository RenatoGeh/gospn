#!/bin/bash

shopt -s nullglob

lits[0]="zero"
lits[1]="one"
lits[2]="two"
lits[3]="three"
lits[4]="four"
lits[5]="five"
lits[6]="six"
lits[7]="seven"
lits[8]="eight"
lits[9]="nine"

for i in `seq 1 10`
do
  blob=""
  for j in `seq 0 9`
  do
    blob="$blob$j/${lits[j]}$i.pgm "
  done
  convert $blob +append "digits_$i.png"
  echo "Appending images: ${blob}into digits_$i.png"
done

res=""
for i in `seq 1 5`
do
  res="${res}digits_$i.png "
done
convert $res -append "digits_left.png"
res=""
for i in `seq 6 10`
do
  res="${res}digits_$i.png "
done
convert $res -append "digits_right.png"
convert digits_left.png digits_right.png +append digits_sample.png

rm digits_left.png digits_right.png
for i in `seq 1 10`
do
  rm digits_$i.png
done
