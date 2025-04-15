package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

func main() {
	// Parse command-line flags
	filePath := flag.String("file", "", "Path to the MessagePack file (required)")
	removeKeys := flag.String("remove", "", "Comma-separated keys to remove (optional)")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Error: The '-file' flag is required.")
		os.Exit(1)
	}

	// Read the file
	data, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Error reading file '%s': %v", *filePath, err)
	}

	// Decode MessagePack data
	var msgMap map[string]string
	err = msgpack.Unmarshal(data, &msgMap)
	if err != nil {
		log.Fatalf("Error decoding MessagePack from file '%s': %v", *filePath, err)
	}

	// Process keys to remove
	keysToRemove := make(map[string]bool)
	if *removeKeys != "" {
		keys := strings.SplitSeq(*removeKeys, ",")
		for key := range keys {
			keysToRemove[strings.TrimSpace(key)] = true
		}
	}

	// Remove entries or inspect content
	if len(keysToRemove) > 0 {
		originalSize := len(msgMap)
		for key := range keysToRemove {
			delete(msgMap, key)
		}
		if len(msgMap) < originalSize {
			// Write the updated map back to the file
			encodedData, err := msgpack.Marshal(msgMap)
			if err != nil {
				log.Fatalf("Error encoding MessagePack: %v", err)
			}

			// Create a backup of the original file
			timestamp := time.Now().Format(time.RFC3339)
			backupFilePath := *filePath + ".bak" + timestamp

			err = os.WriteFile(backupFilePath, data, 0644)
			if err != nil {
				log.Fatalf("Error creating backup file '%s': %v", backupFilePath, err)
			}
			fmt.Printf("Created backup file: %s\n", backupFilePath)

			// Write the updated data to the original file
			err = os.WriteFile(*filePath, encodedData, 0644)
			if err != nil {
				log.Fatalf("Error writing updated file '%s': %v", *filePath, err)
			}
			fmt.Printf("Removed %d entries and updated file '%s'.\n", originalSize-len(msgMap), *filePath)
		} else {
			fmt.Println("No matching keys found to remove.")
		}
	} else {
		// Inspect content
		fmt.Printf("Contents of '%s':\n", *filePath)
		for key, value := range msgMap {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
}
