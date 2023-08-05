package main

import (
		"bufio"
        "context"
		"flag"
		"fmt"
		"os"
        "net/http"
        "net/url"
		"strings"
		"time"

        "github.com/PuerkitoBio/goquery"
        "github.com/chromedp/chromedp"
		"github.com/fatih/color"
)

var (
	urlink         string
	list           string
	domain         string
	headless       bool
	followRedirect bool
	grab           string
	output         string
)

func Banner() {
    name := `Rocky v Beta` 
    author := `[By r3m0t3nu11 https://twitter.com/r3m0t3nu11]`

	fmt.Println(color.RedString(`
	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀       ⠀⠀⠀⡀⣀⣀⣠⠤⠴⠶⠶⠒⠒⠒⠒⠒⠲⣶
	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣤⢴⢾⣿⣟⣷⢤⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿
	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⠴⠛⣍⣴⠼⣿⣻⠟⣿⣟⢭⡷⣤⡀⠀⠀⠀⠀⠀⠀⠀⡇
	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣤⠞⠋⣰⡴⢿⣿⣖⠛⠉⠉⠉⠛⢮⠛⢿⣷⣿⣦⠀⠀⠀⠀⠀⢰⠃
	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⠞⠉⢀⣴⣞⠏⣡⣿⣮⣿⣷⣦⠀⠀⠀⠈⣇⠈⣿⡟⢿⣳⡄⠀⠀⠀⡼⠀
	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡴⠋⠁⢀⣾⠿⢿⣷⡋⣗⣴⣿⠿⠋⠻⣷⡅⠀⢰⠃⣰⡿⠁⠀⢷⡟⡆⠀⢀⡇⠀
	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡴⠋⠀⠀⠀⠉⠛⠷⣾⣍⠀⢿⡿⠃⠀⠀⠀⠘⣿⣶⣣⣴⠋⠀⠀⠀⠀⠹⡽⡆⡾⠀⠀
	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣿⣟⣷⣲⢦⣄⡀⠀⠀⠈⠙⢷⣌⠓⠦⠤⠤⠴⠚⢉⣲⠟⠁⠀⠀⠀⠀⠀⠀⢷⡿⠁⠀⠀
	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⣤⠴⠒⠒⠒⢻⡟⠁⠀⠀⠉⠙⠺⣵⣫⡷⣄⠀⠀⠀⠙⢲⣤⠀⠀⢀⣶⠞⠁⠀⠀⠀⠀⠀⠀⠀⠀⣸⠃⠀⠀⠀
	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡴⠞⠉⠀⠀⠀⠀⠀⣰⣛⡲⠶⢤⣀⠀⠀⠀⢀⣭⣿⣽⢷⣄⠀⠀⠀⠈⢿⡼⠛⠁⠀⠀⠀⠀⠀⠀⠀⠀⢀⡼⠁⠀⠀⠀⠀
	⠀⠀⠀⠀⠀⠀⠀⠀⢀⡴⠋⠀⠀⠀⠀⠀⠀⢀⡾⢁⠟⣹⠛⠲⣄⣙⣶⣞⡽⠋⠀⠈⠻⣼⠳⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡞⠁⠀⠀⠀⠀⠀
	⠀⠀⠀⠀⠀⠀⢀⡴⠋⠀⠀⠀⠀⠀⠀⠀⣠⠏⣰⢋⡾⣱⢏⣴⡿⠋⣷⣟⠳⣄⠀⠀⠀⠈⠳⣜⢳⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⠏⠀⠀⠀⠀⠀⠀⠀
	⠀⠀⠀⠀⠀⢠⠞⠁⠀⠀⢀⣠⠤⠴⠒⢲⠟⠒⠦⢬⣘⢁⡼⠋⢠⣞⡽⢉⡦⣌⠳⡄⠀⠀⠀⠘⢧⠝⢦⠀⠀⡀⠀⠀⠀⠀⢀⡼⠃⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⠀⠀⢀⡴⠁⠀⠀⣠⠖⠉⠀⠀⠀⢠⣯⣤⣀⠀⢀⡼⠋⠀⡰⠋⢚⡴⢫⡞⢙⢢⡈⢦⡀⠀⠀⠀⢫⣭⢧⠀⠀⠀⠀⠀⣠⠟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⠀⣰⠟⠀⠀⢀⠾⠁⠀⠀⠀⠀⠀⡏⠀⠀⠈⣹⠏⣀⡴⠚⠙⢶⣈⠕⢉⡴⢁⡔⠙⢦⠱⣄⠀⠀⠀⣷⠞⢇⠀⠀⣠⠞⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⠀⠛⠛⠛⠒⠋⠀⠀⣀⣠⣤⠄⠒⠛⠒⢺⡟⣑⢾⡉⠳⣄⠀⠀⠙⢦⠉⠐⠋⡠⠊⡈⢧⡙⣆⣀⠼⢻⡏⣿⣀⡜⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⠀⠀⠀⠀⠀⠀⠀⢰⠋⠈⠉⣛⡢⢤⡰⢋⡴⠋⠀⠉⢣⡈⠳⣄⠀⠀⠱⡄⠈⡠⠞⣡⢖⣷⠞⠁⠀⢸⣷⡟⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⠀⠀⠀⠀⠀⠀⠀⠸⣤⠴⠚⠉⢠⡞⣱⠿⡅⠀⠀⠀⠀⠙⢆⠈⢳⡀⠀⠙⣄⠐⢚⡵⠛⠁⠀⠀⠀⣸⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⠀⠀⠀⠀⠀⢀⡴⠊⠁⠀⠀⢀⣠⡝⠁⠀⠈⠲⡀⠀⠀⠀⠈⢳⡀⠹⡄⠀⠸⡞⠋⠀⠀⠀⠀⠀⠀⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⠀⠀⠀⠀⣠⣞⣅⡀⠀⡠⠞⠉⠀⠀⠀⠀⠀⠀⠙⣦⠀⠀⠀⠀⣷⠀⠹⡄⠀⡇⠀⠀⠀⠀⠀⠀⢠⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⠀⠀⠀⡼⠟⣥⣤⠝⠛⠛⢻⠂⠀⠀⠀⠀⢱⠀⠀⠀⢳⡄⠀⢰⠋⠀⠀⣇⡤⡗⠀⠀⠀⠀⠀⠀⢸⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⠀⠀⠀⣠⡴⠛⠁⠀⣀⡤⠟⠀⠀⠀⠀⢀⡿⠀⠀⠀⠀⣹⡀⡘⠉⠓⠋⠉⠀⡇⠀⠀⠀⠀⠀⢀⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⠀⢀⡾⠋⠀⠀⠀⠘⣿⠛⠀⠀⠀⠀⢠⡿⠁⠀⠀⠀⣰⠋⢧⡇⠀⠀⠀⠀⣠⠃⠀⠀⠀⠀⢠⡿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⢠⣾⢾⡇⠀⠀⠀⣸⠁⣀⠀⠾⣽⡶⠋⠀⠀⠀⢀⣴⡃⢀⣸⠃⠀⠀⢀⡴⠋⠀⠀⠀⢀⡴⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⠸⣳⠋⠀⠀⠀⠀⣷⢎⡁⢀⠀⢹⣄⡴⣦⢀⡤⠞⠁⠉⠉⠀⠀⣀⡴⠋⠀⠀⠀⣀⠴⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⢰⠇⠀⠀⠀⣀⣠⢴⡿⢁⡾⣇⡼⠿⠗⠛⠁⠀⠀⠀⠀⠀⠀⢾⠉⠀⠀⢀⣤⠞⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⡎⢀⡤⠒⠋⠻⠶⠯⠕⠊⠙⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡞⢀⣠⠖⠋⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⣷⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠟⠉⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠘⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	
	
	⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    `))

	fmt.Println(color.BlueString(name))
	fmt.Println(color.YellowString(author))
}

