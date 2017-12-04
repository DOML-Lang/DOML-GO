/*
 * Copyright (c) 2017. Braedon Wooding
 * Created under LICENSE, see the file LICENSE for information
 */

 /*
 Statistics:
 - It scales up accurately based on character/line count
 	- i.e. A string of length 84 takes 8902 ns to run, 1193 takes 107848, and 2385 takes 211611 on the same machine
 - Thus the ratio is 0.0094361, 0.01106186, 0.01127068
 	- (demonstrating a small fall off due to the fact that a certain percentage is taken up by the initial)
 - A test was done to find a more accurate average by effectively just averaging a large range of the values
 	- The first one was simply from 1 to 5, second was 5 to 50, and third was 1 to 100
 	- Results being 0.0110691, 0.01042827, and 0.01057732
 		- Reversing back the equation you get for a length of 84, 1193, and 2385; 7588.69/8,055/7941.52, 107777.50678/114400.567/112788.495, and 215464.672/228705.24066/225482.447
 		- Thus the accuracy ratio being: 85%/90%/89%, 99.93%/94%/95.6%, 98%/93%/94%
 - Error value for the average of the ratios increases as the length of the block grows
 */

package parser

import (
	"bufio"
	"doml/internal/core"
	"doml/runtime"
	"errors"
	"os"
	"strconv"
	"strings"
	"unicode"
	"fmt"
)

type parsingInformation struct {
	instructionIndex      int
	instructions          []runtime.Instruction
	maximumAmountOfValues int
	currentVariable       string
	nextRegister          int
	registers             map[string]registerData
}

type registerData struct {
	registerID int
	objectType string
}

