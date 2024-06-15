package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+/"

func main() {
	// Channel to catch Ctrl+C (SIGINT) signal.
	// This allows us to print the best result found so far and the most
	// recent nonce before exiting
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Starting with "a" as the nonce and the worst possible result
	nonce := []byte("a")
	bestResult := "zzzzzz"
	username := "lucasmelin"

	// If a username is provided as an argument, start with that username
	if len(os.Args) > 1 {
		username = os.Args[1]
		fmt.Println("Provided username:", username)
	}
	// If a nonce is provided as an argument, start with that nonce
	if len(os.Args) > 2 {
		nonce = []byte(os.Args[2])
		fmt.Printf("\nStarting with provided nonce: %s", nonce)
	}
	// If the best nonce is provided as an argument, calculate the
	// best result for that nonce and use it as the starting point
	if len(os.Args) > 3 {
		bestNonce := []byte(os.Args[3])
		bestResult = hashUsername(username, bestNonce)
		fmt.Printf("\nPrevious best result: Nonce %s, SHA256: %s\n\n", bestNonce, prettyPrint(bestResult))
	}
	for {
		select {
		case <-sigChan:
			fmt.Println("\nProgram interrupted. Exiting...")
			fmt.Printf("\nBest result found: %s", prettyPrint(bestResult))
			fmt.Printf("\nLast nonce: %s\n", nonce)
			return
		default:
			// Generate next nonce
			nonce = getNextNonce(nonce)
			// Hash the username with the nonce
			hashString := hashUsername(username, nonce)
			if hashString < bestResult {
				bestResult = hashString
				fmt.Printf("Nonce: %s, SHA256: %s\n", nonce, prettyPrint(hashString))
			}
		}
	}
}

// hashUsername hashes the username and nonce with SHA256 and returns the result
func hashUsername(username string, nonce []byte) string {
	user := []byte(username + "/")
	hash := sha256.Sum256(append(user, nonce...))
	return hex.EncodeToString(hash[:])
}

// getNextNonce generates the next nonce string by iterating through the charset
func getNextNonce(currentNonce []byte) []byte {
	// Iterate through nonce string from right to left
	for i := len(currentNonce) - 1; i >= 0; i-- {
		// Directly calculate the index of a character in the charset
		// by subtracting the value of the first character in the charset
		// from the current character.
		index := int(currentNonce[i] - charset[0])
		if index < len(charset)-1 {
			// Increment the rune and return the new nonce
			currentNonce[i] = charset[index+1]
			return currentNonce
		} else {
			// Set the rune to the first in the charset
			currentNonce[i] = charset[0]
		}
	}

	// Add a new rune to the left
	return append([]byte{charset[0]}, currentNonce...)
}

func prettyPrint(s string) string {
	return fmt.Sprint(strings.Join(chunk(s, 8), " "))
}

// chunk splits a string into chunks of a given size
func chunk(s string, chunkSize int) []string {
	if len(s) == 0 {
		return []string{}
	}
	// If the chunk size is greater than the string length, return the string as a single chunk
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string
	currentLen := 0
	currentStart := 0
	for i := range s {
		// If the current chunk is full, add it to the list of chunks
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			// Start a new chunk
			currentStart = i
		}
		currentLen++
	}
	// Add the last chunk
	chunks = append(chunks, s[currentStart:])
	return chunks
}
