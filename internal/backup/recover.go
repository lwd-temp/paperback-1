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
)

// recoverer implements the Recoverer interface.
type recoverer struct {
	master schema.EncryptedMaster
	shards []schema.Shard
}

// Recoverer is used for the reconstruction of the plaintext master document
// using all of the various payloads created for a paperback backup. The key
// methods are AddShard() and Recover().
//
// An example usage of this interface is the following:
// ```
// rc, _ := backup.Recover([]byte(ciphertext))
// _ = rc.AddShard(BundledShard{Shard: shardA, Codeword: cwA})
// _ = rc.AddShard(BundledShard{Shard: shardB, Codeword: cwB})
// _ = rc.AddShard(BundledShard{Shard: shardC, Codeword: cwC})
// plaintext, _ := rc.Recover()
// ```
//
// The reason why the above interface for this package isn't just the trivial
// backup.Recover(ciphertext, shards...) is because, by splitting out the
// AddShard methods, this library can provide much more semantic feedback on
// each shard. If a shard fails to decrypt or (in future) if it is badly
// signed, then AddShard will return an error *immediately* thus providing
// semantic feedback about which shard was "bad".
type Recoverer interface {
	// AddShard adds a shard to the given recovery instance. An error is
	// returned if the shard is invalid, malformed, or otherwise doesn't match
	// the other shards being used in this recovery process.
	//
	// XXX: When we add signing to shards, AddShard may produce an error for
	//      *valid* shards if a maliciously created shard was added first.
	//      Maybe we should include that information in the EncryptedMaster?
	AddShard(shard BundledShard) error

	// Recover returns the decrypted master, or an error if the current
	// recovery process is not in a state where recovery can happen. If an
	// error is returned by Recover(), the caller MAY call Recover() again.
	Recover() (schema.Master, error)
}

// Recover takes an encrypted master and returns a Recoverer which is used for
// reconstruction of the plaintext master document. Each individual shard is
// added separately.
func Recover(master schema.EncryptedMaster) (Recoverer, error) {
	return &recoverer{
		master: master,
		shards: nil,
	}, nil
}

// AddShard adds a shard to the given recovery instance. An error is returned
// if the shard is invalid, malformed, or otherwise doesn't match the other
// shards being used in this recovery process.
func (r *recoverer) AddShard(shard BundledShard) error {
	plainShard, err := shard.ToShard()
	if err != nil {
		return err
	}
	r.shards = append(r.shards, plainShard)
	return nil
}

// Recover returns the decrypted master, or an error if the current recovery
// process is not in a state where recovery can happen. If an error is returned
// by Recover(), the caller MAY call Recover() again.
func (r *recoverer) Recover() (schema.Master, error) {
	// TODO
	return nil, nil
}
