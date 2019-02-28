package main

import (
    "fmt"
)

func main() {
    var benchList *SkipList
    benchList = New()
    for i := 0; i <= 10000000; i++ {
        benchList.Set(float64(i), [1]byte{})
    }
    var sl SkipList
    var el Element
    fmt.Printf("Structure sizes: Skiplist is %v, Element is %v bytes\n", unsafe.Sizeof(sl), unsafe.Sizeof(el))
}

