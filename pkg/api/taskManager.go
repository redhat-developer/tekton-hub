package api

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/Pipelines-Marketplace/backend/pkg/compress"
)

// GetCompressedFiles returns the created zip file of requestedTask
func GetCompressedFiles(requestedTask string) *os.File {
	dir := "catalog" + "/" + requestedTask
	requestedFiles, err := ioutil.ReadDir("catalog" + "/" + requestedTask + "/")
	if err != nil {
		log.Fatal(err)
	}
	finalZipFile, err := compress.ZipFiles("finalZipFile.zip", requestedFiles, dir)
	if err != nil {
		log.Println(err)
	}
	return finalZipFile
}
