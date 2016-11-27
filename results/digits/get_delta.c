#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *args[]) {
  double min, max, t;
  int i, n;

  if (argc < 3) {
    printf("Usage: %s [p]...\n"
        "Prints the lowest and highest p values.\n", args[0]);
    return 1;
  }

  min = max = atof(args[1]);
  for (i=1, n=argc-1; i < n; ++i) {
    t = atof(args[i+1]);
    if (t > max)
      max = t;
    else if (t < min)
      min = t;
  }

  printf("%.2f %.2f\n", min, max);

  return 0;
}
