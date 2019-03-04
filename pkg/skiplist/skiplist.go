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

    if element = prevs[0].next[0]; element != nil && element.key <= key {
        element.value = value
        return element
    }

    element = &Element {
        ElementNode: ElementNode {
            next: make([]*Element, list.randLevel()),
        },
        key: key,
        value: value,
    }

    for i := range element.next {
        element.next[i] = prevs[i].next[i]
        prevs[i].next[i] = element
    }

    list.Length++
    return element
}


func (list *SkipList) Get(key float64) *Element {
    list.mutex.Lock()
    defer list.mutex.Unlock()

    var prev *ElementNode = &list.ElementNode
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

func (list *SkipList) getPrevElementNodes(key float64) []*ElementNode {
    var prev *ElementNode = &list.ElementNode
    var next *Element

    prevs := list.prevNodesCache

    for i := list.maxLevel -1 ; i >= 0; i-- {
        next = prev.next[i]
        for next != nil && key > next.key {
            prev = &next.ElementNode
            next = next.next[i]
        }
        prevs[i] = prev
    }

    return prevs
}

func (list *SkipList) randLevel() (level int) {
    r := float64(list.randSource.Int63()) / (1 << 63)
    level = 1
    for level < list.maxLevel && r < list.probTable[level] {
        level++
    }
    return level
}

func probabilityTable(probability float64, maxLevel int) (table []float64) {
    for i := 1; i <= maxLevel; i++ {
        prob := math.Pow(probability, float64(i - 1))
        table = append(table, prob)
    }
    return table
}

func NewWithMaxLevel(maxLevel int) *SkipList {
    if maxLevel < 1 || maxLevel > 64 {
        panic("Maxlevel of skiplist must be a positive integer <= 64")
    }

    return &SkipList {
        ElementNode: ElementNode{ next: make([]*Element, maxLevel) },
        prevNodesCache: make([]*ElementNode, maxLevel),
        maxLevel: maxLevel,
        randSource: rand.New(rand.NewSource(time.Now().UnixNano())),
        probability: DefaultProbability,
        probTable: probabilityTable(DefaultProbability, maxLevel),
    }
}

func New() *SkipList {
    return NewWithMaxLevel(DefaultMaxLevel)
}
