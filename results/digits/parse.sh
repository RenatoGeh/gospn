#!/bin/bash

# Find the percentages of correct hits on each run.
matches=($(grep -o '\(Percentage of correct hits: \)[0-9]\+\(.[0-9]\+\)*\%' *.txt))
n=${#matches[@]}
j=-1
c_avgs=0
c_points=0
c_diff=0
for i in `seq 0 $(( ( n - 1 ) / 5 ))`; do
  let j=j+5

  # Find the averages of each run.
  if [[ "$c_diff" -eq 10 ]]; then
    avgs[c_avgs]=${matches[j]}
    echo "Average for ${c_avgs}-th run: ${avgs[c_avgs]}"
    let c_avgs=c_avgs+1
    let c_diff=0
  # Find the non-averages (correct points on each iteration).
  else
    points[c_points]=${matches[j]}
    echo "Point ${c_points}: ${points[c_points]}"
    let c_points=c_points+1
    let c_diff=c_diff+1
  fi
done

# Write results to gnuplot data file.
:> avgs.dat
echo "# Gnuplot data file for 9 iterations of p." >> avgs.dat
echo "# p   avg   min   max" >> avgs.dat

# Compile delta C code.
gcc get_delta.c -o get_delta.out

st=0
end=9
# Enumerate y-axis (p=0.1 to p=0.9).
for i in `seq 0 8`; do
  min_max=( `./get_delta.out ${points[@]:$st:$end} `)
  let st=st+10
  let end=end+10
  avg=${avgs[i]}
  len_avg=${#avg}
  lend=$(( len_avg - 1 ))
  j=$(( i + 1 ))
  echo " 0.$j ${avg:0:lend} ${min_max[0]} ${min_max[1]}" >> avgs.dat
done

rm get_delta.out

# Create Gnuplot script.
:> $1.gpi
echo "set title \"Dataset $1 for p=[0.1, 0.9]: percentage of correct classifications\"" >> $1.gpi
echo "set xrange [0:1]" >> $1.gpi
echo "set xtics 0.1" >> $1.gpi
echo "set yrange [90:100]" >> $1.gpi
echo "set format y \"%.0f%%\"" >> $1.gpi
echo "set ytics 1" >> $1.gpi
echo "set key outside" >> $1.gpi
echo "set term png size 800,500" >> $1.gpi
echo "set output \"$1.png\"" >> $1.gpi
echo "plot 'avgs.dat' using 1:2 with linespoints lw 3 title \"Averages\", \
  '' using 1:2:3:4 with errorbars lw 3 title \"Min and max\", \
  '' using 1:2:(sprintf(\"%.2f%%\", \$2)) with labels center offset 0,-1 notitle" >> $1.gpi

# Run script.
gnuplot $1.gpi
feh $1.png
