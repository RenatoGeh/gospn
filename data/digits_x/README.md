## `digits`-expanded Dataset

The `digits` expanded dataset (or `digits_x`) is our own custom made dataset for hand-drawn
characters ranging from the digits 0 to 9 based on the `digits` dataset.

Each subdirectory from `0` to `9` contains their respective digit
instances. Class instances are drawn as binary PGM files, with 0 being
white and 255 being black. This dataset contains only one handwriting
style and thus has little variance in drawing style.

The `digits_x` dataset is merely the `digits` dataset with a
transformation applied to the images. For each image in the original
`digits` dataset, we increased the max pixel value from 2 to 256 and
applied a gaussian blur to increase pixel intensity variety.

Subfolder `compiled` contains a single compiled file version of all
class instances according to GoSPN's `data` file syntax.
