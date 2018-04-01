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
	"github.com/cyphar/paperback/pkg/shamir"
)

// backupper is the internal implementation of Backupper.
type backupper struct {
	master     schema.EncryptedMaster
	passphrase []byte
}

// Backupper is used for the construction of the various payloads needed for a
// paperback backup. The key methods are Master() and Shards(), which produce
// the ready-for-document-embedding "master document" and "key shards"
// (respectively).
//
// An example usage of this interface is the following:
// ```
// plaintext := "This is our secret message."
// bc, _ := backup.New([]byte(plaintext))
// master, _ := bc.Master()
// // The shard keywords are included in the wrapped shamir.Shard struct.
// shards, _ := bc.Shards(k, n)
// ```
type Backupper interface {
	Master() schema.EncryptedMaster
	Shards(k, n uint) ([]BundledShard, error)
}

// Create takes a plaintext and constructs a new Backupper which is used to
// generate the necessary master document and shard payloads. The returned
// payloads are effectively ready-for-printing.
func Create(plaintext schema.Master) (Backupper, error) {
	// Generate our passphrase.
	passphrase, err := crypto.GeneratePassphrase(crypto.DefaultPasswordLength)
	if err != nil {
		return nil, err
	}

	// Create the encrypted document now.
	ciphertext, err := plaintext.Encrypt(passphrase)
	if err != nil {
		return nil, err
	}

	return &backupper{
		master:     ciphertext,
		passphrase: passphrase,
	}, nil
}

// Master returns the master document payload, which is symmetrically PGP
// encrypted with a passphrase. The passphrase shards can be generated with
// Shards().
func (bs *backupper) Master() schema.EncryptedMaster {
	return bs.master
}

// Shards constructs a (k,n)-threshold scheme that contains the passphrase of
// the master document. At least k of the n generated shards must be provided
// during recovery to recover the key. Each shard is also individually
// encrypted with some shard codes, which are bundled with the shard here.
func (bs *backupper) Shards(k, n uint) ([]BundledShard, error) {
	// Generate our shards.
	shards, err := shamir.Split(k, n, bs.passphrase)
	if err != nil {
		return nil, err
	}

	// Construct the bundled shards.
	var bundledShards []BundledShard
	for _, shard := range shards {
		schemaShard := schema.NewShard(shard)
		bundledShard, err := NewBundledShard(schemaShard)
		if err != nil {
			return nil, err
		}
		bundledShards = append(bundledShards, bundledShard)
	}
	return bundledShards, nil
}
