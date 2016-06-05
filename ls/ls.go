package ls

import (
	"io/ioutil"
	"mime"
	"path"
	"path/filepath"
)

type File struct {
	Filename    string
	Etag        string
	ContentType string
}

func GetFiles(directory string) ([]File, error) {
	var output []File
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return output, err
	}
	for _, file := range files {
		name := path.Join(directory, file.Name())
		if file.IsDir() {
			subFiles, err := GetFiles(name)
			if err != nil {
				return output, err
			}
			output = append(output, subFiles...)
		} else {
			output = append(output, File{
				Filename:    name,
				Etag:        "TODO",
				ContentType: mime.TypeByExtension(filepath.Ext(name)),
			})
		}
	}
	return output, nil
}
