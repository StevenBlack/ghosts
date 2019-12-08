package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/juliangruber/go-intersect"
)

func main() {

	// for now, hard-coded hosts files in the filesystem
	h1 := "/Users/Steve/Dropbox/dev/hosts/data/StevenBlack/hosts"
	h2 := "/Users/Steve/Dropbox/dev/hosts/data/yoyo.org/hosts"

	h1bytes, err := ioutil.ReadFile(h1)
	h2bytes, err := ioutil.ReadFile(h2)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	host1lines := strings.Split(string(h1bytes), "\n")
	host1lines = normalize(host1lines)
	fmt.Println("Hosts1 length:", len(host1lines))

	host2lines := strings.Split(string(h2bytes), "\n")
	host2lines = normalize(host2lines)
	fmt.Println("Hosts2 length:", len(host2lines))

        intersection := intersect.Simple(host1lines, host2lines)

        fmt.Println("intersection:", intersection)
        fmt.Println("intersection length:", len(intersection))

}

func normalize(slc []string) []string {
	// scrub and trim each element of the slice
	for i := range slc {
		slc[i] = scrub(slc[i], "127.0.0.1")
		slc[i] = scrub(slc[i], "0.0.0.0")
		slc[i] = strings.TrimSpace(slc[i])
	}
	// remove empty elements
	slc = filter(slc, notempty)
	// remove comments
	slc = filter(slc, notcomment)
	return slc
}

func filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func notempty(s string) bool {
	return len(s) > 0
}

func notcomment(s string) bool {
	return !strings.HasPrefix(s, "#")
}

func scrub(s string, r string) string {
	return strings.ReplaceAll(s, r, "")
}
