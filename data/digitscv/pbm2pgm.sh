#!/bin/bash

shopt -s nullglob
folds=(*)

for dir in "${folds[@]}"; do
  if [[ -d "$dir" ]]; then
    for f in ${dir}/*; do
      if [[ -f $f && $f == *.pbm ]]; then
        filename=$(basename "$f")
        base="${filename%.*}"
        convert -compress none "$dir/$base.pbm" "$dir/$base.pgm"
        sed -i -e 's/255/1/g' $dir/$base.pgm
        echo "Converting $dir/$base.pbm to $dir/$base.pgm"
      fi
    done
  fi
done
