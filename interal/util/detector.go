package util

import (
	"log"
	"os"
	"strings"
)

func IsAWebserver(filePath string) bool {

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Error while processing the file")
	}
	content := strings.ToLower(string(data))

	keywords := []string{
		"net/http", "http.listenandserve", "router", "handlefunc",
		"gin.default", "echo.new", "fiber.new", "mux.newrouter",
	}

	for _, keyword := range keywords {
		if strings.Contains(content, keyword) {
			return true
		}
	}

	return false

}
