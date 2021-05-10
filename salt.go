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
	// byte array containing the salt value, should always have a length of 4,
	// values should have a max value of 15.
	bytes []byte
}

// FromByteArray creates a new instance of the Salt from the byte array.
// Length of bytes array should always be 2.
func FromByteArray(bytes []byte) (*Salt, error) {
	if len(bytes) != 2 {
		return nil, fmt.Errorf("invalid salt")
	}
	return &Salt{bytes: bytesAsSalt(bytes)}, nil
}

// FromBase64 creates a new instance of Salt from a Base 64 representation. Turn
// the Base 64 string into a byte array of length 2 and then unpack the 4
// nibbles from the 2 bytes.
func FromBase64(data string) (*Salt, error) {
	b, err := base64.RawStdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return FromByteArray(b)
}

// For a given index 1-16, return whether the visual representation relating to
// that index should be shown.
func (s Salt) Show(i int) bool {
	for _, v := range s.bytes {
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
	for n, v := range s.bytes {
		if v == byte(i-1) {
			is = append(is, strconv.Itoa(n+1))
		}
	}
	return strings.Join(is[:], " ")
}

// Get the Salt value as a byte array.
func (s Salt) GetBytes() []byte {
	return saltAsBytes(s.bytes)
}

// Get the Base 64 representation of the Salt value.
func (s Salt) ToBase64String() string {
	// Get the byte array.
	b := saltAsBytes(s.bytes)
	// Convert the byte array to a Base64 string.
	return base64.RawStdEncoding.EncodeToString(b)
}

// Get the salt value from a byte array.
func bytesAsSalt(b []byte) []byte {
	// Unpack the 4 nibbles from the two bytes. For each byte, bit shift the
	// byte to the right by 4 to get the first nibble, apply the bitwise AND
	// operator against the byte and a value of 16 (0xF) to get the first 4 bits
	// for the second nibble.
	n1, n2 := b[0]>>4, b[0]&0xF
	n3, n4 := b[1]>>4, b[1]&0xF

	return []byte{n1, n2, n3, n4}
}

// Get a salt value as bytes.
func saltAsBytes(data []byte) []byte {
	// Pack the 4 nibbles representing the selected items into two bytes.
	// Bit shift the first nibble to the left by 4 and the OR the result with
	// the first nibble to produce the first byte. Repeat for the third and
	// fourth nibble to create the second byte.
	b1 := (data[0] << 4) | data[1]
	b2 := (data[2] << 4) | data[3]

	// Return the array.
	return []byte{b1, b2}
}
