package main

import (
    "bufio"
    "context"
    "flag"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "os"
    "os/exec"
    "strings"
    "time"

    "github.com/PuerkitoBio/goquery"
    "github.com/cheggaaa/pb/v3"
    "github.com/chromedp/chromedp"
    "github.com/fatih/color"
)



func processURL(urlStr string, domain string, grab string, foundUrls map[string]string, outputFile *os.File, followRedirect bool, childProcess string, outputFilePath string) error {
    visited := make(map[string]bool)
    queue := []string{urlStr}

    opts := append(chromedp.DefaultExecAllocatorOptions[:],
        chromedp.Flag("headless", true),
        chromedp.Flag("no-sandbox", true),
        chromedp.Flag("disable-dev-shm-usage", true),
    )

    allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
    defer cancel()

    ctx, cancel := chromedp.NewContext(allocCtx)
    defer cancel()

    for len(queue) > 0 {
        urlStr := queue[0]
        queue = queue[1:]

        if visited[urlStr] {
            continue
        }

        visited[urlStr] = true

        req, err := http.NewRequest("GET", urlStr, nil)
        if err != nil {
            fmt.Println("Failed to create HTTP request:", err)
            return err
        }

        if !followRedirect {
            client := &http.Client{
                CheckRedirect: func(req *http.Request, via []*http.Request) error {
                    return http.ErrUseLastResponse
                },
            }
            resp, err := client.Do(req)
            if err != nil {
                fmt.Println("Failed to make request:", err)
                return err
            }
            if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusTemporaryRedirect {
                location := resp.Header.Get("Location")
                if location == "" {
                    fmt.Println("Redirect location not found")
                    continue
                }
                u, err := url.Parse(location)
                if err != nil {
                    fmt.Println("Failed to parse redirect location:", err)
                    continue
                }
                newURL := req.URL.ResolveReference(u).String()
                if !visited[newURL] {
                    queue = append(queue, newURL)
                }
                continue
            }
        }

        var res string
        err = chromedp.Run(ctx,
            chromedp.Navigate(urlStr),
            chromedp.Evaluate(`document.documentElement.innerHTML`, &res),
        )
        if err != nil {
            fmt.Println("Failed to navigate:", err)
            return err
        }

        doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
        if err != nil {
            fmt.Printf("Failed to parse HTML for URL %s: %s\n", urlStr, err)
            continue
        }

        doc.Find("*").Each(func(i int, s *goquery.Selection) {
            attr, exists := s.Attr("href")
            if !exists {
                attr, exists = s.Attr("src")
            }
            if exists && (strings.Contains(attr, "?") || strings.HasPrefix(attr, "/")) {
                u, err := url.Parse(attr)
                if err != nil {
                    return
                }

                if domain == "any" || u.Host == domain || strings.HasSuffix(u.Host, "."+domain) {
                    var data string
                    switch grab {
                    case "query":
                        data = u.RawQuery
                    case "endpoints":
                        data = u.Path
                    default:
                        fmt.Println("Invalid value for -grab flag. Options: query, endpoints")
                        return
                    }

                    if data != "" {
                        if _, ok := foundUrls[attr]; !ok {
                            foundUrls[attr] = urlStr

                            
                            var colorFunc func(format string, a ...interface{}) string
                            switch grab {
                            case "query":
                                colorFunc = color.GreenString
                            case "endpoints":
                                colorFunc = color.MagentaString
                            default:
                                fmt.Println("Invalid value for -grab flag. Options: query, endpoints")
                                return
                            }

                           
                            fmt.Println(colorFunc(fmt.Sprintf("\nData found: %s", data)))
                            fmt.Println(color.BlueString(fmt.Sprintf("Found on page: %s", urlStr)))

                            outputFile.WriteString(fmt.Sprintf("Data found: %s\n", data))
                            outputFile.WriteString(fmt.Sprintf("Found on page: %s\n\n", urlStr))

                            if childProcess != "" {
                                
                                cmd := exec.Command(childProcess, "-url", attr, "-output", outputFilePath)
                                err := cmd.Run()
                                if err != nil {
                                    log.Println("Error running child process:", err)
                                }
                            } else {
                                
                                processURL(attr, domain, grab, foundUrls, outputFile, followRedirect, "", outputFilePath)
                            }
                        }
                    }
                }
            }
        })
    }

    return nil
}

