package archive

import (
	"archive/tar"
	"compress/gzip"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/openpgp"
	"io"
	"os"
	"path/filepath"
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

// Create takes the files in and archives them with tar and gzip
func Create(files []string, archiveFilename string, password []byte) error {
	file, err := os.Create(archiveFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	pgpw, err := openpgp.SymmetricallyEncrypt(file, password, &openpgp.FileHints{IsBinary: true}, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error encrypting with pgp")
	}
	defer pgpw.Close()

	gw := gzip.NewWriter(pgpw)
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
