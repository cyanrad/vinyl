package storage

import (
	"fmt"
	"io"
	"log"
	"main/util"
	"net/http"
	"os"
	"path/filepath"
)

func FetchImage(resourceType util.ResourceType, url string, resourceName string) error {
	log.Printf("Fetching image file for %s resource from %s of type %s", resourceName, url, resourceType)
	path, exsists := checkIfImageExists(resourceType, resourceName)
	if exsists {
		log.Printf("File at %s already exists. Skipping fetch", path)
		return nil
	}

	// Send HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create the file
	filePath := filepath.Join(util.MEDIA_PATH, string(resourceType), resourceName+".jpg")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the image data to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}

	log.Printf("File %s fetched successfully", filePath)
	return nil
}

func checkIfImageExists(resourceType util.ResourceType, resourceName string) (string, bool) {
	// Construct file path
	filePath := filepath.Join(util.MEDIA_PATH, string(resourceType), resourceName+".jpg")

	// Check if file exists
	_, err := os.Stat(filePath)
	if err == nil {
		return filePath, true
	}
	return "", false
}
