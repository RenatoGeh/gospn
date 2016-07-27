#!/bin/bash

for i in {2..5}
do
  go run main.go $i | tee out$i.put
done
