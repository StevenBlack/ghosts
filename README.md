# ghosts

A utility, written in Go, to evaluate, compare, and format hosts files.

As chief curator of [StevenBlack/hosts](https://github.com/StevenBlack/hosts), I often need tools to evaluate and compare hosts files.  This repo provides tooling to help manage this.

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

## Built With

## Contributing

## License
