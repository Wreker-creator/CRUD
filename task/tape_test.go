package task

import (
	"io"
	"testing"
)

func TestTape_Write(t *testing.T) {

	t.Run("check if the data is written properly", func(t *testing.T) {

		file, clean := createTempFile(t, "12345")
		defer clean()

		tape := &tape{file}

		tape.Write([]byte("abc"))

		file.Seek(0, io.SeekStart)
		newFileContents, _ := io.ReadAll(file)

		got := string(newFileContents)
		want := "abc"

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

	})

}
