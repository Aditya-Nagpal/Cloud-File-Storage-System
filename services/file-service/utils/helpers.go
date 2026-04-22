package utils

import (
	"path/filepath"
	"strconv"
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func GetParentPath(userEmail string, parentPath string) string {
	return userEmail + "/" + parentPath
}

func GenerateUniqueID(len int) (string, error) {
	id, err := gonanoid.New(len)
	if err != nil {
		return "", err
	}
	return id, nil
}

func ParseFilename(filename string) (string, string) {
	ext := filepath.Ext(filename) // e.g., ".jpg"
	name := strings.TrimSuffix(filename, ext)
	return name, strings.TrimPrefix(ext, ".") // return "jpg" without the dot
}

func CreateS3Key(userId int64, publicId, extension string) string {
	return "users/" + strconv.FormatInt(userId, 10) + "/files/" + publicId + "." + extension
}
