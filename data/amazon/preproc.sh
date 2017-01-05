#!/bin/bash

:> amazon_unquoted.txt
while read l; do
  n=$(( ${#l} - 1 ))
  echo "${l:1:$n}" >> amazon_unquoted.txt
done
