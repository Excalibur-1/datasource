package sort

import (
	"fmt"
	"testing"
)

func TestSorted(t *testing.T) {
	sorted := Sorted().Desc("id").Asc("user").Desc("code")
	fmt.Println(sorted.ToString())
	fmt.Println(sorted.FirstAscString())
}

func TestSortedProperty(t *testing.T) {
	prop := []string{"id", "user", "code"}
	sorted := Sorted()
	for k, v := range prop {
		if k/2 == 0 {
			sorted.Desc(v)
		} else {
			sorted.Asc(v)
		}
	}
	fmt.Println(sorted.ToString())
	fmt.Println(sorted.FirstAscString())
}
