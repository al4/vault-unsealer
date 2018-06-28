package env_file

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// Properties - map of key-values
type Properties map[string]string

// ReadPropertiesFile - Read a properties file and return a Properties map
func ReadPropertiesFile(filename string) (Properties, error) {
	properties := Properties{}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := strings.TrimSpace(line[equal+1:])
				properties[key] = value
			}
		}
	}

	scanErr := scanner.Err()

	if scanErr != nil {
		log.Fatal(scanErr)
	}

	return properties, scanErr
}
