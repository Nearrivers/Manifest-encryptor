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
	lowerCaseÙ = 'ù'
	upperCaseÙ = 'Ù'
	lowerCaseË = 'ë'
	upperCaseË = 'Ë'
	lowerCaseÏ = 'ï'
	upperCaseÏ = 'Ï'

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
	lowerCaseÙ: lowerCaseU,
	upperCaseÙ: lowerCaseU,
	lowerCaseË: lowerCaseE,
	upperCaseË: upperCaseE,
	lowerCaseÏ: lowerCaseI,
	upperCaseÏ: upperCaseI,
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
	if offset == 0 || offset%26 == 0 {
		return r
	}

	if offset > 26 {
		offset = offset % 26
	}

	if offset > 26 {
		offset = (offset % 26) * -1
	}

	e, ok := equivalents[r]
	if ok {
		r = e
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

		_, err = f.WriteString(encryptedCodex)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func encryptFile(file io.Reader, isReverse bool) (string, error) {
	scanner := bufio.NewScanner(file)

	r := regexp.MustCompile(`(==ROT (-?(\d*)))`)

	var offset int
	var err error

	builder := strings.Builder{}

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			builder.WriteString("\n")
			continue
		}

		// We skip quotes lines
		if line[0] == '>' {
			builder.WriteString(line)
			builder.WriteString("\n")
			continue
		}

		if r.MatchString(line) {
			submatches := r.FindStringSubmatch(line)
			offset, err = strconv.Atoi(submatches[len(submatches)-2])
			if err != nil {
				return "", err
			}

			if isReverse {
				offset *= -1
			}

			builder.WriteString(line)
			builder.WriteString("\n")
			continue
		}

		for _, char := range line {
			builder.WriteRune(swapRune(char, int32(offset)))
		}

		builder.WriteString("\n")
	}

	if err = scanner.Err(); err != nil {
		return "", err
	}

	return builder.String(), nil
}
