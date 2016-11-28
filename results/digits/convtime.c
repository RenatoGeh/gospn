#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *args[]) {
  double secs;
  int h, m;

  if (argc < 2) {
    printf("Usage: %s secs\n"
        "Prints the amount of hours, minutes and seconds (up until the decimal case).\n", args[0]);
    return 1;
  }

  secs = atof(args[1]);

  h = (int) secs / 3600;
  m = (int) secs / 60;
  secs -= h*3600 + m*60;

  printf("%d %d %.2f\n", h, m, secs);

  return 0;
}
