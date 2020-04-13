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

```bash
$ ./ghosts -h
Usage of ./ghosts:
  -c string
    	Hosts list to compare. A full URL, or a local file.
    	Use the -m option for the main comparison list.
    	Use the --clip option to use what is on the system clipboard.
  -clip
    	The comparison hosts are in the system clipboard
  -intersection
    	Return the list of intersection hosts? (default false)
  -ip string
    	Localhost IP address (default "0.0.0.0")
  -m string
    	The main list of hosts to analyze, or serve as a basis for comparison. A full URL, or a local file.
    	 (default "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts")
  -noheader
    	Remove the file header from output? (default false)
  -o	Return the list of hosts? (default false)
  -p	Return a plain output list of hosts, with no IP address prefix? (default false)
  -s	Sort the hosts? (default false)
  -stats
    	display stats? (default true)
  -tld
    	Return the list of TLD and their tally (default false)
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
--------------------------------------------------------------------------------
Base hosts file summary:
--------------------------------------------------------------------------------
Location: https://someonewhocares.org/hosts/zero/hosts
Domains: 14,401
Bytes: 417 kB
TLD tally:
   com: 8,683
   net: 2,231
   info: 479
   ru: 291
   de: 232
   pl: 181
   org: 174
   fr: 147
   nl: 140
   at: 133
   uk: 117

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

Let's compare the **someonewhocares.org** hosts file (14,401 domains) to the one at **mvps.org** (10,473 domains).  The basic report tells us there are 1,548 domains in the interseation of the two.

```
$ ./ghosts -m https://someonewhocares.org/hosts/zero/hosts -c http://winhelp2002.mvps.org/hosts.txt
--------------------------------------------------------------------------------
Base hosts file summary:
--------------------------------------------------------------------------------
Location: https://someonewhocares.org/hosts/zero/hosts
Domains: 14,401
Bytes: 417 kB
--------------------------------------------------------------------------------
--------------------------------------------------------------------------------
Compared hosts file summary:
--------------------------------------------------------------------------------
Location: http://winhelp2002.mvps.org/hosts.txt
Domains: 10,473
Bytes: 392 kB
--------------------------------------------------------------------------------
Intersection: 1,548 domains
```

**Compare two hosts files, local or remote, and LIST their intersection** by specifying `-m <location>` option for the main hosts file, `-c <location>` option for the second comparison file, and add the `--intersection` flag to get the detailed list of the intersecting domains..

Let's compare the **someonewhocares.org** hosts file (14,401 domains) to the one at **mvps.org** (10,473 domains).  The basic report shows us all 1,548 domains in the interseation of the two.

```
$ ./ghosts -m https://someonewhocares.org/hosts/zero/hosts -c http://winhelp2002.mvps.org/hosts.txt --intersection
--------------------------------------------------------------------------------
Base hosts file summary:
--------------------------------------------------------------------------------
Location: https://someonewhocares.org/hosts/zero/hosts
Domains: 14,401
Bytes: 417 kB
--------------------------------------------------------------------------------
--------------------------------------------------------------------------------
Compared hosts file summary:
--------------------------------------------------------------------------------
Location: http://winhelp2002.mvps.org/hosts.txt
Domains: 10,473
Bytes: 392 kB
--------------------------------------------------------------------------------
intersection: [006.free-counter.co.uk 102.112.2o7.net 102.122.2o7.net 122.2o7.net 192.168.112.2o7.net 
1ca.cqcounter.com 1uk.cqcounter.com 1up.us.intellitxt.com 1us.cqcounter.com .... long list ]
Intersection: 1,548 domains
```


## Running the tests

`$ go test` runs the test suite.

## Contributing

TBA.

## License

MIT.

## Related repositories

* [StevenBlack/hosts](https://github.com/StevenBlack/hosts) is my amalgamated hosts file, with custom variants, from various curated sources.
* [StevenBlack/rhosts](https://github.com/StevenBlack/ghosts) hosts tools, written in Rust, just getting started on that.
