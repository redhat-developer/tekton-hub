/*
Package compress provides functions for compressing files
*/
package compress

import (
	"archive/zip"
	"io"
	"log"
	"os"
)

// ZipFiles creates a new zip file and adds files[]
func ZipFiles(filename string, files []os.FileInfo, dir string) (*os.File, error) {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer newZipFile.Close()
	// Create a zip file
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file, dir); err != nil {
			return nil, err
		}
	}
	var a []byte
	newZipFile.Read(a)
	return newZipFile, nil
}

// AddFileToZip adds a given file to zip
func AddFileToZip(zipWriter *zip.Writer, file os.FileInfo, dir string) error {
	fileToZip, err := os.Open(dir + "/" + file.Name())
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = file.Name()

	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	log.Println()
	return err
}
