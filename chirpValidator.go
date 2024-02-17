package main

import (
	"errors"
	"slices"
	"strings"
)

func validateChirp(chirp string) (string, error) {
	if len(chirp) > 140 {
		err := errors.New("chirp is too long")
		return "", err
	}

	censoredSentance := censorWord(chirp)

	return censoredSentance, nil
}

func censorWord(stringToCensor string) string {
	badWordsDictionary := map[string]string{
		"kerfuffle": "*********",
		"sharbert":  "********",
		"fornax":    "******",
	}

	censoredSlice := make([]string, 0)
	words := strings.Split(stringToCensor, " ")

mainLoop:
	for _, word := range words {
		for badWord, replacement := range badWordsDictionary {
			if badWord == strings.ToLower(word) {
				censoredSlice = slices.Insert(censoredSlice, len(censoredSlice), replacement)
				continue mainLoop
			}
		}

		censoredSlice = slices.Insert(censoredSlice, len(censoredSlice), word)
	}

	finalString := strings.Join(censoredSlice, " ")

	return finalString
}
