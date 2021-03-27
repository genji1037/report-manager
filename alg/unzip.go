package alg

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

func Unzip(zippedFile *os.File) (io.ReadCloser, error) {
	fi, err := zippedFile.Stat()
	if err != nil {
		return nil, err
	}

	r, err := zip.NewReader(zippedFile, fi.Size())
	if err != nil {
		return nil, err
	}

	if len(r.File) == 0 {
		return nil, fmt.Errorf("no file exists in zip file")
	}

	rd, err := r.File[0].Open()
	if err != nil {
		return nil, fmt.Errorf("open file failed: %v", err)
	}

	return rd, nil
}
