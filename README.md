# Smart Run-Length Encoding

Implementation of a smart run-length encoding algorithm. The goal of this implementation is
to optimize for the worst-case usage, where the data is highly random. In these cases, adding
1s before most characters nearly doubles the size of the input, leading to a highly inefficient
compression ratio. In this implementation, only values which have a run length which can be
expressed in less characters than will be removed will be prepended with a run length value. 
This results in a worst-case compression ratio of 1.

The drawback of this approach is that ambiguity is introduced, where digits may be part of a run
length or part of the input data. To solve this, a delimitter must be used. This delimitter must
not be a possible value in the input data. As a result, if there is no character which is outside
of the range of possible characters in the input, this algorithm **will not** work.

## Examples
For each example, the chosen delimitter is `.` and the possible input characters are a-z, A-Z, and 0-9.

`asdf` > `asdf`

`23fffffj` > `23.5.fj`

`23fffj` > `23fffj`