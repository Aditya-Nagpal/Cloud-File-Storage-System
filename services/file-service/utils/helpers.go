package utils

import (
	"path/filepath"
	"strconv"
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func GenerateUniqueID(len int) (string, error) {
	id, err := gonanoid.New(len)
	if err != nil {
		return "", err
	}
	return id, nil
}

func ParseFilename(filename string) (string, string) {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	return name, strings.TrimPrefix(ext, ".")
}

func CreateS3Key(userId int64, publicId, extension string) string {
	return "users/" + strconv.FormatInt(userId, 10) + "/files/" + publicId + "." + extension
}
