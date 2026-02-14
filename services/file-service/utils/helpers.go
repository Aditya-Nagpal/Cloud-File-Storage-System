package utils

func GetParentPath(userEmail string, parentPath string) string {
	return userEmail + "/" + parentPath
}