func main() {

        yellow := color.New(color.FgYellow).SprintFunc()
        red := color.New(color.FgRed).SprintFunc()
        blue := color.New(color.FgBlue).SprintFunc()

    asciiArt := `
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
`
    Pb := `Rocky v Beta ` 
    Bp := `[By r3m0t3nu11 https://twittwer.com/r3m0t3nu11]`

    fmt.Println(red(asciiArt))
    fmt.Println(blue(Pb))
    fmt.Println(yellow(Bp))
    urlStr := flag.String("url", "", "The URL to crawl")
    filePath := flag.String("file", "", "The path to a file containing URLs to crawl")
    domain := flag.String("domain", "", "The domain to match (use 'any' for any subdomain)")
    headless := flag.Bool("headless", false, "Run the browser in headless mode")
    followRedirect := flag.Bool("follow-redirect", true, "Follow redirects when making requests")
    grab := flag.String("grab", "query", "The type of data to grab (options: query, endpoints)")
    childProcess := flag.String("child-process", "", "The command to run as a child process")
    outputFilePath := flag.String("output", "output.txt", "The path to the output file")
    flag.Parse()

    if *urlStr == "" && *filePath == "" {
        fmt.Println("Please provide a URL to crawl with the -url flag or a file with the -file flag")
        return
    }

    visited := make(map[string]bool)
    queue := make([]string, 0)

    if *urlStr != "" {
        queue = append(queue, *urlStr)
    }

    if *filePath != "" {
        file, err := os.Open(*filePath)
        if err != nil {
            fmt.Printf("Error opening file %s: %s\n", *filePath, err.Error())
            return
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            urlStr := strings.TrimSpace(scanner.Text())
            if urlStr != "" {
                queue = append(queue, urlStr)
            }
        }

        if err := scanner.Err(); err != nil {
            fmt.Printf("Error reading file %s: %s\n", *filePath, err.Error())
            return
        }
    }

    foundUrls := make(map[string]string)

    outputFile, err := os.Create(*outputFilePath)
    if err != nil {
        fmt.Println("Error creating output file:", err.Error())
        return
    }
    defer outputFile.Close()

    opts := append(chromedp.DefaultExecAllocatorOptions[:],
        chromedp.Flag("headless", *headless),
        chromedp.Flag("no-sandbox", true),
        chromedp.Flag("disable-dev-shm-usage", true),
    )

    allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
    defer cancel()

    ctx, cancel := chromedp.NewContext(allocCtx)
    defer cancel()

    totalPages := len(queue)
    bar := pb.StartNew(totalPages)
    bar.SetWidth(80)

    startTime := time.Now()

    for len(queue) > 0 {
        urlStr := queue[0]
        queue = queue[1:]

        if visited[urlStr] {
            continue
        }

        visited[urlStr] = true

        err := processURL(urlStr, *domain, *grab, foundUrls, outputFile, *followRedirect, *childProcess, *outputFilePath)
        if err != nil {
            fmt.Printf("Failed to process URL %s: %s\n", urlStr, err)
            continue
        }

        req, err := http.NewRequest("GET", urlStr, nil)
        if err != nil {
            fmt.Println("Failed to create HTTP request:", err)
            continue
        }

        if !*followRedirect {
            client := &http.Client{
                CheckRedirect: func(req *http.Request, via []*http.Request) error {
                    return http.ErrUseLastResponse
                },
            }
            resp, err := client.Do(req)
            if err != nil {
                fmt.Println("Failed to make request:", err)
                continue
            }
            if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusTemporaryRedirect {
                location := resp.Header.Get("Location")
                if location == "" {
                    fmt.Println("Redirect location not found")
                    continue
                }
                u, err := url.Parse(location)
                if err != nil {
                    fmt.Println("Failed to parse redirect location:", err)
                    continue
                }
                newURL := req.URL.ResolveReference(u).String()
                if !visited[newURL] {
                    queue = append(queue, newURL)
                }
                continue
            }
        }

        var res string
        err = chromedp.Run(ctx,
            chromedp.Navigate(urlStr),
            chromedp.Evaluate(`document.documentElement.innerHTML`, &res),
        )
        if err != nil {
            fmt.Println("Failed to navigate:", err)
            continue
        }

        doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
        if err != nil {
            fmt.Printf("Failed to parse HTML for URL %s: %s\n", urlStr, err)
            continue
        }

        doc.Find("*").Each(func(i int, s *goquery.Selection) {
            attr, exists := s.Attr("href")
            if !exists {
                attr, exists = s.Attr("src")
            }
            if exists && (strings.Contains(attr, "?") || strings.HasPrefix(attr, "/")) {
                u, err := url.Parse(attr)
                if err != nil {
                    return
                }

                if *domain == "any" || u.Host == *domain || strings.HasSuffix(u.Host, "."+*domain) {
                    var data string
                    switch *grab {
                    case "query":
                        data = u.RawQuery
                    case "endpoints":
                        data = u.Path
                    default:
                        fmt.Println("Invalid value for -grab flag. Options: query, endpoints")
                        return
                    }

                    if data != "" {
                        if _, ok := foundUrls[attr]; !ok {
                            foundUrls[attr] = urlStr

                            
                            var colorFunc func(format string, a ...interface{}) string
                            switch *grab {
                            case "query":
                                colorFunc = color.GreenString
                            case "endpoints":
                                colorFunc = color.MagentaString
                            default:
                                fmt.Println("Invalid value for -grab flag. Options: query, endpoints")
                                return
                            }

                            
                            fmt.Println(colorFunc(fmt.Sprintf("\nData found: %s", data)))
                            fmt.Println(color.BlueString(fmt.Sprintf("Found on page: %s", urlStr)))

                            outputFile.WriteString(fmt.Sprintf("Data found: %s\n", data))
                            outputFile.WriteString(fmt.Sprintf("Found on page: %s\n\n", urlStr))

                            if *childProcess != "" {
                                
                                cmd := exec.Command(*childProcess, "-url", attr, "-output", *outputFilePath)
                                err := cmd.Run()
                                if err != nil {
                                    log.Println("Error running child process:", err)
                                }
                            } else {
                                
                                err := processURL(attr, *domain, *grab, foundUrls, outputFile, *followRedirect, "", *outputFilePath)
                                if err != nil {
                                    fmt.Printf("Failed to process URL %s: %s\n", attr, err)
                                }
                            }
                        }
                    }
                }
            }
        })

        bar.Increment()
    }

    bar.Finish()

    endTime := time.Now()
    elapsedTime := endTime.Sub(startTime)
    fmt.Println("Elapsed time:", elapsedTime)
}
