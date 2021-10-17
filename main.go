package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/dustin/go-humanize"
	"github.com/thoas/go-funk"
)

// Expose the command line flags we support
var mainHosts, compareHosts, ipLocalhost string
var addDefaults, alphaSort, output, plainOutput, stats, intersectionList, tld, noheader, sysclipboard, uniquelist, version bool

type TLDtally struct {
	tld   string
	tally int
}

// Update the version # before every release.
const VERSION = "v0.2.3.Beta"

// A Hosts struct holds all the facets of a collection of hosts.
type Hosts struct {
	Raw          []byte
	Location     string
	Header       []string
	Domains      []string
	TLDs         map[string]int
	TLDtallies   []TLDtally
	Duplicates   []string
	Intersection []string
	Unique       []string
}

// Reset the Hosts structure to an initial, unloaded state.
func (h *Hosts) Reset() bool {
	// zero everything
	h.Raw = []byte{}
	h.Location = ""
	h.Header = []string{}
	h.Domains = []string{}
	h.TLDs = map[string]int{}
	h.TLDtallies = []TLDtally{}
	h.Duplicates = []string{}
	h.Intersection = []string{}
	h.Unique = []string{}

	return true
}

// summarize the hosts
func (h *Hosts) Summary(prefix string) string {
	var summary []string
	sepLen := 40

	summary = append(summary, strings.Repeat("-", sepLen))
	summary = append(summary, prefix+" summary:")
	summary = append(summary, strings.Repeat("-", sepLen))
	summary = append(summary, "Location: "+h.Location)
	summary = append(summary, "Domains: "+humanize.Comma(int64(len(h.Domains))))
	summary = append(summary, "Bytes: "+humanize.Bytes(uint64(int64(len(h.Raw)))))
	if tld {
		var s []string
		for _, t := range h.TLDtallies {
			s = append(s, t.tld+": "+humanize.Comma(int64(t.tally)))
		}
		summary = append(summary, "TLD tally:  ("+strconv.Itoa(len(s))+" unique TLD)\n   "+strings.Join(s, "\n   "))
		// summary = append(summary, strings.Join(s, "\n   "))
	}

	summary = append(summary, strings.Repeat("-", sepLen))

	return strings.Join(summary[:], "\n")
}

