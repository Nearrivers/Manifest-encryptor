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

	lowerCaseE = 'e'
	upperCaseE = 'E'
	lowerCaseC = 'c'
	upperCaseC = 'C'
	lowerCaseU = 'u'
	upperCaseU = 'U'
	lowerCaseI = 'i'
	upperCaseI = 'I'
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
	lowerCaseÊ: lowerCaseE,
	upperCaseÊ: upperCaseE,
	lowerCaseÎ: lowerCaseI,
	upperCaseÎ: lowerCaseI,
	lowerCaseÛ: lowerCaseU,
	upperCaseÛ: upperCaseU,
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

	return encryptedCodex, nil
}
