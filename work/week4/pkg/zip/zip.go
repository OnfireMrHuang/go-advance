package zip

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func GzipFile(source string, targzPath string) error {
	out, err := os.Create(targzPath)
	if err != nil {
		return fmt.Errorf("error creating %s: %v", targzPath, err)
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Println(err)
		}
	}()

	gzWriter, err := gzip.NewWriterLevel(out, gzip.BestCompression)
	if err != nil {
		return err
	}

	defer func() {
		if err := gzWriter.Close(); err != nil {
			log.Println(err)
		}
	}()

	fi, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("%s: stat: %v", source, err)
	}

	_, fname := filepath.Split(source)
	gzWriter.Header.Name = fname
	file, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open file %s. error: %s", source, err.Error())
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	_, err = io.CopyN(gzWriter, file, fi.Size())
	if err != nil && err != io.EOF {
		return fmt.Errorf("%s: copying contents: %v", source, err)
	}
	return nil
}
