package objects

import (
	"mime"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Location     string
	Key          string
	Etag         string
	ACL          string
	CacheControl string
	ContentType  string
}

func GetFiles(files []string, strip int, destination string, cacheControl string, acl string) ([]File, error) {
	var output []File
	for _, file := range files {
		name := filepath.Clean(file)
		fileInfo, err := os.Stat(name)
		if err != nil {
			return output, err
		}
		if !fileInfo.IsDir() {
			output = append(output, File{
				ACL:          acl,
				CacheControl: cacheControl,
				ContentType:  mime.TypeByExtension(filepath.Ext(name)),
				Etag:         "TODO",
				Key:          filepath.Join(destination, StripFromName(name, strip)),
				Location:     name,
			})
		}
	}
	return output, nil
}

func StripFromName(name string, strip int) string {
	return strings.Join(strings.Split(name, "/")[strip:], "/")
}
