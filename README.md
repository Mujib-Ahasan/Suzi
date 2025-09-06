#  Suzi Load Tester
Suzi is a lightweight, flexible, and extensible HTTP load testing tool written in Go.
It’s designed for developers and DevOps engineers who want quick insights into application performance without relying on heavy external tools.

## Features
- ### Multiple Attack Types:
     - Basic Load – send a fixed number of requests at a constant rate.
     - Burst Load – hammer the server with a sudden spike of requests.
     - Ramp-Up Load – gradually increase request rate from low to peak.
     - Random Load - send requests in a random manner.
- ### Latency Percentiles (p50, p90, p95, p99):
  get detailed insights into response time distribution. Useful when need to get over all idea about user experience.
- ### Result Visualization:
     - Generates an interactive HTML report with charts (response time per request).
     - Includes a latency summary table showing percentile breakdown.
- ### Configurable Timeout per request:
  simulate slow responses or protect against server hangs. One of the important feature in day-to-day server testing (basic).
- ### Extensible Design:
  future-ready to add features like email notifications after testing or CI/CD integration.

## Installation
- #### Clone the repo first:
       git clone https://github.com/Mujib-Ahasan/Suzi.git
       cd suzi
- #### Build the binary:
      go build -o suzi.exe  -> for windows
      go build -o suzi -> for Linux/macOS
- #### Run the examples below:
       -------Disable the timeout and plot function just by removing the flags---------
       //Basic attack
      ./suzi -url https://example.com -req 50 -rate 5 -atk basic -timeout 5 -plot true

       //Burst attack
      ./suzi -url https://example.com -req 50 -atk burst -timeout 5 -plot true

       //Ramp-Up attack
      ./suzi -url https://example.com -req 50  -atk rampup -timeout 5 -plot true

       // Random attack
      ./suzi -url https://example.com -req 50 -rate 5 -attack basic -timeout 5 -plot true
## CLI Flags

| Flag        | Default    | Description                              |
| ----------- | ---------- | ---------------------------------------- |
| `-url`      | (required) | Target site or API endpoint to test      |
| `-req`      | 10         | Number of requests to send               |
| `-atk`      | `basic`    | Attack type (`basic`, `burst`, `rampup`, `random`) |
| `-method`   | `GET`      | HTTP method (GET, planning to add POST)  |
| `-rate`     | 1          | Number of requests per second            |
| `-timeout`  | 5          | Request timeout in seconds               |
| `-cpus`     | NumCPU     | Number of CPU cores to use               |
| `-plot`     | `false`    | Generate interactive HTML report         |

## Now The Most Fun One - Why Suzi? 🤔
- Most load testers (like JMeter or Locust) are heavy, complex, and overkill for quick checks.
- Our Suzi Offers:
     - Simple – one binary, easy flags, no config files.
     - Fast – written in Go with goroutines for concurrency.
     - Extensible – easy to add attack strategies and plotting.
     - Developer-Friendly – quick experiments, local testing, and custom benchmark.
## Contributing 
Contributions are always welcome! Feel free to open issues or submit pull requests for new features, enhancements, or bug fixes.
## License
 MIT License – free to use and modify.
