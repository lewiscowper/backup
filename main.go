package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"time"
)

func addFile(tw *tar.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	fi, err := os.Lstat(path)
	if err != nil {
		log.Fatal(err)
	}

	switch mode := fi.Mode(); {
	case mode.IsRegular():
		log.WithFields(log.Fields{
			"file": path,
		}).Info("Adding file")

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

	case mode.IsDir():
		filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			addFile(tw, filePath)
			return nil
		})

		return nil
	case mode&os.ModeSymlink != 0:
		log.WithFields(log.Fields{
			"file": path,
		}).Info("Not adding file, it's a symbolic link")
		return nil
	case mode&os.ModeNamedPipe != 0:
		log.WithFields(log.Fields{
			"file": path,
		}).Info("Not adding file, it's a named pipe")
		return nil
	}

	return nil
}

func computeMd5(filePath string) (string, error) {
	var result string
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		return result, err
	}

	checksum := h.Sum(nil)

	result = hex.EncodeToString(checksum)

	return result, nil
}

func createCheckSum(filename string, hash string) error {
	checksumFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer checksumFile.Close()

	checksumFile.WriteString(hash)

	return nil
}

func createArchive(files []string, archiveFilename string) error {
	file, err := os.Create(archiveFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for i := range files {
		if err := addFile(tw, files[i]); err != nil {
			return err
		}
	}

	return nil
}

func getFilenames(prefix string) (archiveFilename string, checksumFilename string) {
	currentTime := time.Now().Unix()
	filename := fmt.Sprintf("%v-%v", prefix, currentTime)

	archiveFilename = fmt.Sprintf("%v.tar.gz", filename)
	checksumFilename = fmt.Sprintf("%v.md5", filename)

	return archiveFilename, checksumFilename
}

func main() {
	files := os.Args[1:]

	archiveFilename, checksumFilename := getFilenames("backup")

	if err := createArchive(files, archiveFilename); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error creating archive")
	}

	if hash, err := computeMd5(archiveFilename); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error getting checksum")
	} else {
		err := createCheckSum(checksumFilename, hash)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Error creating checksum file")
		}

		log.WithFields(log.Fields{
			"filename": archiveFilename,
		}).Info("Created archive")
	}
}
