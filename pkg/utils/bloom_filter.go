/*
Bloom Filter is a space efficient probabilistic data structure 
which tells you with a probability whether the element is present in set or not
but for sure if it's not

To do: 
Number of Hash Functions (k) and size of bloomfilter is determined by number of keys
we would like to add and also based on how much false positive rate we can accept

@Author: Rohith Uppala
*/

package bloomfilter

import (
    "hash"
    "github.com/spaolacci/murmur3"
    "hash/fnv"
)


type Interface interface {
    AddToBloomFilter(item []byte)
    CheckForElementInBloomFilter(item []byte)
}

const NumberOfHashFunctions = 3

type BloomFilter struct {
    bitset []bool
    k uint // Number of Hash Functions
    NumberOfElements uint // Number of keys in bloom filter
    m uint // size of bloom filter bitset
    hashFuncs []hash.Hash64
}

func New(size uint) *BloomFilter {
    return &BloomFilter {
        bitset: make([]bool, size),
        k: NumberOfHashFunctions,
        NumberOfElements: uint(0),
        m: size,
        hashFuncs: []hash.Hash64{murmur3.New64(), fnv.New64(), fnv.New64a()},
    }
}

func (bf *BloomFilter) AddToBloomFilter(item []byte) {
    hashes := bf.hashValues(item)
    i := uint(0)
    for {
        if i >= bf.k { break }
        position := uint(hashes[i]) % bf.m
        bf.bitset[uint(position)] = true
        i += 1
    }
    bf.NumberOfElements += 1
}

func (bf *BloomFilter) hashValues(item []byte) []uint64 {
    var result []uint64
    for _, hashFunc := range bf.hashFuncs {
        hashFunc.Write(item)
        result = append(result, hashFunc.Sum64())
        hashFunc.Reset()
    }
    return result
}

func (bf *BloomFilter) CheckForElementInBloomFilter(item []byte) (exists bool) {
    hashes := bf.hashValues(item)
    i, exists := uint(0), true
    for {
        if i >= bf.k { break }
        position := uint(hashes[i]) % bf.m
        if !bf.bitset[uint(position)] {
            exists = false
            break
        }
        i += 1
    }
    return
}
