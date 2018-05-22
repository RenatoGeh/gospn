#include <cstdio>
#include <cstdlib>

/* Compresses PGM files.
 *
 * Example:
 *
 * Given an m-bit PGM file F:
 *
 * Let q = 2^m - 1.
 *   | P2
 *   | 46 56
 *   | q
 *   | ...
 *
 * Compresses F into an n-bit PGM file H, with n < m:
 *
 * Let p = 2^n - 1.
 *   | P2
 *   | 46 56
 *   | p
 *   | ...
 */
int main(int argc, char *args[]) {
  int bit = 4, max = 1;
  int r = 0;

  if (argc < 2) {
    printf("Usage: %s bit r\n  bit - number of bits for the image's max value\n  r - r=0 if bit "
        "is the number of bits to be used, else bit is the max value itself.\n",
        args[0]);
    return 1;
  }

  if (argc > 1)
    bit = atoi(args[1]);
  if (argc > 3)
    r = atoi(args[3]);

  if (!r)
    max <<= bit;
  else
    max = bit;

  int w, h, omax;
  scanf("P2 %d %d %d", &w, &h, &omax);

#define MAX_FILENAME_SIZE 50
  char filename[MAX_FILENAME_SIZE];
#undef MAX_FILENAME_SIZE

  if (argc > 2)
    sprintf(filename, "%s_%d-bit.pgm", args[2], bit);
  else
    sprintf(filename, "img_%d-bit.pgm", bit);

  printf("P2\n%d %d\n%d\n", w, h, --max);

  int n = w*h;
  double df = (double) (max+1) / (double) (omax+1);
  for (int i = 0; i < n; ++i) {
    int opx;
    scanf("%d", &opx);
    double npx = (double) opx * df;
    printf("%d", (int) npx);
    if (i % w == w-1)
      putchar('\n');
    else
      putchar(' ');
  }

  return 0;
}
