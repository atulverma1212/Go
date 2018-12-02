package utils

import (
	"sort"
)

// A data structure to hold key/value pairs
type Pair struct {
	Key   string
	Value int
}

// A slice of pairs that implements sort.Interface to sort by values
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func SortMap(data map[string] interface{}) PairList{
	var p PairList
	for k, v := range data {
		p = append(p,Pair{k, v.(int)})
	}
	sort.Sort(p)
	return p
}
