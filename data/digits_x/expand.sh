#!/bin/bash

# Applies transformations to the original digits dataset to conform with the specifications of the
# digits_x dataset.

if [ "$1" == "--clean" ] || [ "$1" == "-c" ] || [ "$4" == "--clean" ] || [ "$4" == "-c" ]; then
  echo "Cleaning..."
  for i in `seq 0 9`; do
    rm $i -rf
  done
  if [ "$1" == "--clean" ] || [ "$1" == "-c" ]; then
    exit 0
  fi
fi

# Path to the original dataset (root directory of digits).
d_path=$1

# Blur parameters passed to ImageMagick.
b="$2"

# Depth (max pixel value in bits) passed to ImageMagick.
d="$3"

echo "Preparing..."
shopt -s nullglob
for i in `seq 0 9`; do
  mkdir -p "$i"
  cp -r $d_path/$i ./
  cd $i
  for f in *.pgm; do
    mv "$f" "_$f"
  done
  cd ..
done

g++ pgm_compress.cpp -o pgm_compress.out
echo "Applying transformations..."
for i in `seq 0 9`; do
  cd $i
  for f in *.pgm; do
    echo "$f, ${f:1:${#f}-1}"
    sed 's/\<1\>/15/g' "$f" | convert -blur "$b" -compress none - - | ../pgm_compress.out "$d" 1 > "./${f:1:${#f}-1}"
    rm "$f"
  done
  cd ..
done
rm pgm_compress.out
