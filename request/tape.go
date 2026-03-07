package request

import (
	"io"
	"os"
)

/*

	There is a problem with how delete is implemented right now, when we delete something
	we seek to the start of the file again and because of that when data is being written again
	old bytes get left behind

	Basically, overwrite is happening instead of rewriting. Solution?

	Seek to start, truncate the file (basically empty it), rewrite from beginning.

	So instead of each method of FileSytemStore doing database.Seek(), we wrap the file in this tap
	and every write automatically handles the seek and truncate part.

	Doesn't require us to keep repeating something.

*/

type tape struct {
	// file io.ReadWriteSeeker changed to os.File as it contains truncate functionality.
	file *os.File
}

func (t *tape) Write(p []byte) (n int, err error) {
	t.file.Seek(0, io.SeekStart)
	t.file.Truncate(0) // required to reset the file.
	return t.file.Write(p)
}
