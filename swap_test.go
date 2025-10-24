package main

import (
	"strings"
	"testing"
)

func TestRuneLoop(t *testing.T) {
	cases := []struct {
		runeToOffset, expectedRune rune
		offset                     int32
		description                string
	}{
		{
			runeToOffset: 'a',
			expectedRune: 'a',
			offset:       0,
			description:  "Offset of 0, same letter",
		},
		{
			runeToOffset: 'a',
			expectedRune: 'b',
			offset:       1,
			description:  "Offset of 1, lower case",
		},
		{
			runeToOffset: 'A',
			expectedRune: 'B',
			offset:       1,
			description:  "Offset of 1, upper case",
		},
		{
			runeToOffset: 'b',
			expectedRune: 'a',
			offset:       -1,
			description:  "Offset of -1, upper case",
		},
		{
			runeToOffset: 'B',
			expectedRune: 'A',
			offset:       -1,
			description:  "Offset of -1, upper case",
		},
		{
			runeToOffset: 'Z',
			expectedRune: 'A',
			offset:       1,
			description:  "Offset of 1 on Z, upper case",
		},
		{
			runeToOffset: 'z',
			expectedRune: 'a',
			offset:       1,
			description:  "Offset of 1 on Z, lower case",
		},
		{
			runeToOffset: 'é',
			expectedRune: 'f',
			offset:       1,
			description:  "Offset of 1 on é, lower case",
		},
		{
			runeToOffset: 'É',
			expectedRune: 'F',
			offset:       1,
			description:  "Offset of 1 on É, upper case",
		},
		{
			runeToOffset: 'A',
			expectedRune: 'Z',
			offset:       -1,
			description:  "Offset of 1 on A, upper case",
		},
		{
			runeToOffset: 'a',
			expectedRune: 'z',
			offset:       -1,
			description:  "Offset of 1 on a, lower case",
		},
		{
			runeToOffset: 'ç',
			expectedRune: 'd',
			offset:       1,
			description:  "Offset of 1 on ç, lower case",
		},
		{
			runeToOffset: 'Ç',
			expectedRune: 'D',
			offset:       1,
			description:  "Offset of 1 on Ç, upper case",
		},
		{
			runeToOffset: 'A',
			expectedRune: 'A',
			offset:       26,
			description:  "Offset of 26 on A, upper case",
		},
		{
			runeToOffset: 'A',
			expectedRune: 'A',
			offset:       260,
			description:  "Offset of 260 on A, upper case",
		},
		{
			runeToOffset: 'A',
			expectedRune: 'A',
			offset:       -260,
			description:  "Offset of -260 on A, upper case",
		},
		{
			runeToOffset: 'o',
			expectedRune: 'h',
			offset:       45,
			description:  "Offset of 45 on o, lower case",
		},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			result := swapRune(tt.runeToOffset, tt.offset)

			if result != tt.expectedRune {
				t.Errorf("with %d (%s), got %d (%s), want %d (%s)", tt.runeToOffset, string(tt.runeToOffset), result, string(result), tt.expectedRune, string(tt.expectedRune))
			}
		})
	}
}

func TestEncryptFile(t *testing.T) {
	cases := []struct {
		fileContent, description, expectedFileContent string
	}{
		{
			fileContent:         "> This is a quote line. It should be ignored",
			expectedFileContent: "> This is a quote line. It should be ignored\n",
			description:         "Quote line. Content should stay the same",
		},
		{
			fileContent:         "==ROT 45",
			expectedFileContent: "==ROT 45\n",
			description:         "Rotation line. Content should stay the same",
		},
		{
			fileContent: `
==ROT 2

> Ahaha

Hello world`,
			expectedFileContent: `
==ROT 2

> Ahaha

Jgnnq yqtnf` + "\n",
			description: "Typical exemple of a file. Offsets letters by 2 to the right",
		},
		{
			fileContent: `
==ROT -2

> Ohoho

Hello world`,
			expectedFileContent: `
==ROT -2

> Ohoho

Fcjjm umpjb` + "\n",
			description: "Other typical exemple of a file. Offsets letters by 2 to the left",
		},
	}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			encryption, err := encryptOrDecryptFile(strings.NewReader(tt.fileContent))
			if err != nil {
				t.Errorf("got error %v, but didn't expect one", err)
			}

			if string(encryption) != tt.expectedFileContent {
				t.Errorf("got :\n%s \n want :\n%s \n", string(encryption), tt.expectedFileContent)
			}
		})
	}
}
