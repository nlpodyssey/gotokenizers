// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bpemodel

// DefaultCacheCapacity is the default capacity for BpeModel internal cache.
const DefaultCacheCapacity = 10_000

type WordCache struct {
	capacity int
}

// NewCache returns a new Cache initialized with the given capacity.
//
// If capacity is set to zero, the cache becomes ineffective (is disabled).
func NewCache(capacity int) *WordCache {
	return &WordCache{
		capacity: capacity,
	}
}

// NewDefaultCache returns a new Cache initialized with the default capacity.
func NewDefaultCache() *WordCache {
	return NewCache(DefaultCacheCapacity)
}

func (c *WordCache) SetValues(keys []string, values []*Word) {
}

func (c *WordCache) GetValues(keys []string) []*Word {
	words := make([]*Word, len(keys))
	for i := range keys {
		words[i] = nil
	}
	return words
}
