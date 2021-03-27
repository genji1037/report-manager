package alg

import (
	"fmt"
	"io"
	"os"
	"report-manager/logger"
)

func WriteTemp(rd io.Reader, filename string) (*os.File, error) {
	dir := os.TempDir()
	err := os.MkdirAll(dir, 0644)
	if err != nil {
		return nil, fmt.Errorf("mkdir failed: %v", err)
	}
	path := dir + filename
	logger.Infof("start write tmp file %s", path)
	defer logger.Infof("finish write tmp file %s", path)
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("open file %s failed: %v", path, err)
	}
	bufsize := 1024
	buf := make([]byte, bufsize)
	for {
		n, err := rd.Read(buf)
		if err != nil {
			if err == io.EOF {
				_, err = f.Write(buf[:n])
				if err != nil {
					return nil, fmt.Errorf("write file failed: %v", err)
				}
				return f, nil
			}
			return nil, fmt.Errorf("read from rd failed: %v", err)
		}
		_, err = f.Write(buf[:n])
		if err != nil {
			return nil, fmt.Errorf("write file failed: %v", err)
		}
	}
}
