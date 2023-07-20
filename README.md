# Rocky
Rocky is a Go interface tool for recursively crawling websites and extracting URLs and walkthrough and extract [query,endpoints].

## Installation

### Go
```
go install github.com/r3m0t3nu11/Rocky@master
```

## Usage
```
Usage: ./Rocky [options] (-url <target>|-file targets.txt)

Usage of rocky:
  -child-process string
    	The command to run as a child process
  -domain string
    	The domain to match (use 'any' for any subdomain)
  -file string
    	The path to a file containing URLs to crawl
  -follow-redirect
    	Follow redirects when making requests (default true)
  -grab string
    	The type of data to grab (options: query, endpoints) (default "query")
  -headless
    	Run the browser in headless mode
  -output string
    	The path to the output file (default "output.txt")
  -url string
    	The URL to crawl

```



## Example

Target `*.example.com`
```
➜ rocky -url https://example.com
``` 



Target all urls in given file
```
➜ rocky -file urls.txt
```

Target `*.example.com` and `*.google.com` (if found)
```
➜ rocky -url https://example.com -domain google.com
```

Target all domains that contain `example`
```
➜ rocky -url https://example.com -domain example
```

Target `*.example.com` but grab any domain `*.*.ltd`  
```
➜ rocky -url https://example.com -domain any
```
Target `example.com` use headless mode 
```
➜ rocky -url https://example.com -headless
```

Target `example.com` follow redirect 
```
➜ rocky -url https://example.com -follow-redirect
```
