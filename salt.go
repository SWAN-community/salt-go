/* ****************************************************************************
 * Copyright 2021 51 Degrees Mobile Experts Limited (51degrees.com)
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 * ***************************************************************************/

package salt

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

// Salt struct.
type Salt struct {
	Data []byte // byte array containing the salt value, should always have a
	// length of 4
}

// Create a new instance of Salt from a Base 64 representation. Turn the Base 64
// string into a byte array of length 2 and then unpack the 4 nibbles from the
// 2 bytes.
func NewSalt(data string) (*Salt, error) {
	var s Salt

	b, err := base64.RawStdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	// Length of binary data should always be 2.
	if len(b) != 2 {
		return nil, fmt.Errorf("invalid salt")
	}

	// Unpack the 4 nibbles from the two bytes. For each byte, bit shift the
	// byte to the right by 4 to get the first nibble, apply the bitwise AND
	// operator against the byte and a value of 16 (0xF) to get the first 4 bits
	// for the second nibble.
	n1, n2 := b[0]>>4, b[0]&0xF
	n3, n4 := b[1]>>4, b[1]&0xF

	s.Data = []byte{n1, n2, n3, n4}

	return &s, nil
}

// For a given index 1-16, return whether the visual representation relating to
// that index should be shown.
func (s Salt) Show(i int) bool {
	for _, v := range s.Data {
		if v == byte(i-1) {
			return true
		}
	}
	return false
}

// For a given index 1-16, return a number or numbers as a string indicating the
// order of selection. If the index was not selected, then an empty string.
func (s Salt) Number(i int) string {
	var is []string
	for n, v := range s.Data {
		if v == byte(i-1) {
			is = append(is, strconv.Itoa(n+1))
		}
	}
	return strings.Join(is[:], " ")
}

// Get the Base 64 representation of the Salt value.
func (s Salt) ToBase64String() string {

	// Pack the 4 nibbles representing the selected items into two bytes.
	// Bit shift the first nibble to the left by 4 and the OR the result with
	// the first nibble to produce the first byte. Repeat for the third and
	// fourth nibble to create the second byte.
	b1 := (s.Data[0] << 4) | s.Data[1]
	b2 := (s.Data[2] << 4) | s.Data[3]

	// Form the array.
	b := []byte{b1, b2}

	// Convert the byte array to a Base64 string.
	return base64.RawStdEncoding.EncodeToString(b)
}
