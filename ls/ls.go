package ls

import (
	"io/ioutil"
	"mime"
	"path/filepath"
	"strings"
)

type File struct {
	Location    string
	Key         string
	Etag        string
	ContentType string
}

func GetFiles(directory string, strip int, destination string) ([]File, error) {
	directory = filepath.Clean(directory)
	var output []File
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return output, err
	}
	for _, file := range files {
		name := filepath.Join(directory, file.Name())
		if file.IsDir() {
			subFiles, err := GetFiles(name, strip, destination)
			if err != nil {
				return output, err
			}
			output = append(output, subFiles...)
		} else {
			output = append(output, File{
				Key:         filepath.Join(destination, StripFromName(name, strip)),
				Location:    name,
				Etag:        "TODO",
				ContentType: mime.TypeByExtension(filepath.Ext(name)),
			})
		}
	}
	return output, nil
}

func StripFromName(name string, strip int) string {
	return strings.Join(strings.Split(name, "/")[strip:], "/")
}
