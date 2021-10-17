# ghosts

`ghosts` is a utility to evaluate, compare, and format hosts files.  It's written in [Go](https://golang.org/).

Here is what `ghosts` does:

* Summarize any hosts file retrieved over HTTP, or from a local file.
* Compare two hosts files, and determine their intersection.
* Compare a reference hosts file with a list of hosts presently in your system clipboard.
* List the tally of TLDs in the hosts file.
* Output the hosts as a plain list of domains, or with IP4 pefix.
* Sort the hosts coherently by domain, TLD, subdomain, subsubdomain, and so on.

## Getting started

### Get help just as you might expect

```
$ ./ghosts -h
Usage of ./ghosts:
  -c string
    	Hosts list to compare.
    	A shortcut code, full URL, or a local file.
    	Use the -m option for the main comparison list.
    	Use the -clip option to use what is on the system clipboard.

    	Shortcut codes
    	==============
    	The following shortcut codes can be used to select among preset main lists.

    	Amalgamated list shortcuts:
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

    	Source list shortcuts:
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
    	-c mvps                  // //winhelp2002.mvps.or
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

  -clip
    	The comparison hosts are in the system clipboard
  -d	Include default hosts at the top of file.
  -intersection
    	Return the list of intersection hosts? (default false)
  -ip string
    	Localhost IP address (default "0.0.0.0")
  -m string
    	The main list of hosts to analyze, or serve as a basis for comparison.
    	A shortcut code, a full URL, or a local file.
    	See the -c flag for the list of shortcut codes. (default "base")
  -noheader
    	Remove the file header from output? (default false)
  -o	Return the list of hosts? (default false)
  -p	Return a plain output list of hosts, with no IP address prefix? (default false)
  -s	Sort the hosts? (default false)
  -stats
    	display stats? (default true)
  -tld
    	Return the list of TLD and their tally (default false)
  -unique
    	List the unique domains in the comparison list
  -v	
		Return the current version
```

### Summarize statistics from any hosts file

**If you specify no hosts file**, by default a summary of [StevenBlack/hosts](https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts) is produced.

```
$ ./ghosts
--------------------------------------------------------------------------------
Base hosts file summary:
--------------------------------------------------------------------------------
Location: https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts
Domains: 54,702
Bytes: 1.7 MB
--------------------------------------------------------------------------------
```
**Specify any hosts file to summarize** by using the `-m` option, like this:

```
$ ./ghosts -m https://someonewhocares.org/hosts/zero/hosts
--------------------------------------------------------------------------------
Base hosts file summary:
--------------------------------------------------------------------------------
Location: https://someonewhocares.org/hosts/zero/hosts
Domains: 14,401
Bytes: 417 kB
--------------------------------------------------------------------------------
```
**Additionally produce a top-level-domain (TLD) report** by using the `-tld` option, like this:

```
$ ./ghosts -m https://someonewhocares.org/hosts/zero/hosts -tld
----------------------------------------
Base hosts file summary:
----------------------------------------
Location: https://someonewhocares.org/hosts/zero/hosts
Domains: 16,933
Bytes: 483 kB
TLD tally:  (177 unique TLD)
   com: 10,093
   net: 2,634
   info: 563
   ru: 296
   de: 263
   org: 241
   pl: 186
   nl: 184
   uk: 158

skipping many lines for brevity

   bo: 1
   rw: 1
   guru: 1
   ae: 1
   men: 1
   ga: 1
   watch: 1
   ac: 1
```

**Compare two hosts files, local or remote, and assess their intersection** by specifying `-m <location>` option for the main hosts file and `-c <location>` option for the second comparison file.

Let's compare the **someonewhocares.org** hosts file (15,474 domains) to the one at **mvps.org** (8,730 domains).  The basic report tells us there are 1,354 domains in the interseation of the two.

Here we use **shortcut presets** to specify the two lists to compare, but we could have specified full URLs for either source.

```
$ ./ghosts -m someonewhocares -c mvps
----------------------------------------
Base hosts file summary:
----------------------------------------
Location: https://someonewhocares.org/hosts/zero/hosts
Domains: 15,474
Bytes: 445 kB
----------------------------------------
----------------------------------------
Compared hosts file summary:
----------------------------------------
Location: https://winhelp2002.mvps.org/hosts.txt
Domains: 8,730
Bytes: 335 kB
----------------------------------------
Intersection: 1,354 domains
```

**Compare two hosts files, local or remote, and LIST their intersection** by specifying `-m <location>` option for the main hosts file, `-c <location>` option for the second comparison file, and add the `--intersection` flag to get the detailed list of the intersecting domains.

Let's compare the **someonewhocares.org** hosts file (14,401 domains) to the one at **mvps.org** (10,473 domains).  The basic report shows us all 1,548 domains in the interseation of the two.

Here we use **shortcut presets** to specify the two lists to compare, but we could have specified full URLs for either source.

```
$ ./ghosts -m someonewhocares -c mvps --intersection
--------------------------------------------------------------------------------
Base hosts file summary:
--------------------------------------------------------------------------------
Location: https://someonewhocares.org/hosts/zero/hosts
Domains: 15,474
Bytes: 445 kB
--------------------------------------------------------------------------------
--------------------------------------------------------------------------------
Compared hosts file summary:
--------------------------------------------------------------------------------
Location: http://winhelp2002.mvps.org/hosts.txt
Domains: 8,730
Bytes: 335 kB
--------------------------------------------------------------------------------
intersection: [006.free-counter.co.uk 102.112.2o7.net 102.122.2o7.net 122.2o7.net 192.168.112.2o7.net
1ca.cqcounter.com 1uk.cqcounter.com 1up.us.intellitxt.com 1us.cqcounter.com .... long list ]
Intersection: 1,354 domains
```

**Compare two hosts files, local or remote, and list what's unique in the second file** by specifying `-m <location>` option for the main hosts file, `-c <location>` option for the second comparison file, and add the `--unique` flag to get the list of domains in the comparison file that are not in the main hoss file.


### Output a list of domains in hosts format, or as a plaintext list

To list domains, use the `-o [optional file]` option.  If you provide no file mame, the list goes to `stdout`.

To sort the domains, use the `-s` flag.

#### Hosts format output

The output is in hosts format by default.  The default IP address is `0.0.0.0`, and you can change that with the `-ip` option. for example, `-ip 127.0.0.1` will list the hosts with a `127.0.0.1` prefix.

The hosts output will include the header of the original source file.  This header typically includes information about the author, and sometimes includes a copyright statement.  To not include the header, use the `-noheader` flag.

To include the list of default hosts at the top of the hosts output, use the `d` flag.  By default the list of hosts does not include the loopback hosts at the top of the list.

These are the default hosts that will be included with the `-d` flag.

```
127.0.0.1 localhost
127.0.0.1 localhost.localdomain
127.0.0.1 local
255.255.255.255 broadcasthost
::1 localhost
::1 ip6-localhost
::1 ip6-loopback
fe80::1%lo0 localhost
ff00::0 ip6-localnet
ff00::0 ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
ff02::3 ip6-allhosts
0.0.0.0 0.0.0.0
```

#### Plaintext domains output

To get a plaintext list of domains, use the `-p` flag.



## Running the tests

`$ go test` runs the test suite.
`$ gotest` runs colorized tests.

## Contributing

TBA.

## License

MIT.

## Related repositories

* [StevenBlack/hosts](https://github.com/StevenBlack/hosts) is my amalgamated hosts file, with custom variants, from various curated sources.
* [StevenBlack/rhosts](https://github.com/StevenBlack/rhosts) hosts tools, written in Rust, just getting started on that.
