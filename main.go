package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/juliangruber/go-intersect"
)

type Hosts struct {
	raw     []byte
	url     string
	File    string
	domains []string
}

func (h Hosts) normalize(slc []string) []string {
	for i := range slc {
		slc[i] = h.scrub(slc[i], "127.0.0.1")
		slc[i] = h.scrub(slc[i], "0.0.0.0")
		slc[i] = strings.TrimSpace(slc[i])
	}
	// remove empty elements
	slc = h.filter(slc, h.notempty)
	// remove comments
	slc = h.filter(slc, h.notcomment)

	return slc
}

func (h *Hosts) loadfile(file string) {
	h.File = file
	bytes, err := ioutil.ReadFile(file)
	h.checkerror(err)
	h.raw = bytes
	tempdomains := strings.Split(string(bytes), "\n")
	h.domains = h.normalize(tempdomains)
}

func (h Hosts) length() int {
	return len(h.domains)
}

func (h Hosts) filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func (h Hosts) notempty(s string) bool {
	return len(s) > 0
}

func (h Hosts) notcomment(s string) bool {
	return !strings.HasPrefix(s, "#")
}

func (h Hosts) scrub(s string, r string) string {
	return strings.ReplaceAll(s, r, "")
}

func (h Hosts) replace(s string, r string, n string) string {
	return strings.ReplaceAll(s, r, n)
}

func (h Hosts) checkerror(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	hf1 := Hosts{}
	hf1.loadfile("/Users/Steve/Dropbox/dev/hosts/data/StevenBlack/hosts")
	hf2 := Hosts{}
	hf2.loadfile("/Users/Steve/Dropbox/dev/hosts/data/yoyo.org/hosts")

	intersection := intersect.Simple(hf1.domains, hf2.domains)

	fmt.Println("intersection:", intersection)
	fmt.Println("intersection length:", len(intersection))
}