func ParseArguments() {
	flag.StringVar(&urlink, "url", "", "The URL to crawl")
	flag.StringVar(&list, "file", "", "The path to a file containing URLs to crawl")
	flag.StringVar(&domain, "domain", "", "The domain to match (use 'any' for any subdomain)")
	flag.BoolVar(&headless, "headless", false, "Run the browser in headless mode")
	flag.BoolVar(&followRedirect, "follow-redirect", true, "Follow redirects when making requests")
	flag.StringVar(&grab, "grab", "query", "The type of data to grab (options: query, endpoints)")
	flag.StringVar(&output, "output", "output.txt", "The path to the output file")

	flag.Parse()
}

func GenerateURLsQueue(urlink string, list string) []string {
    queue := make([]string, 0)

    if urlink == "" && list == "" {
        fmt.Println("Please provide a URL to crawl with the -url flag or a file with the -file flag")
        return queue
    }

    if urlink != "" {
        queue = append(queue, urlink)
    }

    if list != "" {
        file, err := os.Open(list)
        if err != nil {
            fmt.Printf("Error opening file %s: %s\n", list, err.Error())
            return queue
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            link := strings.TrimSpace(scanner.Text())
            if link != "" {
                queue = append(queue, link)
            }
        }

        if err := scanner.Err(); err != nil {
            fmt.Printf("Error reading file %s: %s\n", list, err.Error())
            return queue
        }
    }

    return queue
}

