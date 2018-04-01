/*
 * paperback: resilient paper backups for the very paranoid
 * Copyright (C) 2018 Aleksa Sarai <cyphar@cyphar.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package bip39

import (
	"bytes"
	"fmt"
	"strconv"
)

const byteChunkSize = 8

// chunkBits takes an array of bytes and re-chunks it into the given chunk size
// (in bits). The entire array is treated as big-endian, so the bytes are
// chunked in the same order they would be iterated over.
func chunkBits(data []byte, chunkSize int) []uint64 {
	if chunkSize <= 0 || chunkSize > 64 {
		panic("chunkBits cannot store <=0bit or >64bit chunks")
	}

	// This is all done in the worst possible way. I'm sorry. In future I will
	// implement it properly using bitwise logic -- the reason for using
	// strings here is that the implementation is far more obviously correct
	// than handling weird chunkSize values.

	bitBuilder := new(bytes.Buffer)
	for _, bits := range data {
		_, _ = bitBuilder.WriteString(fmt.Sprintf("%.*b", byteChunkSize, bits))
	}
	bitString := bitBuilder.String()

	var chunks []uint64
	for bitString != "" {
		endIdx := chunkSize
		if endIdx > len(bitString) {
			endIdx = len(bitString)
		}
		chunk, err := strconv.ParseUint(bitString[:endIdx], 2, chunkSize)
		if err != nil {
			panic("generated bits from bitBuilder are not base-2: " + err.Error())
		}
		bitString = bitString[endIdx:]
		chunks = append(chunks, chunk)
	}
	return chunks
}

// combineBits takes a set of bits cut up into the given chunkSize and returns
// the reconstructed bytes.
func combineBits(bits []uint64, chunkSize int) []byte {
	if chunkSize <= 0 || chunkSize > 64 {
		panic("combineBits cannot restore <=0bit or >64bit chunks")
	}

	// As with chunkBits this is done in the worst possible way.

	byteBuilder := new(bytes.Buffer)
	for _, chunk := range bits {
		_, _ = byteBuilder.WriteString(fmt.Sprintf("%.*b", chunkSize, chunk))
	}
	byteString := byteBuilder.String()

	var bytes []byte
	for byteString != "" {
		endIdx := byteChunkSize
		if endIdx > len(byteString) {
			endIdx = len(byteString)
		}
		value, err := strconv.ParseUint(byteString[:endIdx], 2, byteChunkSize)
		if err != nil {
			panic("generated bits from byteBuilder are not base-2: " + err.Error())
		}
		byteString = byteString[endIdx:]
		bytes = append(bytes, byte(value&0xFF))
	}
	return bytes
}
