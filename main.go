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

	"github.com/dustin/go-humanize"
	"github.com/juliangruber/go-intersect"
)

// Expose the command line flags we support
var inputHosts, compareHosts, ipLocalhost string
var alphaSort, output, plainOutput, stats, intersectionList, tld bool

type TLDtally struct {
	tld   string
	tally int
}

// A Hosts struct holds all the facets of a collection of hosts.
type Hosts struct {
	Raw        []byte
	Location   string
	Domains    []string
	TLDs       map[string]int
	TLDtallies []TLDtally
	Duplicates []string
}

// Reset the Hosts structure to an initial, unloaded state.
func (h *Hosts) Reset() bool {
	// zero everything
	h.Raw = []byte{}
	h.Location = ""
	h.Domains = []string{}
	h.TLDs = map[string]int{}
	h.TLDtallies = []TLDtally{}
	h.Duplicates = []string{}

	return true
}

// summarize the hosts
func (h *Hosts) Summary(prefix string) string {
	var summary []string
	sepLen := 80

	summary = append(summary, strings.Repeat("-", sepLen))
	summary = append(summary, prefix+" summary:")
	summary = append(summary, strings.Repeat("-", sepLen))
	summary = append(summary, "Location: "+h.Location)
	summary = append(summary, "Domains: "+humanize.Comma(int64(len(h.Domains))))
	summary = append(summary, "Bytes: "+humanize.Bytes(uint64(int64(len(h.Raw)))))
	if tld {
		var s []string
		for _, t := range h.TLDtallies {
			s = append(s, t.tld + ": " + humanize.Comma(int64(t.tally)))
		}
		summary = append(summary, "TLD tally:\n   "+strings.Join(s, "\n   "))
	}

	summary = append(summary, strings.Repeat("-", sepLen))

	return strings.Join(summary[:], "\n")
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
	slc = h.filter(slc, h.notEmpty)

	// Step: remove line if it doesn't begin with an IP address
	var ipSlice []string
	for i := range slc {
		words := strings.Fields(slc[i])
		if net.ParseIP(words[0]) == nil {
			continue
		}
		// removing the ip address
		ipSlice = append(ipSlice, strings.Join(words[1:], " "))
	}
	slc = ipSlice

	// we could bail at this juncture
	if len(slc) == 0 {
		return slc
	}

	// Step: split multi-host lines
	var outSlice []string
	for i := range slc {
		newSlice := strings.Split(slc[i], " ")
		outSlice = append(outSlice, newSlice...)
	}
	slc = outSlice

	// regular string sort for deduplication
	sort.Sort(sort.StringSlice(slc))

	// deduplicate
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

	// tally TLDs
	h.TLDs = make(map[string]int)
	h.TLDtallies = []TLDtally{}
	m := map[string]int{}
	n := map[int][]string{}
	for i := range slc {
		ss := strings.Split(slc[i], ".")
		if len(ss) > 1 {
			s := ss[len(ss)-1]
			_, ok := m[s]
			if ok {
				m[s] = m[s] + 1
			} else {
				m[s] = 1
			}
		}
	}
	var a []int
	for k, v := range m {
		n[v] = append(n[v], k)
	}
	for k := range n {
		a = append(a, k)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(a)))

	for _, k := range a {
		for _, s := range n[k] {
			h.TLDs[s] = k
			h.TLDtallies = append(h.TLDtallies, TLDtally{s, k})
		}
	}

	// custom domain sorting
	if alphaSort {
		sort.Sort(domainSort(slc))
	}

	// Stash our slice of domains.
	h.Domains = slc

	if output {
		prefix := ipLocalhost
		for i := range slc {
			if plainOutput {
				fmt.Println(slc[i])
			} else {
				fmt.Println(prefix, slc[i])
			}
		}
	}
	return slc
}

// Load (generically) a list of hosts into the Hosts struc
func (h *Hosts) Load(location string) int {
	// a wrapper to provide a clean loading interface
	clean := strings.ToLower(location)
	if strings.HasPrefix(clean, "http") {
		return h.loadURL(location)
	}
	return h.Loadfile(location)
}

