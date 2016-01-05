// Package bulk has functions to load and save image metadata to a CSV file as
// needed by the 'Bulk Add CSV' Shimmie2 extension. The CSV file is assumed to
// have the following format:
//
//  "/path/to/image.jpg","spaced tags","source","rating s/q/e","/path/thumbnail.jpg"
//
// The last record (thumbnail) is left empty as thumbnails can easily be
// generated by the server.
//
// The package assumes that all images and the CSV file are under a certain
// directory path that is used as input in many package functions.
package bulk

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Image holds the metadata of each image from the CSV file.
type Image struct {
	ID     int
	Name   string
	Tags   Tags
	Source string
	Rating string
}

// Tags are the space separated tags for each image. It's used (instead of
// []string) to print the tags easier in the HTML template with a custom
// String() method.
type Tags []string

func (t Tags) String() string {
	return strings.Join(t, " ")
}

var supportedExt = []string{"gif", "jpeg", "jpg", "png", "swf"}

func isSupportedType(name string) bool {
	fname := strings.ToLower(name)
	for _, ext := range supportedExt {
		// The only possible returned error is ErrBadPattern, when pattern is
		// malformed. Patterns like *.jpg are never malformed so we ignore the
		// error.
		matches, _ := filepath.Match("*."+ext, fname)
		if matches {
			return true
		}
	}
	return false
}

// LoadImages takes the path of a folder as input and scans it for files with
// one of the following extensions: "gif", "jpeg", "jpg", "png", "swf"
//
// It returns a slice of images without metadata, using the filename as Name
// and the order the files were found as an increasing ID starting from 0.
//
// In case of a CSV file, the image metadata can be read using LoadCSV and
// combined with the images (discovered by LoadImages) using Combine.
func LoadImages(dir string) ([]Image, error) {
	var images []Image

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	id := 0
	for _, f := range files {
		if !f.IsDir() {
			if isSupportedType(f.Name()) {
				img := Image{ID: id, Name: f.Name()}
				images = append(images, img)
				id++
			}
		}
	}
	return images, nil
}

// LoadCSV loads the image metadata from a CSV file that is open for reading.
// The metadata are returned as slice of images and should be combined with the
// slice of images discovered by LoadImages by calling Combine.
func LoadCSV(file io.Reader) ([]Image, error) {
	images := []Image{}

	r := csv.NewReader(file)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(record) != 5 {
			return nil, fmt.Errorf("invalid csv file format")
		}
		// Image filepath (first column) should exist otherwise we cannot match
		// the metadata with the images found under the directory.
		if record[0] != "" {
			img := Image{
				Name:   filepath.Base(record[0]),
				Tags:   strings.Split(record[1], " "),
				Source: record[2],
				Rating: record[3],
			}
			images = append(images, img)
		}
	}
	return images, nil
}

// Combine takes the metadata of imagesWithInfo and copies them to images
// returning the combined result.
func Combine(images, imagesWithInfo []Image) []Image {
	for _, info := range imagesWithInfo {
		if info.Name == "" {
			continue
		}
		img := findByName(images, info.Name)
		if img != nil {
			img.Source = info.Source
			img.Rating = info.Rating
			img.Tags = info.Tags
		}
	}
	return images
}

type byName []Image

func (img byName) Len() int           { return len(img) }
func (img byName) Swap(i, j int)      { img[i], img[j] = img[j], img[i] }
func (img byName) Less(i, j int) bool { return img[i].Name < img[j].Name }

func findByName(image []Image, name string) *Image {
	sort.Sort(byName(image))
	i := sort.Search(len(image), func(i int) bool { return image[i].Name >= name })
	if i < len(image) && image[i].Name == name {
		return &image[i]
	}
	return nil
}

type byID []Image

func (img byID) Len() int           { return len(img) }
func (img byID) Swap(i, j int)      { img[i], img[j] = img[j], img[i] }
func (img byID) Less(i, j int) bool { return img[i].ID < img[j].ID }

// FindByID takes a slice of images, sorts them by ID and then finds then one
// with ID id.
func FindByID(image []Image, id int) *Image {
	sort.Sort(byID(image))
	i := sort.Search(len(image), func(i int) bool { return image[i].ID >= id })
	if i < len(image) && image[i].ID == id {
		return &image[i]
	}
	return nil
}

// CurrentPrefix reads from an open CSV file, reads the first line and uses the
// base directory of the provided dir path to find the path prefix used in the
// CSV file.
//
// As an example if the provided dir path is '/localpath/pics' and the first line
// has '/serverpath/pics/pic1' then the returned current prefix will be
// '/serverpath'.
func CurrentPrefix(dir string, file io.Reader) (string, error) {
	r := csv.NewReader(file)
	record, err := r.Read()
	if err == io.EOF {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	folder := filepath.Base(dir)
	sep := fmt.Sprintf("%c", filepath.Separator)
	prefix := sep
	parts := strings.Split(record[0], sep)
	for _, p := range parts {
		if p == folder {
			break
		}
		prefix = filepath.Join(prefix, p)
	}
	return prefix, nil
}

// Save will write the image metadata to a CSV file with csvFilename under the
// dir path. It will keep the base of the dir path and replace the prefix with
// the provided one.
func Save(images []Image, dir, csvFilename, prefix string) error {
	csvFile := filepath.Join(dir, csvFilename)
	f, err := os.Create(csvFile)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	w.WriteAll(toRecords(images, dir, prefix))

	if err := w.Error(); err != nil {
		return fmt.Errorf("error writing csv: %v", err)
	}
	return nil
}

func toRecords(images []Image, dir, prefix string) [][]string {
	var records [][]string
	for _, img := range images {
		record := toRecord(img, dir, prefix)
		records = append(records, record)
	}
	return records
}

func toRecord(img Image, dir, prefix string) []string {
	var record []string
	record = append(record, filepath.Join(prefix, filepath.Base(dir), img.Name))
	record = append(record, strings.Join(img.Tags, " "))
	record = append(record, img.Source)
	record = append(record, img.Rating)
	record = append(record, "")
	return record
}
