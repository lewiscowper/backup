package main // import "github.com/lewiscowper/backup"

import (
	"github.com/lewiscowper/backup/archive"
	"github.com/lewiscowper/backup/checksum"
	"github.com/lewiscowper/backup/credentials"
	"github.com/lewiscowper/backup/filenames"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	files := os.Args[1:]

	archiveFilename, checksumFilename := filenames.Get("backup")

	password, err := credentials.Get()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error capturing password")
	}

	if err := archive.Create(files, archiveFilename, password); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error creating archive")
	}

	if err := checksum.CreateFile(archiveFilename, checksumFilename); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error creating checksum")
	}

	log.WithFields(log.Fields{
		"filename": archiveFilename,
	}).Info("Created archive")
}
