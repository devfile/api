package genutils

import (
	"sort"

	"k8s.io/apimachinery/pkg/version"
)

type sortByKubeLikeVersion []string

func (a sortByKubeLikeVersion) Len() int      { return len(a) }
func (a sortByKubeLikeVersion) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sortByKubeLikeVersion) Less(i, j int) bool {
	return version.CompareKubeAwareVersionStrings(a[i], a[j]) < 0
}

// SortKubeLikeVersion sorts the provided versions according to "kube-like" versioning order.
// "Kube-like" versions start with a "v", then are followed by a number (the major version),
// then optionally the string "alpha" or "beta" and another number (the minor version). These are sorted first
// by GA > beta > alpha (where GA is a version with no suffix such as beta or alpha), and then by comparing
// major version, then minor version. An example sorted list of versions:
// v10, v2, v1, v11beta2, v10beta3, v3beta1, v12alpha1, v11alpha2, foo1, foo10.
func SortKubeLikeVersion(versions []string) {
	sort.Sort(sortByKubeLikeVersion(versions))
}

// LatestKubeLikeVersion retrieves the latest version from the the provided versions, according to "kube-like" versioning order.
// "Kube-like" versions start with a "v", then are followed by a number (the major version),
// then optionally the string "alpha" or "beta" and another number (the minor version). These are sorted first
// by GA > beta > alpha (where GA is a version with no suffix such as beta or alpha), and then by comparing
// major version, then minor version. An example sorted list of versions:
// v10, v2, v1, v11beta2, v10beta3, v3beta1, v12alpha1, v11alpha2, foo1, foo10.
func LatestKubeLikeVersion(versions []string) string {
	if len(versions) == 0 {
		return ""
	}
	versionsToSort := make([]string, len(versions))
	copy(versionsToSort, versions)
	sort.Sort(sortByKubeLikeVersion(versionsToSort))
	return versionsToSort[len(versionsToSort)-1]
}