func CreateOutputFile(output string) (*os.File, error) {
    outputFile, err := os.Create(output)
    if err != nil {
        fmt.Println("Error creating output file:", err.Error())
        return nil, err
    }

    return outputFile, nil
}

func main() {
	Banner()

	ParseArguments()

	queue := GenerateURLsQueue(urlink, list)
	
	outputFile, err := CreateOutputFile(output)
    if err != nil {
        return
    }

    defer outputFile.Close()

    startTime := time.Now()
	
	process := processURL(queue, domain, grab, *outputFile, followRedirect)

	if process != nil {
		fmt.Printf("Failed to process URL %s: %s\n", urlink, process)
		return
	}

    endTime := time.Now()
    elapsedTime := endTime.Sub(startTime)
    fmt.Println("Elapsed time:", elapsedTime)
}

func NewChromeDPContext(options ...chromedp.ExecAllocatorOption) (context.Context, context.CancelFunc, error) {
	options = append(options,
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	ctx, cancel := chromedp.NewContext(allocCtx)
	return ctx, cancel, nil
}

func processURL(queue []string, domain string, grab string, outputFile os.File, followRedirect bool) error {
    visited := make(map[string]bool)

    ctx, cancelCtx, err := NewChromeDPContext()
	if err != nil {
		fmt.Println("Failed to create chrome driver context:", err)
	}

	defer cancelCtx()

    for len(queue) > 0 {

        urlink := queue[0]
        queue = queue[1:]

        if visited[urlink] {
            continue
        }

        visited[urlink] = true

        if !followRedirect {
			req, err := http.NewRequest("GET", urlink, nil)
			if err != nil {
				fmt.Println("Failed to create HTTP request:", err)
			}

            client := &http.Client{
                CheckRedirect: func(req *http.Request, via []*http.Request) error {
                    return http.ErrUseLastResponse
                },
            }

            res, err := client.Do(req)
			defer res.Body.Close()

			if err != nil {
				return err
			}

            if res.StatusCode == http.StatusFound || res.StatusCode == http.StatusTemporaryRedirect {
                location := res.Header.Get("Location")
                if location == "" {
                    fmt.Println("Redirect location not found")
                    continue
                }

                parsedUrl, err := url.Parse(location)
                if err != nil {
                    fmt.Println("Failed to parse redirect location:", err)
                    continue
                }

                newURL := req.URL.ResolveReference(parsedUrl).String()

                if !visited[newURL] {
                    queue = append(queue, newURL)
                }

                continue
            }
        }

        var response string
        err = chromedp.Run(
            ctx,
            chromedp.Navigate(urlink),
            chromedp.Evaluate(`document.documentElement.innerHTML`, &response),
        )

        if err != nil {
            fmt.Println("Failed to navigate:", err)
            continue
        }

        doc, err := goquery.NewDocumentFromReader(strings.NewReader(response))
        if err != nil {
            fmt.Printf("Failed to parse HTML for URL %s: %s\n", urlink, err)
            continue
        }

		
        doc.Find("*").Each(func(i int, s *goquery.Selection) {
            var data string
            attr, exists := s.Attr("href")

			if !visited[attr] {
				queue = append(queue, attr)
			}

            if exists && (strings.Contains(attr, "?") || strings.HasPrefix(attr, "/")) {
                var parsedUrl, err = url.Parse(attr)

                if err != nil {
                    return
                }

                if domain == "any" || parsedUrl.Host == domain {
                    switch grab {
                    case "query":
                        data = parsedUrl.RawQuery
                    case "endpoints":
                        data = parsedUrl.Path
                    default:
                        fmt.Println("Invalid value for -grab flag. Options: query, endpoints")
                        return
                    }
                }
            }

            if data != "" {
				fmt.Println(color.GreenString(fmt.Sprintf("\nData found: %s", data)))
				fmt.Println(color.BlueString(fmt.Sprintf("Found on page: %s", urlink)))

				outputFile.WriteString(fmt.Sprintf("Data found: %s\n", data))
				outputFile.WriteString(fmt.Sprintf("Found on page: %s\n\n", urlink))

				processURL(queue, domain, grab, outputFile, followRedirect)
            }
        })
    }

    return nil
}