func ParseDOMLFromFile(filePath string) (runtimeInstance runtime.Runtime, err error, lineNum int) {
	file, err := os.Open(filePath)
	if err != nil {
		return runtime.Runtime{}, err, 0
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	instructions, err, lineCount := parseInstructions(scanner)

	if err != nil {
		return runtime.Runtime{}, err, lineCount
	}

	return runtime.NewEmptyRuntime(instructions), nil, lineCount
}

func ParseDOMLFromString(text string) (runtimeInstance runtime.Runtime, err error, lineNum int) {
	scanner := bufio.NewScanner(strings.NewReader(text))
	instructions, err, lineCount := parseInstructions(scanner)

	if err != nil {
		return runtime.Runtime{}, err, lineCount
	}

	return runtime.NewEmptyRuntime(instructions), nil, lineCount
}

// Allowing only '.' and '->'
func checkIsValidIdentifier(text string) bool {
	var dashPrevious bool

	for _, char := range text {
		if char == '-' && !dashPrevious {
			dashPrevious = true
		} else if char == '>' && dashPrevious {
			dashPrevious = false
		} else if char != '.' && !unicode.IsLetter(char) {
			return false
		} else if dashPrevious {
			return false
		}
	}

	return true
}

// This just parses a set of instructions from a scanner
// Doesn't currently track lines, I don't think its possible with how I setup scanner
// But I could also just do some other way of tracking it to handle it from that situation??
// Furthermore the code is bloated as hell, and we could remove all the bloat by spreading it into different functions
func parseInstructions(scanner *bufio.Scanner) (instructions *[]runtime.Instruction, err error, lineCount int) {
	scanner.Split(DOMLSplitFunc)

	parsingInfo := parsingInformation{2, make([]runtime.Instruction, 2), 0, "", 0, make(map[string]registerData)}
	parsingInfo.instructions[0] = runtime.Instruction{Opcode: core.MakeReg}
	parsingInfo.instructions[1] = runtime.Instruction{Opcode: core.MakeSpace}
	var text string
	var newline bool
	var dontScan bool
	lineCount = 1

	quickScan := func() bool {
		for scanner.Scan() {
			text = scanner.Text()
			if text == "\n" {
				newline = true
				lineCount += 1
			} else {
				return true
			}
		}
		return false
	}

	for dontScan || quickScan() {
		dontScan = false

		if text == "/" {
			if !quickScan() {
				return nil, scanner.Err(), lineCount
			}

			if text == "*" {
				if !quickScan() {
					return nil, scanner.Err(), lineCount
				}

				// Multi-line
				commentCount := 1
				for commentCount > 0 && quickScan() {
					if text == "*" {
						if quickScan() {
							if text == "/" {
								commentCount -= 1
							}
						} else {
							// REF: 1
							// No need to set comment count
							// Or anything else since if we are here
							// Its going to be > 0
							break
						}
					} else if text == "/" {
						if quickScan() {
							if text == "*" {
								commentCount += 1
							}
						} else {
							// Same as REF: 1
							break
						}
					}
				}

				if commentCount > 0 {
					return nil, errors.New("didn't finish multi-line comment"), lineCount
				} else {
					continue
				}
			} else if text == "/" {
				// Singular line
				newline = false
				for scanner.Scan() {
					if text == "\n" {
						newline = true
						lineCount += 1
						break
					}
				}

				if newline == true {
					dontScan = true
					continue
				} else {
					// We are at the end of the buffer
					break
				}
			} else {
				// Error
				return nil, errors.New("invalid character"), lineCount
			}
		}

		if text == "@" {
			// Creation Statement
			if !quickScan() {
				return nil, scanner.Err(), lineCount
			}

			creationVariable := text

			if !quickScan() {
				return nil, scanner.Err(), lineCount
			}

			if !checkIsValidIdentifier(creationVariable) {
				return nil, errors.New("not a valid identifier"), lineCount
			}

			if text == "=" {
				if !quickScan() {
					return nil, scanner.Err(), lineCount
				}

				constructor := text

				if !quickScan() {
					return nil, scanner.Err(), lineCount
				}

				if text == "..." {
					parsingInfo.currentVariable = creationVariable
				}

				parsingInfo.registers[creationVariable] = registerData{registerID: parsingInfo.nextRegister, objectType: constructor}
				parsingInfo.instructions = append(parsingInfo.instructions, runtime.Instruction{Opcode: core.New, Parameter: "new " + constructor},
					runtime.Instruction{Opcode: core.RegObj, Parameter: parsingInfo.nextRegister})
				parsingInfo.nextRegister += 1
			} else {
				return nil, errors.New("missing a '='"), lineCount
			}
		} else if text == ";" || newline {
			// Assignment Statement
			newline = false
			if text == ";" && !quickScan() {
				return nil, scanner.Err(), lineCount
			}

			var assignmentSplit []string

			if text[0] == '.' {
				if parsingInfo.currentVariable != "" {
					assignmentSplit = []string{parsingInfo.currentVariable, text[1:]}
				} else {
					return nil, errors.New("missing a '...' on previous line"), lineCount
				}
			} else if index := strings.Index(text, "."); index != -1 {
				assignmentSplit = []string{text[:index], text[index+1:]}
			} else {
				return nil, errors.New("missing a '.'"), lineCount
			}

			if !quickScan() {
				return nil, scanner.Err(), lineCount
			}

			// Values
			if text == "=" {
				if !quickScan() {
					return nil, scanner.Err(), lineCount
				}

				values := 0
				index := 0

				for {
					if text[0] == '"' {
						index = strings.LastIndex(text[:], "\"")
						if index == -1 {
							return nil, errors.New("missing ending '\"'"), lineCount
						}

						parsingInfo.instructions = append(parsingInfo.instructions, runtime.Instruction{Opcode: core.PushStr, Parameter: text[1:index]})
						values += 1

						if !quickScan() {
							return nil, scanner.Err(), lineCount
						}

						if text != "," {
							dontScan = true
							break
						} else {
							if !quickScan() {
								return nil, scanner.Err(), lineCount
							}

							continue
						}
					} else {
						var value interface{}
						var opcode byte

						if text == "false" {
							value = false
							opcode = core.PushBool
						} else if text == "true" {
							value = true
							opcode = core.PushBool
						} else {
							var characters = []rune(text)
							if unicode.IsDigit(characters[0]) {
								if len(characters) > 1 {
									// Could be integer, number, or binary/octal/hex
									if unicode.ToLower(characters[1]) == 'x' {
										// Hex
										text = text[2:]
										value, err = strconv.ParseInt(text, 16, 64)
										if err != nil {
											return nil, err, lineCount
										}
										opcode = core.PushInt
									} else if unicode.ToLower(characters[1]) == 'b' {
										// Binary
										text = text[2:]
										value, err = strconv.ParseInt(text, 2, 64)
										if err != nil {
											return nil, err, lineCount
										}
										opcode = core.PushInt
									} else if unicode.ToLower(characters[1]) == 'o' {
										// Octal
										text = text[2:]
										value, err = strconv.ParseInt(text, 8, 64)
										if err != nil {
											return nil, err, lineCount
										}
										opcode = core.PushInt
									} else if strings.ContainsRune(text, '.') {
										// Floating point
										value, err = strconv.ParseFloat(text, 64)
										if err != nil {
											return nil, err, lineCount
										}
										opcode = core.PushNum
									} else {
										// Integer
										value, err = strconv.ParseInt(text, 10, 64)
										if err != nil {
											return nil, err, lineCount
										}
										opcode = core.PushInt
									}
								} else {
									// Integer
									value, err = strconv.ParseInt(text, 10, 64)
									if err != nil {
										return nil, err, lineCount
									}
									opcode = core.PushInt
								}
							} else if characters[0] == '$' {
								value, opcode, err = ParseDecimal(text)

								if err != nil {
									return nil, err, lineCount
								}
							}
						}

						parsingInfo.instructions = append(parsingInfo.instructions, runtime.Instruction{Opcode: opcode, Parameter: value})
						values += 1

						if !quickScan() {
							// Sometimes they can end here
							break
						}

						if text != "," {
							dontScan = true
							break
						} else {
							if !quickScan() {
								return nil, scanner.Err(), lineCount
							}

							continue
						}
					}
				}

				objectInfo := parsingInfo.registers[assignmentSplit[0]]
				parsingInfo.instructions = append(parsingInfo.instructions, runtime.Instruction{Opcode: core.PushObj, Parameter: objectInfo.registerID},
																			runtime.Instruction{Opcode: core.Set, Parameter: fmt.Sprintf("set %[1]v::%[2]v", objectInfo.objectType, assignmentSplit[1])})
				if values > parsingInfo.maximumAmountOfValues {
					parsingInfo.maximumAmountOfValues = values
				}
			} else {
				return nil, errors.New("missing '='"), lineCount
			}
		} else {
			return nil, errors.New("invalid character/s: " + text), lineCount
		}
	}

	parsingInfo.instructions[0].Parameter = parsingInfo.nextRegister
	parsingInfo.instructions[1].Parameter = parsingInfo.maximumAmountOfValues

	return &parsingInfo.instructions, scanner.Err(), lineCount
}
