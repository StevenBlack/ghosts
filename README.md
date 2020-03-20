# ghosts

`ghosts` is a utility to evaluate, compare, and format hosts files.

Here is what it does:

* Summarize any hosts file retrieved over HTTP, or local file.
* Summarize two hosts files, and determine their intersection.
* Output the hosts as a plain list of domains, or with IP4 pefix.
* Sort the hosts coherently by domain, TLD, subdomain, subsubdomain, and so on.

## Getting started

### Usage examples

Get help just as you might expect:

```bash
$ ./ghosts -h
Usage of ./ghosts:
  -c string
    	Hosts list to compare. A full URL, or a local file.
  -i string
    	The main list of hosts to analyze, or serve as a basis for comparison. A full URL, or a local file. (default "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts")
  -intersection
    	Return the list of intersection hosts? (default false)
  -ip string
    	Localhost IP address (default "0.0.0.0")
  -o	Return the list of hosts? (default false)
  -p	Return a plain output list of hosts? (default false)
  -s	Sort the hosts? (default false)
  -stats
    	display stats? (default true)
  -tld
    	Return the list of TLD and their tally (default false)
```

## Running the tests

`$ go test` runs the test suite.

## Built With

This utility is written with Go.

## Contributing

## License

MIT.
