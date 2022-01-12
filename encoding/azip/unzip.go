package azip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return nil, err
			}
			continue
		}

		// Make File
		if err2 := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err2 != nil {
			return filenames, err2
		}

		outFile, err3 := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err3 != nil {
			return filenames, err
		}

		rc, err3 := f.Open()
		if err3 != nil {
			return filenames, err3
		}

		_, err4 := io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		if err5 := outFile.Close(); err5 != nil {
			return nil, err5
		}
		rc.Close()

		if err4 != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
