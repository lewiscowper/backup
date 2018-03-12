package filenames

import (
	"fmt"
	"time"
)

// Get returns the appropriate filenames for the archive and the checksum files
func Get(prefix string) (archiveFilename string, checksumFilename string) {
	currentTime := time.Now().Unix()
	filename := fmt.Sprintf("%v-%v", prefix, currentTime)

	archiveFilename = fmt.Sprintf("%v.tar.gz", filename)
	checksumFilename = fmt.Sprintf("%v.sha512", filename)

	return archiveFilename, checksumFilename
}
