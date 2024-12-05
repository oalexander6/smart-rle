# Smart Run-Length Encoding

Package rle implements the necessary functionality to encode and decode byte arrays
using a smart run-length encoding scheme. In this scheme, consecutive identical
bytes will be replaced by <delimiter><length><delimiter><original byte> if and only
if this results in a shorter output.

In order for this encoding to work, there must be a byte value that is never present
in the input to use as a delimiter. If this byte is ever found in the input to encode,
it will return an error.

The goal of this implementation is to optimize the worst-case compression ratio of
this encoding. The result of encoding an array of bytes with this implementation
will be, at worst, the same length as the input. In a traditional RLE scheme, each
byte is prefixed by a count, which could result in an output that is up to double
the size of the input for a highly random input.

An additional optimization is used by encoding the run lengths using base-62, where
the digits 0-9, a-z, and A-Z are all used. This increases the number of lengths that
can be described with a given number of characters, allowing for the swapping of
consecutive bytes with a run length to occur more frequently than with base-10. This
also further improves the compression ratio.

The delimiter is always encoded as the first byte of the output. If the input is empty,
then the delimiter is not specified and an empty output is returned. For decoding,
the delimiter is always pulled from the first byte.

## Examples
For each example, the chosen delimitter is `.` and the possible input characters are a-z, A-Z, and 0-9.

   `asdf` > `asdf`
   `asdddf` > `asdddf` (no swap because .3.d is longer than ddd)
   `asdddddf` > `as.5.df`

## Todo
Eliminate the restriction on not having the delimiter in the input by allowing it to be present
if two are specified in a row.