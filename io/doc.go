/*
	Package io is an Input/Output package.
	GoSPN reads and writes from .data files. To run GoSPN we must first convert a dataset into a data
	file. For now, GoSPN supports converting PGM and PBM image files into data files.

	Converting PGM files is done by io/pgm.go, whilst PBM files are handled by io/pbm.go. Function
	names are (supposed to be) intuitive: the input format (e.g. PGM) followed by a suffix to
	indicate whether it is a folder or not (e.g. F) to the output format data (e.g. PGMFToData).
	The Buffered variant is for big datasets. Instead of saving every file stream in memory, we
	concurrently run each stream according to the number of CPUs in the user's machine.

	We differentiate Data from Evidence. Data is supposed to contain the classification labels, that
	is, data is the training set. Evidence removes the instance's labels and acts as test set.

	For output we follow the same format as input. VarSetToPGM, for instance, takes a variable
	instantiation set and converts it into a PGM image. This is useful for image completion.

	Other output functions include DrawGraphTools and DrawGraph. DrawGraphTools draws the given SPN
	into a graph-tool python script. This script can be run just like any pythons script. After
	doing so, a new image of the SPN will be generated. Note that this requires the graph-tool
	library (https://graph-tool.skewed.de/). DrawGraph uses Graphviz to draw the graph. You can then
	run the resulting dot script with sfdp, neato or any other layout program. This requires the
	graphviz library (http://www.graphviz.org/).

	WriteToFile writes the SPN to a file. TODO: ReadFromFile should read the SPN from a .mdl file.
*/
package io
