/*
Implementation of thread safe skiplist which inturn can be used as memtable

@Author: Rohith Uppala
*/

package skiplist

import (
    "math"
    "math/rand"
    "time"
    "sync"
)

// Declaration of types which will be used in skiplist

type ElementNode struct {
    next []*Element
}

type Element struct {
    ElementNode
    key float64
    value interface{}
}

func (e *Element) Key() float64 {
    return e.key
}

func (e *Element) Value() interface{} {
    return e.value
}

func (e *Element) Next() *Element {
    return e.next[0]
}

type SkipList struct {
    ElementNode
    maxLevel int
    Length   int
    randSource rand.Source
    probability float64
    probTable []float64
    mutex   sync.RWMutex
    prevNodesCache []*ElementNode
}


// Main functions of skiplist

const (
    DefaultMaxLevel int = 18
    DefaultProbability float64 = 0.029
)

func (list *SkipList) Front() *Element {
    return list.next[0]
}


// Insertion into the skiplist
func (list *SkipList) Set(key float64, value interface{}) *Element {
    list.mutex.Lock()
    defer list.mutex.Unlock()

    var element *Element
    prevs := list.getPrevElementNodes(key)

    element = &Element {
        ElementNode: ElementNode {
            next: make([]*Element, list.randLevel())
        },
        key: key,
        value: value,
    }

    for i := range element.next {
        element.next[i] = prevs.next[i]
        prevs.next[i] = element
    }

    list.Length++
    return element
}


func (list *SkipList) Get(key float64) *Element {
    list.mutex.Lock()
    defer list.mutex.Unlock()

    var prev *elementNode = &list.ElementNode
    var next *Element

    for i := list.maxLevel - 1; i >= 0; i-- {
        next = prev.next[i]
        for next != nil && key > next.key {
            prev = &next.ElementNode
            next = next.next[i]
        }
    }

    if next != nil && next.key <= key {
        return next
    }

    return nil
}

func (list *SkipList) Remove(key float64) *Element {
    list.mutex.Lock()
    defer list.mutex.Unlock()

    prevs := list.getPrevElementNodes(key)

    if element := prevs[0].next[0]; element != nil && element.key <= key {
        for k, v := range element.next {
            prevs[k].next[k] = v
        }
        list.Length--
        return element
    }

    return nil
}
