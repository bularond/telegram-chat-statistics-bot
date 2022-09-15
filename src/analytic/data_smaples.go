package analytic

import (
	"sort"
	"time"
)

type PersonList []*Person

func (p PersonList) Len() int           { return len(p) }
func (p PersonList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PersonList) Less(i, j int) bool { return p[i].MessageCount < p[j].MessageCount }

func (cs *ChatStatistics) GetMostActiveProfile(count int) PersonList {
	if count > len(cs.Persons) {
		count = len(cs.Persons)
	}

	persons := make(PersonList, len(cs.Persons))
	pointer := 0
	for _, person := range cs.Persons {
		persons[pointer] = person
		pointer++
	}

	sort.Sort(sort.Reverse(persons))

	return persons[:count]
}

type WordsPair struct {
	Key   string
	Value int
}

type WordsPairList []WordsPair

func (p WordsPairList) Len() int           { return len(p) }
func (p WordsPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p WordsPairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func (cs *ChatStatistics) GetMostPopularWords(count int) WordsPairList {
	if count > len(cs.WorsdMap) {
		count = len(cs.WorsdMap)
	}

	list := make(WordsPairList, 0)
	for key, value := range cs.WorsdMap {
		if len(key) > 6 {
			list = append(list, WordsPair{key, value})
		}
	}

	sort.Sort(sort.Reverse(list))

	return list[:count]
}

type DatePair struct {
	Key   time.Time
	Value int
}

type DatePairList []DatePair

func (p DatePairList) Len() int           { return len(p) }
func (p DatePairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p DatePairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func (cs *ChatStatistics) GetMostPopularDate(count int) DatePairList {
	if count > len(cs.DateMap) {
		count = len(cs.DateMap)
	}

	list := make(DatePairList, len(cs.DateMap))
	pointer := 0
	for key, value := range cs.DateMap {
		list[pointer] = DatePair{key, value}
		pointer++
	}

	sort.Sort(sort.Reverse(list))

	return list[:count]
}
