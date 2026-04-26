package utils

import (
	"fmt"
	"time"
)

func CalculateAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()

	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		age--
	}

	return age
}

func CreateS3KeyForDisplayPicture(userId int64) string {
	return fmt.Sprintf("displayPictures/users/%d/", userId)
}
