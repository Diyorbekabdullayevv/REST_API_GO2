package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func GetUserID(path, route string) (int, bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 && parts[0] == route {
		userID, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("Failed to convert ID to integer value:", err)
			return 0, false
		}
		return userID, true
	}
	return 0, false
}
