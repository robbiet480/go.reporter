package reporter

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

// FilesystemBackend is a struct that stores the default report storage location
type FilesystemBackend struct {
	storageLocation string // The absolute path to the location of the Reporter JSON, usually ~/Dropbox/Apps/Reporter-App/
}

// GetLatestReport searches the storageLocation to find the latest report file.
// It searches based on filename, not on modified or created time, because
// both can be updated after/before the date in the filename.
func (fs *FilesystemBackend) GetLatestReport() (File, error) {
	var reporterFile File
	files, err := ioutil.ReadDir(fs.storageLocation)
	if err != nil {
		return reporterFile, err
	}
	var latestDate time.Time
	var latestFile os.FileInfo
	for _, file := range files {
		if strings.Contains(file.Name(), "-reporter-export.json") {
			filenameDate, err := dateForFilename(file.Name())
			if err != nil {
				return reporterFile, err
			}
			if filenameDate.After(latestDate) {
				latestDate = filenameDate
				latestFile = file
			}
		}
	}
	filePath := filepath.Join(fs.storageLocation, latestFile.Name())
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return reporterFile, err
	}
	return File{
		Name:             latestFile.Name(),
		Path:             filePath,
		Source:           "filesystem",
		ModifiedTime:     latestFile.ModTime(),
		TimeFromFilename: latestDate,
		Contents:         string(fileContents),
	}, nil
}

// GetReportForPath returns a File for the file at the full path specified.
func (fs *FilesystemBackend) GetReportForPath(path string) (File, error) {
	var reporterFile File
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return reporterFile, err
	}
	osOpen, err := os.Open(path)
	if err != nil {
		return reporterFile, err
	}
	fileStat, err := osOpen.Stat()
	if err != nil {
		return reporterFile, err
	}
	filenameDate, err := dateForFilename(path)
	if err != nil {
		return reporterFile, err
	}
	return File{
		Name:             fileStat.Name(),
		Path:             path,
		Source:           "filesystem",
		ModifiedTime:     fileStat.ModTime(),
		TimeFromFilename: filenameDate,
		Contents:         string(file),
	}, nil
}

// GetReportForTime returns a File for the file with the date given in the filename
func (fs *FilesystemBackend) GetReportForTime(date time.Time) (File, error) {
	fileName := fmt.Sprintf("%s-reporter-export.json", date.Format("2006-01-02"))
	filePath := filepath.Join(fs.storageLocation, fileName)
	return fs.GetReportForPath(filePath)
}

// ListReports lists all available reports
func (fs *FilesystemBackend) ListReports() ([]File, error) {
	var allFiles []File
	files, err := ioutil.ReadDir(fs.storageLocation)
	if err != nil {
		return allFiles, err
	}
	for _, file := range files {
		if strings.Contains(file.Name(), "-reporter-export.json") {
			filenameDate, err := dateForFilename(file.Name())
			if err != nil {
				return allFiles, err
			}
			filePath := filepath.Join(fs.storageLocation, file.Name())
			singleFile := File{
				Name:             file.Name(),
				Path:             filePath,
				Source:           "filesystem",
				ModifiedTime:     file.ModTime(),
				TimeFromFilename: filenameDate,
			}
			allFiles = append(allFiles, singleFile)
		}
	}
	return allFiles, nil
}

// NewFilesystemBackend returns a new local filesystem backend to read JSON from.
// If a storageLocation isn't provided, the default location is
//   ~/Dropbox/Apps/Reporter-App/
func NewFilesystemBackend(storageLocation string) (*FilesystemBackend, error) {
	if storageLocation == "" {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}
		storageLocation = filepath.Join(usr.HomeDir, "Dropbox/Apps/Reporter-App/")
	}
	return &FilesystemBackend{storageLocation}, nil
}
