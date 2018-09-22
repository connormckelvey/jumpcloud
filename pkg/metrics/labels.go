package metrics

import (
	"hash/fnv"
	"sort"
)

type labelMap map[string]string
type hashedLabels = uint64

func (ls labelMap) keys() (keys []string) {
	for k := range ls {
		keys = append(keys, k)
	}
	return
}

func (ls labelMap) hash() hashedLabels {
	keys := ls.keys()
	sort.Strings(keys)

	hasher := fnv.New64()
	for _, k := range keys {
		hasher.Write([]byte(k))
		hasher.Write([]byte(ls[k]))
	}
	return hashedLabels(hasher.Sum64())
}
