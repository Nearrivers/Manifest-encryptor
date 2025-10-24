package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	lowerCaseA = 97
	lowerCaseZ = 122
	upperCaseA = 65
	upperCaseZ = 90

	lowerCaseÉ = 233
	upperCaseÉ = 201
	lowerCaseÀ = 224
	upperCaseÀ = 192
	lowerCaseÇ = 231
	upperCaseÇ = 199
	lowerCaseÈ = 232
	upperCaseÈ = 200

	lowerCaseE = 101
	upperCaseE = 69
	lowerCaseC = 99
	upperCaseC = 67
)

var equivalent = map[rune]rune{
	lowerCaseÉ: lowerCaseE,
	upperCaseÉ: upperCaseE,
	lowerCaseÈ: lowerCaseE,
	upperCaseÈ: upperCaseE,
	lowerCaseÀ: lowerCaseA,
	upperCaseÀ: upperCaseA,
	lowerCaseÇ: lowerCaseC,
	upperCaseÇ: upperCaseC,
}

// TODO: Les -1 ne fonctionnent toujours pas (alors que les -20 oui ?)
func swap(r, lowerLimit, upperLimit rune, offset int32) rune {
	newLetter := r + offset

	var remainder int32
	if offset < 0 {
		remainder = lowerLimit % newLetter
	} else {
		remainder = newLetter % upperLimit
	}

	if remainder < newLetter && remainder != 0 {
		if offset > 0 {
			r = lowerLimit + remainder - 1
		} else {
			r = upperLimit - (remainder - 1)
		}

		return r
	}

	return r + offset
}

func swapRune(r rune, offset int32) rune {
	e, ok := equivalent[r]
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

func main() {
	err := filepath.WalkDir("./", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// NOTE: À activer une fois le programme 100% fonctionnel
		// if !strings.Contains(d.Name(), "ntrée") {
		// 	return nil
		// }

		if d.IsDir() || filepath.Ext(d.Name()) != ".md" || strings.Contains(d.Name(), "Codex") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		encryptedCodex, err := encryptFile(file)
		if err != nil {
			return err
		}

		f, err := os.Create(fmt.Sprintf("Codex %s", filepath.Base(d.Name())))
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

func encryptFile(file io.Reader) ([]byte, error) {
	scanner := bufio.NewScanner(file)

	// NOTE: Inside my files, there will be lines that start with ==ROT n
	// with n a number. This number will be used to offset letters found inside
	// the file
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

		// If
		if r.MatchString(line) {
			submatches := r.FindStringSubmatch(line)
			offset, err = strconv.Atoi(submatches[len(submatches)-2])
			if err != nil {
				return []byte{}, err
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

	fmt.Println(string(encryptedCodex))

	return encryptedCodex, nil
}
