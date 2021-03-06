package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func isCommentString(s string) bool {
	return (strings.HasPrefix(s, "\"\"\"") || strings.HasPrefix(s, "'''"))
}

// IsEmptyOrUndocumented returns true if the array contains
// an empty comment or if it's three lines and contains the
// "Undocumented" comment. The lastString parameter is the last
// string read before the comment was encountered.
func IsEmptyOrUndocumented(lastString string, cb []string) bool {

	// Exceptions can contain an empty string. Won't catch subclass names but it's good enough.
	if strings.Contains(lastString, "Exception") || strings.Contains(lastString, "Error") {
		return false
	}

	for idx, i := range cb {
		im := strings.TrimSpace(i)
		if idx == 1 && isCommentString(im) && len(strings.TrimSpace(cb[idx-1])) == 3 {
			return true
		}
	}

	if len(cb) == 3 {
		switch strings.TrimSpace(cb[1]) {
		case "Undocumented":
			return true
		case "Undocumented.":
			return true
		}
	}

	return false
}

// StripEmptyOrIrrelevantComments scans the reader into a string array
// buffer and returns it, sans empty comments and one line comments
// containing the "Undocumented" comment.
// TODO: investigate replacing the buffer output with an io.Writer parameter
func StripEmptyOrIrrelevantComments(r io.Reader) (buffer []string) {
	scanner := bufio.NewScanner(r)
	var s string
	var m string
	dbuf := []string{}
	shouldFix := false
	var lastString string

	for scanner.Scan() {
		s = scanner.Text()
		m = strings.TrimSpace(s)
		if isCommentString(m) {
			dbuf = append(dbuf, s)
			if len(dbuf) > 1 {
				isEmpty := IsEmptyOrUndocumented(lastString, dbuf)
				if isEmpty {
					shouldFix = true
				} else {
					for _, i := range dbuf {
						buffer = append(buffer, i)
					}
				}
				dbuf = dbuf[:0]
			}
			continue
		}
		lastString = m

		if len(dbuf) > 0 {
			dbuf = append(dbuf, s)
			continue
		}
		buffer = append(buffer, s)
	}

	if !shouldFix {
		return nil
	}

	return
}

func visit(path string, i os.FileInfo, err error) error {

	m, err := filepath.Match("*.py", i.Name())

	if err != nil || !m {
		return err
	}

	f, err := os.OpenFile(path, os.O_RDWR, i.Mode())
	defer f.Close()

	if err != nil {
		panic(err)
	}

	buffer := StripEmptyOrIrrelevantComments(f)
	if buffer == nil {
		return nil
	}

	if _, err := f.Seek(0, os.SEEK_SET); err != nil {
		panic(err)
	}

	w := bufio.NewWriter(f)
	for _, s := range buffer {
		fmt.Fprintln(w, s)
	}
	w.Flush()
	offset, _ := f.Seek(0, os.SEEK_CUR)
	f.Truncate(offset)

	return nil
}

func main() {
	err := filepath.Walk(".", visit)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
