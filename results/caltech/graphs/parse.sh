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

# Create Gnuplot script.
out=${1}_percs.gpi
:> $out
echo "set title \"Dataset $1 for p=[0.1, 0.9]: percentage of correct classifications\"" >> $out
echo "set xrange [0:1]" >> $out
echo "set format x \"%.1f\"" >> $out
echo "set xtics 0.1" >> $out
echo "set xlabel \"(p)\nPartition for cross-validation\"" >> $out
echo "set yrange [:100]" >> $out
echo "set format y \"%.0f%%\"" >> $out
echo "set ytics 5" >> $out
echo "set ylabel \"Correct classifications\n(\%)\"" >> $out
echo "set grid" >> $out
echo "set key outside" >> $out
echo "set term png size 800,500" >> $out
echo "set output \"${1}_percs.png\"" >> $out
echo "plot 'avgs.dat' using 1:2 with linespoints lw 3 title \"Averages\", \
  '' using 1:2:3:4 with errorbars lw 3 title \"Min and max\", \
  '' using 1:2:(sprintf(\"%.2f%%\", \$2)) with labels center offset 0,-1 notitle" >> $out

# Run script.
gnuplot $out
#feh ${1}_percs.png

# Find user time elapsed in seconds.
matches=($(grep -o 'User time (seconds): [0-9]\+\(.[0-9]\+\)*' *.txt))
n=${#matches[@]}
j=0

for i in `seq 0 $(( n - 1 ))`; do
  if [[ $(( i % 4 )) == 3 ]]; then
    secs[j]=${matches[i]}
    echo "Run $j lasted: ${matches[i]} seconds."
    let j=j+1
  fi
done

gcc convtime.c -o convtime.out

:> time.dat
echo "# Gnuplot user time data file for 9 iterations of p." >> time.dat
echo "# p  h  min  sec" >> time.dat

n=${#secs[@]}
min_max=( `./get_delta.out ${secs[@]} -f` )
for i in `seq 1 $n`; do
  j=$(( i - 1 ))
  cvt=`./convtime.out ${secs[j]}}`
  echo "0.$i $cvt" >> time.dat
done

out=${1}_time.gpi
:> $out
echo "set title \"Dataset $1 for p=[0.1, 0.9]: user time\"" >> $out
echo "set xrange [0:1]" >> $out
echo "set format x \"%.1f\"" >> $out
echo "set xtics 0.1" >> $out
echo "set xlabel \"(p)\nPartition for cross-validation\"" >> $out
echo "set timefmt \"%H %M %S\"" >> $out
echo "set ydata time" >> $out
echo "set format y \"%tH:%tM:%.0tS\"" >> $out
echo "set ylabel \"Elapsed time\n(hours:minutes:seconds)\"" >> $out
echo "set yrange [$(( min_max[0] - 90 )):$(( min_max[1] + 90 ))]" >> $out
echo "set grid" >> $out
echo "set key outside" >> $out
echo "set term png size 800,400" >> $out
echo "set output \"${1}_time.png\"" >> $out
echo "plot 'time.dat' using 1:2:3:4 with linespoints lw 3 title \"Total running time\", \
  '' using 1:2:(sprintf(\"(%d:%d:%02.0f)\", \$2, \$3, \$4)):3 with labels center offset -1.5,0.5 font \
  ',8' notitle" >> $out

# Run script.
gnuplot $out
#feh ${1}_time.png

rm convtime.out

# Find memory usage.
matches=($(grep -o 'Maximum resident set size (kbytes): [0-9]\+' *.txt))
n=${#matches[@]}
j=0
for i in `seq 0 $(( n - 1 ))`; do
  if [[ $(( i % 6 )) == 5 ]]; then
    mem[j]=${matches[i]}
    let j=j+1
  fi
done

n=${#mem[@]}
min_max=( `./get_delta.out ${mem[@]} -f` )

:> mem.dat
echo "# Gnuplot memory data file for 9 iterations of p." >> mem.dat
echo "# p  kbytes mbytes" >> mem.dat

for i in `seq 1 9`; do
  j=$(( i - 1 ))
  kb=${mem[j]}
  mb=`bc <<< "scale=2; $kb / 1000.0"`
  echo "KB: $kb = MB: $mb"
  echo " 0.$i $kb $mb" >> mem.dat
done

# Create Gnuplot script.
out=${1}_mem.gpi
:> $out

echo "set title \"Dataset $1 for p=[0.1, 0.9]: memory usage\"" >> $out
echo "set format x2 \"%.1f\"" >> $out
echo "set auto x2" >> $out
echo "set x2tics 0.1" >> $out
echo "set x2label \"Partition for cross-validation\n(p)\"" >> $out
echo "set auto y" >> $out
#echo "set format y \"%f\"" >> $out
#echo "set ytics 1" >> $out
echo "set ylabel \"Total memory used\n(mB)\"" >> $out
echo "set style data histogram" >> $out
echo "set style fill solid 0.1" >> $out
echo "set boxwidth 0.07" >> $out
echo "set xlabel \"(mB)\nMemory values in megabytes\"" >> $out
echo "set grid" >> $out
echo "set key outside" >> $out
echo "set term png size 1000,500" >> $out
echo "set output \"${1}_mem.png\"" >> $out
echo "plot 'mem.dat' using 1:3:x2ticlabel(1):xticlabel(3) with boxes title 'RAM used'" >> $out

gnuplot $out

# Preview graphs.
feh *.png

rm get_delta.out