// Load a file of hosts into the Hosts struc
func (h *Hosts) Loadfile(file string) int {
	// loading hosts from the file system
	h.Reset()
	bytes, err := ioutil.ReadFile(file)
	h.checkError(err)
	h.Location = file
	h.Raw = bytes
	h.process()
	return len(bytes)
}

// Load hosts into the Hosts struc from a URL
func (h *Hosts) loadURL(url string) int {
	// loading hosts from a url
	h.Reset()
	client := http.Client{
		Timeout: time.Duration(5000 * time.Millisecond),
	}
	resp, err := client.Get(url)
	h.checkError(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	h.checkError(err)

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

func (h Hosts) notEmpty(s string) bool {
	return len(s) > 0
}

func (h Hosts) notComment(s string) bool {
	return !strings.HasPrefix(s, "#")
}

func (h Hosts) scrub(s string, r string) string {
	return strings.ReplaceAll(s, r, "")
}

func (h Hosts) replace(s string, r string, n string) string {
	return strings.ReplaceAll(s, r, n)
}

func (h Hosts) checkError(err error) {
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
		out = padRight(cslice[parts-2], length, pad)
		out += padRight(cslice[parts-1], length, pad)
		reverseSlice := reverse(cslice)
		if parts > 2 {
			slc := reverseSlice[2:]
			for i := range slc {
				out += padRight(slc[i], length, " ")
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

func padRight(str string, length int, pad string) string {
	return str + times(pad, length-len(str))
}

func reverse(a []string) []string {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	return a
}

func FlagSet() {
	defaultSource := "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts"
	// -i, --input: The first hosts file to load, serving as a basis for what happens subsequently.  Default is StevenBlack hosts.
	flag.StringVar(&inputHosts, "i", defaultSource, "The main list of hosts to analyze, or serve as a basis for comparison")
	flag.StringVar(&inputHosts, "input", defaultSource, "The main list of hosts to analyze, or serve as a basis for comparison")

	// -c, --compare: The second hosts file to load in order to compare, or merge, with the first hosts file.
	flag.StringVar(&compareHosts, "c", "", "Hosts list to compare")
	flag.StringVar(&compareHosts, "compare", "", "Hosts list to compare")

	flag.BoolVar(&output, "o", false, "Return the list of hosts?")
	flag.BoolVar(&output, "output", false, "Return the list of hosts")

	flag.BoolVar(&intersectionList, "intersection", false, "Return the list of intersection hosts")

	flag.BoolVar(&plainOutput, "p", false, "Return a plain output list of hosts?")
	flag.BoolVar(&plainOutput, "plainOutput", false, "Return a plain output list of hosts")

	flag.BoolVar(&tld, "tld", false, "Return the list of TLD and their tally")

	// these flags are not yet implemented

	flag.StringVar(&ipLocalhost, "ip", "0.0.0.0", "Localhost IP address")
	flag.StringVar(&ipLocalhost, "ipaddress", "0.0.0.0", "Localhost IP address")

	flag.BoolVar(&alphaSort, "s", false, "Sort the hosts?")
	flag.BoolVar(&alphaSort, "sort", false, "Sort the hosts?")

	flag.BoolVar(&stats, "stats", true, "display stats?")

	flag.Parse()
}

func main() {

	FlagSet()

	hf1 := Hosts{}
	hf1.Load(inputHosts)

	if stats && !output {
		fmt.Println(hf1.Summary("Base hosts file"))
	}

	if len(compareHosts) > 0 {
		hf2 := Hosts{}
		hf2.Load(compareHosts)
		if stats && !output {
			fmt.Println(hf2.Summary("Compared hosts file"))
		}

		intersection := intersect.Simple(hf1.Domains, hf2.Domains)

		if intersectionList {
			fmt.Println("intersection:", intersection)
		}
		fmt.Println("Intersection:", humanize.Comma(int64(len(intersection))), "domains")
	}
}
