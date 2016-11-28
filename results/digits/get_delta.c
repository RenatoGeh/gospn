#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main(int argc, char *args[]) {
  double min, max, t;
  int i, n, f;

  if (argc < 3) {
    printf("Usage: %s [p]... [-f]\n"
        "Prints the lowest and highest p values. Floors if -f.\n", args[0]);
    return 1;
  }

  f = !strcmp(args[argc-1], "-f");
  min = max = atof(args[1]);
  for (i=1, n=f?argc-2:argc-1; i < n; ++i) {
    t = atof(args[i+1]);
    if (t > max)
      max = t;
    else if (t < min)
      min = t;
  }

  if (f)
    printf("%d %d\n", (int) min, (int) max);
  else
    printf("%.2f %.2f\n", min, max);

  return 0;
}
