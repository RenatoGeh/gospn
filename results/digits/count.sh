#!/bin/bash

matches=($(grep -o 'Correct classifications: [0-9]\+\/[0-9]\+' *.txt))
n=${#matches[@]}
c=0
t=0

for i in `seq 0 $(( n - 1 ))`; do
  if [[ $(( i % 33 )) == 32 ]]; then
    def_ifs=$IFS
    IFS='/'
    tokens=( ${matches[i]} )
    IFS=$def_ifs
    let c=c+tokens[0]
    let t=t+tokens[1]
    echo "Corrects: $c vs Total: $t"
  fi
done

perc=`bc <<< "scale=3; ($c*100 / $t)"`
echo "Percentage: $perc"

