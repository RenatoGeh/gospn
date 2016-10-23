#!/bin/bash

suffix[0]="-LEFT.png"
suffix[1]="-BOTTOM.png"
suffix[2]="-TOP.png"
suffix[3]="-RIGHT.png"

orientation="-"

if [ $1 == "h" ]; then
  orientation="+"
else
  orientation="-"
fi

for f in *-TOP.png
do
  base=${f%-*}
  echo "Appending image: $base..."
  convert "$base${suffix[0]}" "$base${suffix[1]}" "$base${suffix[2]}" "$base${suffix[3]}" ${orientation}append "${base}.png"
done
