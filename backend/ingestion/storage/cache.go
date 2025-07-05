package storage

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"main/util"
	"os"
	"path/filepath"
)

const SEP = "-"

type Cache struct {
	filepath string
}

func NewCache(filepath string) *Cache {
	c := Cache{
		filepath: filepath,
	}

	return &c
}

func (c *Cache) Get(resource util.ResourceType, ingestion util.IngestionSource, key string, object any) (bool, error) {
	if object == nil {
		return false, errors.New("nil object passed to Cache.ReadObject")
	}

	key = generateCacheKey(resource, ingestion, key)
	filePath := filepath.Join(c.filepath, key) + ".gob"

	// Check if file exists
	_, err := os.Stat(filePath)
	if err != nil {
		return false, nil // the file doesn't exist
	}

	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(object); err != nil {
		return true, err
	}

	return true, nil
}

func (c *Cache) Delete(resource util.ResourceType, ingestion util.IngestionSource, key string) error {
	key = generateCacheKey(resource, ingestion, key)
	filePath := filepath.Join(c.filepath, key) + ".gob"

	// Check if file exists
	_, err := os.Stat(filePath)
	if err != nil {
		return errors.New("calling Cache.Delete on a file that doesn't exist: " + filePath)
	}

	err = os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) Store(resource util.ResourceType, ingestion util.IngestionSource, key string, value any) error {
	key = generateCacheKey(resource, ingestion, key)

	// writing gob file for easy reconstruction
	if err := c.writeObject(key, value); err != nil {
		return err
	}

	// writing human readable json for people to have their data
	// this file is never read by the program after it's creation
	if util.SOURCE == util.SOURCE_SPOTIFY {
		if err := c.writeJSON(key, value); err != nil {
			return err
		}
	}

	return nil
}

func generateCacheKey(resource util.ResourceType, ingestion util.IngestionSource, key string) string {
	return string(resource) + SEP + string(ingestion) + SEP + key
}

func (c *Cache) writeObject(cacheKey string, object any) error {
	file, err := os.Create(filepath.Join(c.filepath, cacheKey) + ".gob")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(object); err != nil {
		return err
	}

	return nil
}

func (c *Cache) writeJSON(cacheKey string, object any) error {
	file, err := os.Create(filepath.Join(c.filepath, cacheKey) + ".json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(object); err != nil {
		return err
	}

	return nil
}
