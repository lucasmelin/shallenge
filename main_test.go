package main

import (
	"testing"
)

func TestHashUsername(t *testing.T) {
	result := hashUsername("lucasmelin", "aaaaaaaaaaaaabVMavi")
	expected := "0000000008abe5171e888ba62fef3b9794d8aa59c71e17922bd4407fc2f80e19"
	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

func TestGetNextNonce(t *testing.T) {
	result := getNextNonce("abc")
	expected := "abd"
	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

func TestGetNextNonceWithCarry(t *testing.T) {
	result := getNextNonce("///")
	expected := "aaaa"
	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

func TestPrettyPrint(t *testing.T) {
	result := prettyPrint("abcdefghijklmnopqrstuvwxzyz0123456789+/")
	expected := "abcdefgh ijklmnop qrstuvwx zyz01234 56789+/"
	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

func TestChunk(t *testing.T) {
	result := chunk("abcdefgh", 4)
	expected := []string{"abcd", "efgh"}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %s but got %s", expected[i], v)
		}
	}
}

func TestChunkWithUnevenSize(t *testing.T) {
	result := chunk("abcdefgh", 3)
	expected := []string{"abc", "def", "gh"}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Expected %s but got %s", expected[i], v)
		}
	}
}

func BenchmarkNonceCalculation(b *testing.B) {
	nonce := "a"
	for i := 0; i < b.N; i++ {
		nonce = getNextNonce(nonce)
	}
}

func BenchmarkHashCalculation(b *testing.B) {
	nonce := "a"
	for i := 0; i < b.N; i++ {
		nonce = getNextNonce(nonce)
		hashUsername("lucasmelin", nonce)
	}
}
