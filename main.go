package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/juliangruber/go-intersect"
)

type Hosts struct {
	raw        []byte
	location   string
	domains    []string
	duplicates []string
}

func (h *Hosts) reset() bool {
	h.raw = []byte{}
	h.location = ""
	h.domains = []string{}
	h.duplicates = []string{}

	return true
}

func (h *Hosts) process() []string {
	slc := strings.Split(string(h.raw), "\n")

	// Step: basic cleanup
	for i := range slc {
		// remove embedded comments
		slc[i] = strings.Split(slc[i], "#")[0]

		// remove all extra whitespace
		words := strings.Fields(slc[i])
		slc[i] = strings.Join(words, " ")
	}

	slc = h.filter(slc, h.notempty)

	// Step: remove line if it doesn't begin with an IP address
	var ipslc []string
	for i := range slc {
		words := strings.Fields(slc[i])
		if net.ParseIP(words[0]) == nil {
			continue
		}
		// removing the ip address
		ipslc = append(ipslc, strings.Join(words[1:], " "))
	}
	slc = ipslc

	if len(slc) == 0 {
		return slc
	}

	// Step: split multi-host lines
	var outslc []string
	for i := range slc {
		newslc := strings.Split(slc[i], " ")
		outslc = append(outslc, newslc...)
	}
	slc = outslc

	// sort
	sort.Sort(sort.StringSlice(slc))

	//deduplicate
	j := 0
	for i := 1; i < len(slc); i++ {
		if slc[j] == slc[i] {
			h.duplicates = append(h.duplicates, slc[j])
			continue
		}
		j++
		slc[j] = slc[i]
	}
	slc = slc[:j+1]

	h.domains = slc
	return slc
}

func (h *Hosts) loadfile(file string) int {
	h.reset()
	bytes, err := ioutil.ReadFile(file)
	h.checkerror(err)
	h.location = file
	h.raw = bytes
	h.process()
	return len(bytes)
}

func (h *Hosts) loadurl(url string) int {
	h.reset()
	client := http.Client{
		Timeout: time.Duration(5000 * time.Millisecond),
	}
	resp, err := client.Get(url)
	h.checkerror(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	h.checkerror(err)

	h.location = url
	h.raw = body
	h.process()
	return len(body)
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
	hf1.loadurl("https://raw.githubusercontent.com/StevenBlack/hosts/master/data/StevenBlack/hosts")
	hf2 := Hosts{}
	hf2.loadurl("http://winhelp2002.mvps.org/hosts.txt")

	intersection := intersect.Simple(hf1.domains, hf2.domains)

	fmt.Println("intersection:", intersection)
	fmt.Println("intersection length:", len(intersection))
}
