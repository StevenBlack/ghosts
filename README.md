# ghosts

A utility, written in Go, to evaluate, compare, and format hosts files.

As chief curator of [StevenBlack/hosts](https://github.com/StevenBlack/hosts), I often need tools to evaluate and compare hosts files.  This repo provides tooling to help manage this.

## Getting started

### Usage examples

Get help just as you might expect:

```bash
$ ./ghosts --help

Usage of ./ghosts:
  -c string
    	Hosts list to compare
  -compare string
    	Hosts list to compare
  -i string
    	The main list of hosts to analyze, or serve as a basis for comparison (default "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts")
  -input string
    	The main list of hosts to analyze, or serve as a basis for comparison (default "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts")
  -intersection
    	Return the list of intersection hosts
  -ip string
    	Localhost IP address (default "0.0.0.0")
  -ipaddress string
    	Localhost IP address (default "0.0.0.0")
  -o	Return the list of hosts?
  -output
    	Return the list of hosts
  -p	Return a plain output list of hosts?
  -plainOutput
    	Return a plain output list of hosts
  -s	Sort the hosts?
  -sort
    	Sort the hosts?
  -stats
    	display stats? (default true)
  -tld
      Return the list of TLD and their tally
```

## Running the tests

## Built With

## Contributing

## License
