package reporter

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/stacktic/dropbox"
)

// DropboxBackend is a struct that stores the Dropbox client and default report storage location
type DropboxBackend struct {
	*dropbox.Dropbox
	StorageLocation string // The absolute path to the location of the Reporter JSON, usually /Apps/Reporter-App/
}

// GetLatestReport searches the storageLocation to find the latest report file.
// It searches based on filename, not on modified or created time, because
// both can be updated after/before the date in the filename.
func (db *DropboxBackend) GetLatestReport() (File, error) {
	var reporterFile File
	metadata, err := db.Metadata(db.StorageLocation, true, false, "", "", 10000)
	if err != nil {
		return reporterFile, err
	}
	var newestTime time.Time
	var newestPath string
	for _, file := range metadata.Contents {
		if strings.Contains(filepath.Base(file.Path), "-reporter-export.json") {
			filenameDate, err := dateForFilename(file.Path)
			if err != nil {
				return reporterFile, err
			}
			if filenameDate.After(newestTime) {
				newestTime = filenameDate
				newestPath = file.Path
			}
		}
	}

	return db.GetReportForPath(newestPath)
}

// GetReportForPath returns a File for the file at the full path specified.
func (db *DropboxBackend) GetReportForPath(filePath string) (File, error) {
	var reporterFile File
	reader, _, err := db.Download(filePath, "", 0)
	if err != nil {
		return reporterFile, err
	}
	defer reader.Close()
	file, readErr := ioutil.ReadAll(reader)
	if readErr != nil {
		return reporterFile, readErr
	}

	metadata, err := db.Metadata(filePath, false, false, "", "", 1)
	if readErr != nil {
		return reporterFile, readErr
	}

	filenameDate, err := dateForFilename(filePath)
	if err != nil {
		return reporterFile, err
	}

	return File{
		Name:             filepath.Base(filePath),
		Path:             filePath,
		Source:           "dropbox",
		ModifiedTime:     time.Time(metadata.Modified),
		TimeFromFilename: filenameDate,
		Contents:         string(file),
	}, nil
}

// GetReportForTime returns a File for the file with the date given in the filename
func (db *DropboxBackend) GetReportForTime(date time.Time) (File, error) {
	filePath := fmt.Sprintf("%s%s-reporter-export.json", db.StorageLocation, date.Format("2006-01-02"))
	return db.GetReportForPath(filePath)
}

// ListReports lists all available reports
func (db *DropboxBackend) ListReports() ([]File, error) {
	var allFiles []File
	metadata, err := db.Metadata(db.StorageLocation, true, false, "", "", 10000)
	if err != nil {
		return allFiles, err
	}
	for _, file := range metadata.Contents {
		if strings.Contains(filepath.Base(file.Path), "-reporter-export.json") {
			filenameDate, err := dateForFilename(file.Path)
			if err != nil {
				return allFiles, err
			}
			allFiles = append(allFiles, File{
				Name:             filepath.Base(file.Path),
				Path:             file.Path,
				Source:           "dropbox",
				ModifiedTime:     time.Time(file.Modified),
				TimeFromFilename: filenameDate,
			})
		}
	}

	return allFiles, nil
}

// NewDropboxBackend returns a new Dropbox backend to read JSON from.
// You must provide an accessToken, which you can get by creating an app
// in the Dropbox API and then pressing Generate.
// Access tokens do not expire.
// If a storageLocation isn't provided, the default location is
//   /Apps/Reporter-App/
func NewDropboxBackend(accessToken, storageLocation string) (*DropboxBackend, error) {
	if accessToken == "" {
		return nil, errors.New("No access token provided for Dropbox backend")
	}
	db := dropbox.NewDropbox()
	db.SetAccessToken(accessToken)
	if storageLocation == "" {
		storageLocation = "/Apps/Reporter-App/"
	}
	return &DropboxBackend{db, storageLocation}, nil
}
