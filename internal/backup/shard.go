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

package backup

import (
	"github.com/cyphar/paperback/internal/schema"
	"github.com/cyphar/paperback/pkg/crypto"
)

// BundledShard is a shard bundled with its shard codes, used purely to wrap
// the two main pieces needed for each key shard when it is output into a
// document.
type BundledShard struct {
	Shard    schema.EncryptedShard
	Codeword schema.Codeword
}

// NewBundledShard takes a shard, generates a new set of codewords and then
// encrypts the shard. The returned shard is bundled with the set of codewords.
func NewBundledShard(shard schema.Shard) (BundledShard, error) {
	// Generate our new passphrase.
	cw, err := crypto.GenerateCodeword(crypto.DefaultCodewordLength)
	if err != nil {
		return BundledShard{}, err
	}

	// Embed the shard inside a schema.Shard and encrypt it.
	encryptedShard, err := shard.Encrypt(cw)
	if err != nil {
		return BundledShard{}, err
	}

	// We're done.
	return BundledShard{
		Shard:    encryptedShard,
		Codeword: cw,
	}, nil
}

// ToShard takes an encrypted and bundled shard, and then returns the original
// (unencrypted) shard or an error.
func (bs BundledShard) ToShard() (schema.Shard, error) {
	return bs.Shard.Decrypt(bs.Codeword)
}
