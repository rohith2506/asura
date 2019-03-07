package bloomfilter

import (
    "testing"
    "fmt"
)

func TestBasicBloomFilterFunctionality(t *testing.T) {
    var testBloomFilter *BloomFilter
    testBloomFilter = New(100)

    for i := 0; i <= 100; i++ {
        testBloomFilter.AddToBloomFilter([]byte{'r'})
    }
    for i := 101; i <= 110; i++  {
        exists := testBloomFilter.CheckForElementInBloomFilter([]byte{'0'})
        fmt.Println(exists)
        if !exists {
            t.Fatalf("Implementation is wrong")
        }
    }
}
