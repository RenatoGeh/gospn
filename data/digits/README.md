## `digits` Dataset

The `digits` dataset is our own custom made dataset for hand-drawn
characters ranging from the digits 0 to 9.

Each subdirectory from `0` to `9` contains their respective digit
instances. Class instances are drawn as binary PGM files, with 0 being
white and 1 being black. This dataset contains only one handwriting
style and thus has little variance in drawing style.

Subfolder `compiled` contains a single compiled file version of all
class instances according to GoSPN's `data` file syntax.

Bash scripts contained in this directory include:

- `concat.sh`
    * Concatenates class instances into a single image file.
- `merge.sh`
    * Merges the contents of two directories into a single destination
      folder.
- `pbm2pgm.sh`
    * Converts PBM into PGM files.
