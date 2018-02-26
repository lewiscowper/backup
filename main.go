package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func addFile(tw *tar.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if stat, err := file.Stat(); err == nil {
		// now lets create the header as needed for this file within the tarball
		header := new(tar.Header)
		header.Name = path
		header.Size = stat.Size()
		header.Mode = int64(stat.Mode())
		header.ModTime = stat.ModTime()
		// write the header to the tarball archive
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// copy the file data to the tarball
		if _, err := io.Copy(tw, file); err != nil {
			return err
		}
	}
	return nil
}

func computeMd5(filePath string) ([]byte, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}

	result = h.Sum(nil)

	return result, nil
}

func convertHashToString(hash []byte) string {
	return hex.EncodeToString(hash[:16])
}

func main() {
	// get the current date for timestamp purposes
	current_time := time.Now().Unix()
	filename := fmt.Sprintf("%v-%v", "backup", current_time)

	// set up the output file
	file, err := os.Create(fmt.Sprintf("%v.tar.gz", filename))
	if err != nil {
		log.Fatalln(err)
	}
	// set up the gzip writer
	gw := gzip.NewWriter(file)
	tw := tar.NewWriter(gw)
	// grab the paths that need to be added in
	paths := os.Args[1:]

	// add each file as needed into the current tar archive
	for i := range paths {
		if err := addFile(tw, paths[i]); err != nil {
			log.Fatalln(err)
		}
	}

	// close our writers in reverse order before calculating checksum
	tw.Close()
	gw.Close()
	file.Close()

	if b, err := computeMd5(fmt.Sprintf("%v.tar.gz", filename)); err != nil {
		fmt.Printf("Err: %v", err)
	} else {
		checksumFile, err := os.Create(fmt.Sprintf("%v.md5", filename))
		if err != nil {
			log.Fatalln(err)
		}
		defer checksumFile.Close()

		hash := convertHashToString(b)

		checksumFile.WriteString(hash)

		log.Printf("Backup created with filename %v.tar.gz, and checksum %v", filename, hash)
	}
}
