package ls

import (
	"io/ioutil"
	"mime"
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
		if !file.IsDir() {
			filename := file.Name()
			output = append(output, File{
				Filename:    filename,
				Etag:        "TODO",
				ContentType: mime.TypeByExtension(filepath.Ext(filename)),
			})
		}
	}
	return output, nil
}
