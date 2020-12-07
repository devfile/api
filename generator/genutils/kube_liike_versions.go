package genutils

import (
	"sort"

	"k8s.io/apimachinery/pkg/version"
)



type sortByKubeLikeVersion []string
		
func (a sortByKubeLikeVersion) Len() int           { return len(a) }
func (a sortByKubeLikeVersion) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByKubeLikeVersion) Less(i, j int) bool { return version.CompareKubeAwareVersionStrings(a[i], a[j]) < 0 }

func SortKubeLikeVersion(versions []string) {
	sort.Sort(sortByKubeLikeVersion(versions))
}

func LatestKubeLikeVersion(versions []string) string {
	if len(versions) == 0 {
		return ""
	}
	versionsToSort := versions
	sort.Sort(sortByKubeLikeVersion(versionsToSort))
	return versionsToSort[len(versionsToSort)-1]	
}
