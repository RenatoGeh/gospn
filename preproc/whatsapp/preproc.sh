#!/bin/bash

cat "$1" | perl -pe 's/[0-9]+\/[0-9]+\/[0-9]+\, [0-9]+\:[0-9]+ - .+: //g' > result.txt
cat result.txt | perl -pe 's/[0-9]+\/[0-9]+\/[0-9]+\, [0-9]+\:[0-9]+ - .+//g' > t.out
cat t.out | perl -pe 's/\<Media omitted\>//g' > result.txt
cat result.txt | perl -pe 's/https:\/\/([^\s]+)//g' > t.out
cat t.out | perl -pe 's/((http(s)*|ftp):\/\/|www)([^\s]+)//g' > result.txt
cat result.txt | perl -pe 's/^\s*\n//gm' > out.txt
mv out.txt result.txt
