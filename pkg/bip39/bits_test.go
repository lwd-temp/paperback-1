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
	"crypto/rand"
	"io"
	"testing"
)

// TestRoundtrip just does a bunch of randomised round-trip testing.
func TestRoundtrip(t *testing.T) {
	const nTrials = 16
	for L := 0; L < 32; L++ {
		input := make([]byte, L)
		for trial := 0; trial < nTrials; trial++ {
			if L > 0 {
				n, err := io.ReadFull(rand.Reader, input)
				if n != L || err != nil {
					t.Errorf("failed to read random (L=%d,trial=%d): %v", L, trial, err)
					continue
				}
			}

			for chunkSize := 1; chunkSize < 16; chunkSize++ {
				chunks := chunkBits(input, chunkSize)
				// Make sure the chunks make sense.
				for idx, chunk := range chunks {
					if chunk > (1 << uint(chunkSize)) {
						t.Errorf("chunked slice %v has incorrect chunk at (L=%d,trial=%d,chunkSize=%d,idx=%d): %v", input, L, trial, chunkSize, idx, chunk)
						continue
					}
				}
				output := combineBits(chunks, chunkSize)
				if !bytes.Equal(input, output) {
					t.Errorf("chunking round-trip failed (L=%d,trial=%d,chunkSize=%d): expected %v got %v", L, trial, chunkSize, input, output)
				}
			}
		}
	}
}
