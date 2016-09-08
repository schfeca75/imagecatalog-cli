package imagecatalog

import (
	"regexp"
	"sort"
	"strconv"
)

type ByVersion []string

func (version ByVersion) Len() int {
	return len(version)
}

func (version ByVersion) Swap(i, j int) {
	version[i], version[j] = version[j], version[i]
}

func (version ByVersion) Less(i, j int) bool {
	splitRegexp := regexp.MustCompile("[\\.\\-]")
	vals1 := splitRegexp.Split(version[i], -1)
	vals2 := splitRegexp.Split(version[j], -1)

	i = 0
	// set index to first non-equal ordinal or length of shortest version string
	for i < len(vals1) && i < len(vals2) && vals1[i] == vals2[i] {
		i++
	}
	// compare first non-equal ordinal number
	if i < len(vals1) && i < len(vals2) {
		val1, _ := strconv.Atoi(vals1[i])
		val2, _ := strconv.Atoi(vals2[i])
		return val1 < val2
	}
	// the strings are equal or one string is a substring of the other, then shorter wins
	// e.g. "2.4.2.0" is newer than "2.4.2.0-9999"
	return len(vals2) > len(vals1)
}

func SortCbVersionKeys(vmap map[string]CbImageInfo) []string {
	sortedVersions := []string{}
	for version := range vmap {
		sortedVersions = append(sortedVersions, version)
	}
	sort.Sort(ByVersion(sortedVersions))
	return sortedVersions
}

func SortImVersionKeys(vmap map[string]ImageInfo) []string {
	sortedVersions := []string{}
	for version := range vmap {
		sortedVersions = append(sortedVersions, version)
	}
	sort.Sort(ByVersion(sortedVersions))
	return sortedVersions
}

func SortVersions(varray []string) []string {
	sort.Sort(ByVersion(varray))
	return varray
}
