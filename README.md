# Go version

 go1.20.4 

# Rocky Beta
Rocky is a Go interface tool for recursively crawling websites  extracting URLs from website and walkthrough from previuos urls agin to get as much as possible [query,endpoints].

![Screenshot from 2023-07-20 21-07-37](https://github.com/r3m0t3nu11/Rocky/assets/26588044/0b1db9a8-fb03-4052-ad5d-de65c38a2881)
![Screenshot from 2023-07-20 21-09-09](https://github.com/r3m0t3nu11/Rocky/assets/26588044/2db89224-27f0-48ed-a3d3-63ff6f23852b)
![Screenshot from 2023-07-20 21-10-02](https://github.com/r3m0t3nu11/Rocky/assets/26588044/c0dc9d49-e72c-4e32-8c88-acc115ec36a1)



### Go
```
git clone https://github.com/r3m0t3nu11/Rocky.git&&cd Rocky;go build rocky.go
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
