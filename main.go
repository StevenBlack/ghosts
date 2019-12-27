package main

import (
	"flag"
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

// Expose the command line flags we support
var inputhosts, comparehosts, ipLocalhost string
var dedup, alphasort, output, plain bool

// A Hosts struc holds all the facets of a collection of hosts.
type Hosts struct {
	Raw        []byte
	Location   string
	Domains    []string
	Duplicates []string
}

// Reset the Hosts struc to an initial, unloaded state.
func (h *Hosts) Reset() bool {
	// zero everything
	h.Raw = []byte{}
	h.Location = ""
	h.Domains = []string{}
	h.Duplicates = []string{}

	return true
}

func (h *Hosts) process() []string {
	// make a slice with the lines from the Raw domains
	slc := strings.Split(string(h.Raw), "\n")

	// Step: basic cleanup
	for i := range slc {
		// remove embedded comments
		slc[i] = strings.Split(slc[i], "#")[0]

		// remove all extra whitespace
		words := strings.Fields(slc[i])
		slc[i] = strings.Join(words, " ")

		// lowercase everything
		slc[i] = strings.ToLower(slc[i])
	}

	// discard blank lines
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

	// we could bail at this juncture
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

	// regular string sort for deduplication
	sort.Sort(sort.StringSlice(slc))

	// deduplicate
        if dedup {
                j := 0
                for i := 1; i < len(slc); i++ {
                        if slc[j] == slc[i] {
                                h.Duplicates = append(h.Duplicates, slc[j])
                                continue
                        }
                        j++
                        slc[j] = slc[i]
		}
                slc = slc[:j+1]
	}
	// custom domain sorting
	sort.Sort(domainSort(slc))

	// Stash our slice of domains.
	h.Domains = slc
	return slc
}

// Load (generically) a list of hosts into the Hosts struc
func (h *Hosts) Load(location string) int {
	// a wrapper to provide a clean loading interface
	clean := strings.ToLower(location)
	if strings.HasPrefix(clean, "http") {
		return h.Loadurl(location)
	}
	return h.Loadfile(location)
}

// Load a file of hosts into the Hosts struc
func (h *Hosts) Loadfile(file string) int {
	// loading hosts from the file system
	h.Reset()
	bytes, err := ioutil.ReadFile(file)
	h.checkerror(err)
	h.Location = file
	h.Raw = bytes
	h.process()
	return len(bytes)
}

// Load hosts into the Hosts struc from a URL
func (h *Hosts) Loadurl(url string) int {
	// loading hosts from a url
	h.Reset()
	client := http.Client{
		Timeout: time.Duration(5000 * time.Millisecond),
	}
	resp, err := client.Get(url)
	h.checkerror(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	h.checkerror(err)

	h.Location = url
	h.Raw = body
	h.process()
	return len(body)
}

func (h Hosts) length() int {
	return len(h.Domains)
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

// Structure and functions for custom domain sorting.
type domainSort []string

func (s domainSort) Len() int {
	return len(s)
}
func (s domainSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s domainSort) Less(i, j int) bool {
	hf := Hosts{}
	return hf.Normalize(s[i]) < hf.Normalize(s[j])
}

// Normalize the host string for sorting
func (h *Hosts) Normalize(c string) string {
	pad := " "
	length := 50
	cslice := strings.Split(c, ".")
	parts := len(cslice)
	out := c
	if parts > 1 {
		out = padr(cslice[parts-2], length, pad)
		out += padr(cslice[parts-1], length, pad)
		rslice := reverse(cslice)
		if parts > 2 {
			slc := rslice[2:]
			for i := range slc {
				out += padr(slc[i], length, " ")
			}
		}
	}
	return out
}

func times(str string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(str, n)
}

func padr(str string, length int, pad string) string {
	return str + times(pad, length-len(str))
}

func reverse(a []string) []string {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	return a
}

func main() {

	// -i, --input: The first hosts file to load, serving as a basis for what happens subsequently.  Default is my ad-hoc list.
	flag.StringVar(&inputhosts, "i", "https://raw.githubusercontent.com/StevenBlack/hosts/master/data/StevenBlack/hosts", "The main list of hosts to analyze, or serve as a basis for comparison")
	flag.StringVar(&inputhosts, "input", "https://raw.githubusercontent.com/StevenBlack/hosts/master/data/StevenBlack/hosts", "The main list of hosts to analyze, or serve as a basis for comparison")

	// -c, --compare: The second hosts file to load in order to compare, or merge, with the first hosts file.
	flag.StringVar(&comparehosts, "c", "", "Hosts list to compare")
	flag.StringVar(&comparehosts, "compare", "", "Hosts list to compare")

	flag.BoolVar(&alphasort, "s", false, "Sort the hosts?")
	flag.BoolVar(&alphasort, "sort", false, "Sort the hosts?")

	flag.BoolVar(&dedup, "d", true, "De duplicate hosts?")
	flag.BoolVar(&dedup, "dedupe", true, "De duplicate hosts?")

	flag.Parse()

	hf1 := Hosts{}
	hf1.Load(inputhosts)

	if len(comparehosts) > 0 {
		hf2 := Hosts{}
		hf2.Load(comparehosts)
		intersection := intersect.Simple(hf1.Domains, hf2.Domains)

		fmt.Println("intersection:", intersection)
		fmt.Println("intersection length:", len(intersection))
	}
}
