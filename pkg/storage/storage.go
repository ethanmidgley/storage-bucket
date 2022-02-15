package storage

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ethanmidgley/storage-bucket/pkg/config"
)

// CheckBucket will make check to see if there is an active bucket
func CheckBucket() bool {
	_, err := os.Stat(config.Conf.PathPrefix + "/" + config.Conf.Yaml.Bucket.Location)
	return !os.IsNotExist(err)
}

// CreateBucket will make a folder for the bucket according to config.yaml
func CreateBucket() error {
	if err := os.Mkdir(config.Conf.PathPrefix+"/"+config.Conf.Yaml.Bucket.Location, 0755); err != nil {
		return err
	}
	log.Printf("Bucket created at %s", config.Conf.Yaml.Bucket.Location)
	return nil
}

// FileExist will check to see if a file exists with the key provided
func FileExist(key string) bool {

	_, err := os.Stat(fmt.Sprintf("%s/%s/%s", config.Conf.PathPrefix, config.Conf.Yaml.Bucket.Location, key))
	return !os.IsNotExist(err)

}

// Delete will delete a file from the bucket location
func Delete(key string) error {

	return os.Remove(fmt.Sprintf("%s/%s/%s", config.Conf.PathPrefix, config.Conf.Yaml.Bucket.Location, key))

}

// Export will take the bucket location and will zip it so it can be reused in another location
func Export() error {

	if config.Conf.Yaml.Bucket.Export.Allowed {
		files, err := ioutil.ReadDir(config.Conf.PathPrefix + "/" + config.Conf.Yaml.Bucket.Location)
		if err != nil {
			return err
		}

		if len(files) == 0 {
			return errors.New("cannot export empty bucket")
		}

		if config.Conf.Yaml.Bucket.Export.Compression == "tar" {

			newArchive, err := os.Create(fmt.Sprintf("%s/%s.tar.gz", config.Conf.PathPrefix, config.Conf.Yaml.Bucket.Location))
			if err != nil {
				return err
			}
			defer newArchive.Close()

			gw := gzip.NewWriter(newArchive)
			defer gw.Close()
			tw := tar.NewWriter(gw)
			defer tw.Close()

			for _, file := range files {
				err := addToArchive(tw, fmt.Sprintf("%s/%s/%s", config.Conf.PathPrefix, config.Conf.Yaml.Bucket.Location, file.Name()))
				if err != nil {
					return err
				}
			}

		} else {

			newZipFile, err := os.Create(fmt.Sprintf("%s/%s.zip", config.Conf.PathPrefix, config.Conf.Yaml.Bucket.Location))
			if err != nil {
				return err
			}
			defer newZipFile.Close()

			zipwriter := zip.NewWriter(newZipFile)
			defer zipwriter.Close()

			for _, file := range files {
				err := addFileToZip(zipwriter, fmt.Sprintf("%s/%s/%s", config.Conf.PathPrefix, config.Conf.Yaml.Bucket.Location, file.Name()))
				if err != nil {
					return err
				}
			}

		}

	}
	return nil
}
