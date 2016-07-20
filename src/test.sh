#!/bin/bash

for i in {0..10}
do
  go run main.go | tee out$i.put
done