func (h *Hosts) process() []string {
	// make a slice with the lines from the Raw domains
	slc := strings.Split(string(h.Raw), "\n")

	// Step: preserve the header
	for i := range slc {
		tst := strings.TrimSpace(slc[i])
		if strings.HasPrefix(tst, "#") || len(tst) == 0 {
			h.Header = append(h.Header, slc[i])
		} else {
			break
		}
	}

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

	// Step: discard blank lines
	slc = h.filter(slc, h.notEmpty)

	// step: line match regex for ip address, domain, or host
	// This regex matches domain, or host
	r, _ := regexp.Compile("((^(?:[a-z_0-9](?:[a-z_0-9-]{0,61}[a-z_0-9])?\\.)+[a-z_0-9][a-z_0-9-]{0,61}[a-z_0-9]$)|((^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?))(\\s+((?:[a-z_0-9](?:[a-z_0-9-]{0,61}[a-z_0-9])?\\.)+[a-z_0-9][a-z_0-9-]{0,61}[a-z_0-9]\\s*)+$)))")
	var matchSlice []string
	for i := range slc {
		if r.MatchString(slc[i]) {
			words := strings.Fields(slc[i])
			if net.ParseIP(words[0]) == nil {
				// no IP segment - handle case of multi-host line
				newSlice := strings.Split(slc[i], " ")
				matchSlice = append(matchSlice, newSlice...)
			} else {
				// remove the IP segment
				newSlice := strings.Split(strings.Join(words[1:], " "), " ")
				matchSlice = append(matchSlice, newSlice...)
			}
		}
	}
	slc = matchSlice

	// we could bail at this juncture
	if len(slc) == 0 {
		return slc
	}

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
		// first, the header
		if !noheader {
			for i := range h.Header {
				fmt.Println(h.Header[i])
			}
		}

		// add defaults
		if addDefaults && !plainOutput {
			fmt.Println("127.0.0.1 localhost")
			fmt.Println("127.0.0.1 localhost.localdomain")
			fmt.Println("127.0.0.1 local")
			fmt.Println("255.255.255.255 broadcasthost")
			fmt.Println("::1 localhost")
			fmt.Println("::1 ip6-localhost")
			fmt.Println("::1 ip6-loopback")
			fmt.Println("fe80::1%lo0 localhost")
			fmt.Println("ff00::0 ip6-localnet")
			fmt.Println("ff00::0 ip6-mcastprefix")
			fmt.Println("ff02::1 ip6-allnodes")
			fmt.Println("ff02::2 ip6-allrouters")
			fmt.Println("ff02::3 ip6-allhosts")
			fmt.Println("0.0.0.0 0.0.0.0")
			fmt.Println("")
		}

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
	var client = http.Client{
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

// Load hosts from the clipboard
func (h *Hosts) LoadClipboard(clip string) int {
	// loading hosts from the file system
	h.Reset()
	bytes := []byte(clip)
	h.Location = "clipboard"
	h.Raw = bytes
	h.process()
	return len(bytes)
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
	defaultMainHosts := "base"
	flag.StringVar(&compareHosts, "c", "", `Hosts list to compare.
A shortcut code, full URL, or a local file.
Use the -m option for the main comparison list.
Use the -clip option to use what is on the system clipboard.

Shortcut codes
==============
The following shortcut codes can be used to select among preset main lists.

Amalgamated lists' shortcuts:
-c b or -m base // use Steven Black's base amalgamated list.
-c f    // use alternates/fakenews/hosts
-c fg   // use alternates/fakenews-gambling/hosts
-c fgp  // use alternates/fakenews-gambling-porn/hosts
-c fgps // use alternates/fakenews-gambling-porn-social/hosts
-c fgs  // use alternates/fakenews-gambling-social/hosts
-c fp   // use alternates/fakenews-porn/hosts
-c fps  // use alternates/fakenews-porn-social/hosts
-c fs   // use alternates/fakenews-social/hosts
-c g    // use alternates/gambling/hosts
-c gp   // use alternates/gambling-porn/hosts
-c gps  // use alternates/gambling-porn-social/hosts
-c gs   // use alternates/gambling-social/hosts
-c p    // use alternates/porn/hosts
-c ps   // use alternates/porn-social/hosts
-c s    // use alternates/social/hosts

Source lists' shortcuts:
-c adaway                // adaway.github.io
-c add2o7net             // FadeMind add.2o7Net hosts
-c adddead               // FadeMind add.Dead hosts
-c addrisk               // FadeMind add.Risk hosts
-c addspam               // FadeMind add.Spam hosts
-c adguard               // AdguardTeam cname-trackers
-c baddboyz              // mitchellkrogza Badd-Boyz-Hosts
-c clefspear             // Clefspeare13 pornhosts
-c digitalside           // davidonzo Threat-Intel
-c fakenews              // marktron/fakenews
-c hostsvn               // bigdargon hostsVN
-c kadhosts              // PolishFiltersTeam
-c metamask              // MetaMask eth-phishing hosts
-c mvps                  // winhelp2002.mvps.or
-c orca                  // orca.pet notonmyshift hosts
-c shady                 // hreyasminocha shady hosts
-c sinfonietta-gambling
-c sinfonietta-porn
-c sinfonietta-snuff
-c sinfonietta-social
-c someonewhocares       // Sam Pollock someonewhocares.org
-c stevenblack           // Steven Black ad-hoc list
-c tiuxo-porn
-c tiuxo-social
-c tiuxo                 // tiuxo list.
-c uncheckyads           // FadeMind  UncheckyAds
-c urlhaus               // urlhaus.abuse.ch
-c yoyo                  // Peter Lowe yoyo.org
`)
	flag.BoolVar(&sysclipboard, "clip", false, "The comparison hosts are in the system clipboard")
	flag.BoolVar(&addDefaults, "d", false, "Include default hosts at the top of file.")
	flag.BoolVar(&intersectionList, "intersection", false, "Return the list of intersection hosts? (default false)")
	flag.BoolVar(&uniquelist, "unique", false, "List the unique domains in the comparison list")
	flag.StringVar(&ipLocalhost, "ip", "0.0.0.0", "Localhost IP address")
	flag.StringVar(&mainHosts, "m", defaultMainHosts, `The main list of hosts to analyze, or serve as a basis for comparison.
A shortcut code, a full URL, or a local file.
See the -c flag for the list of shortcut codes.`)
	flag.BoolVar(&noheader, "noheader", false, "Remove the file header from output? (default false)")
	flag.BoolVar(&output, "o", false, "Return the list of hosts? (default false)")
	flag.BoolVar(&plainOutput, "p", false, "Return a plain output list of hosts, with no IP address prefix? (default false)")
	flag.BoolVar(&alphaSort, "s", false, "Sort the hosts? (default false)")
	flag.BoolVar(&stats, "stats", true, "display stats?")
	flag.BoolVar(&tld, "tld", false, "Return the list of TLD and their tally (default false)")
	flag.BoolVar(&version, "v", false, "Return the current version")
	flag.Parse()
}

func main() {

	FlagSet()

	hf1 := Hosts{}
	listShortcuts := map[string]string{
		"b":                    "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts",
		"base":                 "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts",
		"f":                    "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/fakenews/hosts",
		"fg":                   "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/fakenews-gambling/hosts",
		"fgp":                  "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/fakenews-gambling-porn/hosts",
		"fgps":                 "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/fakenews-gambling-porn-social/hosts",
		"fgs":                  "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/fakenews-gambling-social/hosts",
		"fp":                   "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/fakenews-porn/hosts",
		"fps":                  "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/fakenews-porn-social/hosts",
		"fs":                   "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/fakenews-social/hosts",
		"g":                    "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/gambling/hosts",
		"gp":                   "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/gambling-porn/hosts",
		"gps":                  "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/gambling-porn-social/hosts",
		"gs":                   "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/gambling-social/hosts",
		"p":                    "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/porn/hosts",
		"ps":                   "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/porn-social/hosts",
		"s":                    "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/social/hosts",
		"adaway":               "https://raw.githubusercontent.com/AdAway/adaway.github.io/master/hosts.txt",
		"add2o7net":            "https://raw.githubusercontent.com/FadeMind/hosts.extras/master/add.2o7Net/hosts",
		"adddead":              "https://raw.githubusercontent.com/FadeMind/hosts.extras/master/add.Dead/hosts",
		"addrisk":              "https://raw.githubusercontent.com/FadeMind/hosts.extras/master/add.Risk/hosts",
		"addspam":              "https://raw.githubusercontent.com/FadeMind/hosts.extras/master/add.Spam/hosts",
		"adguard":              "https://raw.githubusercontent.com/AdguardTeam/cname-trackers/master/combined_disguised_trackers_justdomains.txt",
		"baddboyz":             "https://raw.githubusercontent.com/mitchellkrogza/Badd-Boyz-Hosts/master/hosts",
		"clefspear":            "https://raw.githubusercontent.com/Clefspeare13/pornhosts/master/0.0.0.0/hosts",
		"digitalside":          "https://raw.githubusercontent.com/davidonzo/Threat-Intel/master/lists/latestdomains.piHole.txt",
		"fakenews":             "https://raw.githubusercontent.com/marktron/fakenews/master/fakenews",
		"hostsvn":              "https://raw.githubusercontent.com/bigdargon/hostsVN/master/option/hosts-VN",
		"kadhosts":             "https://raw.githubusercontent.com/PolishFiltersTeam/KADhosts/master/KADhosts.txt",
		"metamask":             "https://raw.githubusercontent.com/MetaMask/eth-phishing-detect/master/src/hosts.txt",
		"mvps":                 "https://winhelp2002.mvps.org/hosts.txt",
		"orca":                 "https://orca.pet/notonmyshift/hosts.txt",
		"shady":                "https://raw.githubusercontent.com/shreyasminocha/shady-hosts/main/hosts",
		"sinfonietta-gambling": "https://raw.githubusercontent.com/Sinfonietta/hostfiles/master/gambling-hosts",
		"sinfonietta-porn":     "https://raw.githubusercontent.com/Sinfonietta/hostfiles/master/pornography-hosts",
		"sinfonietta-snuff":    "https://raw.githubusercontent.com/Sinfonietta/hostfiles/master/snuff-hosts",
		"sinfonietta-social":   "https://raw.githubusercontent.com/Sinfonietta/hostfiles/master/social-hosts",
		"someonewhocares":      "https://someonewhocares.org/hosts/zero/hosts",
		"stevenblack":          "https://raw.githubusercontent.com/StevenBlack/hosts/master/data/StevenBlack/hosts",
		"tiuxo-porn":           "https://raw.githubusercontent.com/tiuxo/hosts/master/porn",
		"tiuxo-social":         "https://raw.githubusercontent.com/tiuxo/hosts/master/social",
		"tiuxo":                "https://raw.githubusercontent.com/tiuxo/hosts/master/ads",
		"uncheckyads":          "https://raw.githubusercontent.com/FadeMind/hosts.extras/master/UncheckyAds/hosts",
		"urlhaus":              "https://urlhaus.abuse.ch/downloads/hostfile/",
		"yoyo":                 "https://pgl.yoyo.org/adservers/serverlist.php?hostformat=hosts&mimetype=plaintext&useip=0.0.0.0",
	}

	_, shortCode := listShortcuts[mainHosts]
	if shortCode {
		mainHosts = listShortcuts[mainHosts]
	}

	if version {
		fmt.Println("The current version is:", VERSION)
		os.Exit(0)
	}

	hf1.Load(mainHosts)

	if stats && !output {
		fmt.Println(hf1.Summary("Base hosts file"))
	}

	if len(compareHosts) > 0 {
		_, shortCode := listShortcuts[compareHosts]
		if shortCode {
			compareHosts = listShortcuts[compareHosts]
		}

		hf2 := Hosts{}
		hf2.Load(compareHosts)
		if stats && !output {
			fmt.Println(hf2.Summary("Compared hosts file"))
		}

		hf2.Intersection = funk.IntersectString(hf1.Domains, hf2.Domains)
		if intersectionList {
			// for now, unceremoniously dump the intersecting domains.
			fmt.Println("intersection:", hf2.Intersection)
		}
		fmt.Println("Intersection:", humanize.Comma(int64(len(hf2.Intersection))), "domains")

		if uniquelist {
			_, hf2.Unique = funk.DifferenceString(hf2.Intersection, hf2.Domains)
			fmt.Println(strings.Repeat("-", 40))
			fmt.Println("Unique in comparison list — ", humanize.Comma(int64(len(hf2.Unique))), "domains", hf2.Unique)
		}
	} else if sysclipboard {
		hf2 := Hosts{}
		clip, _ := clipboard.ReadAll()
		hf2.LoadClipboard(clip)
		if stats && !output {
			fmt.Println(hf2.Summary("Compared hosts from clipboard"))
		}

		hf2.Intersection = funk.IntersectString(hf1.Domains, hf2.Domains)

		if intersectionList {
			// for now, unceremoniously dump the intersecting domains.
			fmt.Println("intersection:", hf2.Intersection)
		}
		fmt.Println("Intersection:", humanize.Comma(int64(len(hf2.Intersection))), "domains")

		if uniquelist {
			_, hf2.Unique = funk.DifferenceString(hf2.Intersection, hf2.Domains)
			fmt.Println("unique in comparison list:", hf2.Unique)
		}
	}
}
