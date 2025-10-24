package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	lowerCaseA = 'a'
	lowerCaseZ = 'z'
	upperCaseA = 'A'
	upperCaseZ = 'Z'

	lowerCaseÉ = 'é'
	upperCaseÉ = 'É'
	lowerCaseÀ = 'à'
	upperCaseÀ = 'À'
	lowerCaseÇ = 'ç'
	upperCaseÇ = 'Ç'
	lowerCaseÈ = 'è'
	upperCaseÈ = 'È'

	// Fuck le passé simple
	lowerCaseÊ = 'ê'
	upperCaseÊ = 'Ê'
	lowerCaseÎ = 'î'
	upperCaseÎ = 'Î'
	lowerCaseÛ = 'û'
	upperCaseÛ = 'Û'
	lowerCaseÂ = 'â'
	upperCaseÂ = 'Â'

	lowerCaseE = 'e'
	upperCaseE = 'E'
	lowerCaseC = 'c'
	upperCaseC = 'C'
	lowerCaseU = 'u'
	upperCaseU = 'U'
	lowerCaseI = 'i'
	upperCaseI = 'I'
)

var equivalents = map[rune]rune{
	lowerCaseÉ: lowerCaseE,
	upperCaseÉ: upperCaseE,
	lowerCaseÈ: lowerCaseE,
	upperCaseÈ: upperCaseE,
	lowerCaseÀ: lowerCaseA,
	upperCaseÀ: upperCaseA,
	lowerCaseÇ: lowerCaseC,
	upperCaseÇ: upperCaseC,
	lowerCaseÊ: lowerCaseE,
	upperCaseÊ: upperCaseE,
	lowerCaseÎ: lowerCaseI,
	upperCaseÎ: upperCaseI,
	lowerCaseÛ: lowerCaseU,
	upperCaseÛ: upperCaseU,
	lowerCaseÂ: lowerCaseA,
	upperCaseÂ: upperCaseA,
}

func swap(r, lowerLimit, upperLimit rune, offset int32) rune {
	newLetter := r + offset

	if offset < 0 && newLetter < lowerLimit {
		return upperLimit - (lowerLimit - newLetter - 1)
	}

	if offset > 0 && newLetter > upperLimit {
		return lowerLimit + (newLetter - upperLimit - 1)
	}

	return newLetter
}

func swapRune(r rune, offset int32) rune {
	e, ok := equivalents[r]
	if ok {
		r = e
	}

	switch {
	case offset == 0:
		return r
	case offset > 26:
		offset = offset % 26
	case offset < -26:
		offset = (offset % 26) * -1
	}

	if r >= upperCaseA && r <= upperCaseZ {
		return swap(r, upperCaseA, upperCaseZ, offset)
	}

	if r >= lowerCaseA && r <= lowerCaseZ {
		return swap(r, lowerCaseA, lowerCaseZ, offset)
	}

	return r
}

// Walks the current dir in search of markdown files that respect the pattern
func main() {
	if len(os.Args) == 1 {
		log.Fatal("This program needs a file pattern (regex not supported)")
	}

	pattern := os.Args[1]

	err := filepath.WalkDir("./", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !strings.Contains(d.Name(), pattern) || d.IsDir() || filepath.Ext(d.Name()) != ".md" {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		// If pattern == "Codex", then the program will decrypt its own output
		// This is achieved by inverting the offset everytime
		isReverse := pattern == "Codex"

		encryptedCodex, err := encryptFile(file, isReverse)
		if err != nil {
			return err
		}

		fileName := "Codex"
		if isReverse {
			fileName = "Origine"
		}
		f, err := os.Create(fmt.Sprintf("%s %s", fileName, filepath.Base(d.Name())))
		if err != nil {
			return err
		}

		_, err = f.Write(encryptedCodex)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func encryptFile(file io.Reader, isReverse bool) ([]byte, error) {
	scanner := bufio.NewScanner(file)

	r := regexp.MustCompile(`(==ROT (-?(\d*)))`)

	var offset int
	var err error

	encryptedCodex := make([]byte, 0)

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			encryptedCodex = append(encryptedCodex, byte('\n'))
			continue
		}

		// We skip quotes lines
		if line[0] == '>' {
			encryptedCodex = append(encryptedCodex, []byte(line)...)
			encryptedCodex = append(encryptedCodex, byte('\n'))
			continue
		}

		if r.MatchString(line) {
			submatches := r.FindStringSubmatch(line)
			offset, err = strconv.Atoi(submatches[len(submatches)-2])
			if err != nil {
				return []byte{}, err
			}

			if isReverse {
				offset *= -1
			}

			encryptedCodex = append(encryptedCodex, []byte(line)...)
			encryptedCodex = append(encryptedCodex, byte('\n'))
			continue
		}

		for _, char := range line {
			encryptedCodex = append(encryptedCodex, byte(swapRune(char, int32(offset))))
		}

		encryptedCodex = append(encryptedCodex, byte('\n'))
	}

	if err = scanner.Err(); err != nil {
		return []byte{}, err
	}

	return encryptedCodex, nil
}
