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

// extender
type extender struct {
}

type Extender interface {
	AddShard(shard BundledShard) error
	Extend()
}

// Extend allows you to extend a set of bundled shards and produce more bundled
// shards that are compatbile with the original set. For a (k,n)-threshold
// scheme you must have at least k shards.
func Extend(shards []BundledShard) ([]BundledShard, error) {
	// TODO
	return nil, nil
}
