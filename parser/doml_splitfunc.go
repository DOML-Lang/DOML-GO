/*
 * Copyright (c) 2017. Braedon Wooding
 * Created under LICENSE, see the file LICENSE for information
 */

package parser

import (
	"unicode/utf8"
	"unicode"
)

func DOMLSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	dotCount := 0
	stringMode := false

	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		// Handles a \n in the preceding spaces
		if r == '\n' {
			return start + width, []byte{'\n'}, nil
		} else if !unicode.IsSpace(r) {
			break
		}
	}

	// Scan until space or symbol, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if stringMode {
			if r == '"' {
				return i + width, data[start:i+width], nil
			}
		} else {
			switch r {
			case '"':
				if i != start {
					return i, data[start:i], nil
				} else {
					stringMode = true
				}
			case '.':
				// This handles the '...' case
				dotCount += 1
				if dotCount == 3 {
					if i != start + width * 2 {
						return i - width, data[start:i-width], nil
					} else {
						return i + width, data[start:i+width], nil
					}
				}
				continue
			case ',', '\n', '@', ';', '=', '/', '*', '\\':
				// Handling cases where we spilt;
				// - Newlines for line counters
				// - ',' '@' ';' '=' for easier parsing
				// OtherCount exists to determine if we can just pass the current character
				// Or if we have to pass the previous characters then come back next time to pass this character
				if i != start {
					return i, data[start:i], nil
				} else {
					return i + width, []byte{byte(r)}, nil
				}
			default:
				// If space pass the range excluding the space but to advance including the space
				if unicode.IsSpace(r) {
					return i + width, data[start:i], nil
				}
			}
		}
	}

	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}

	// Request more data.
	return start, nil, nil
}