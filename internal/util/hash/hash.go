package hash

import "hash/fnv"

func NormalizedHash(expression string) int {
	h := fnv.New32a()
	h.Write([]byte(expression))
	return int(h.Sum32() % 100)
}
