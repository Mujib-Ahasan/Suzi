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
- ### Email Option:
  Email notifications after testing or CI/CD integration can be sent for better review.

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
| `-email`     | `false`    | Send results via email        |
| `-emailTo`     | `"you@local.test"`    | Receiver email-id       |
| `-emailFrom`     | `"Suzi <noreply@gmail.com>"`    | Sender email-id         |
| `-smtpHost`     | `"localhost"`    | For production use "smtp.gmail.com"     |
| `-smtpPort`     | `1025`    |For Production use 465 (donot use 587)     |
| `-smtp-user`     | `os.Getenv("SMTP_USER")`    | Set to the email address that owns the app password (set env var for this)|
| `-smtp-pass`     | `os.Getenv("SMTP_PASS")`    | App Password generated from Gmail (set env var for this)  |
| `-smtpTLS`     | `false`    | True for production        |
| `-smtp-retries`     | `3`    | Email send retries         |
| `-smtp-timeout`     | `10`    |Email send timeout in seconds      |

## Email Option
- For Local testing use Docker `docker run --name mailhog -d -p 8025:8025 -p 1025:1025 mailhog/mailhog`. After starting the container open localhost at port 8025. </br>
  Run the command `./suzi -url https://example.com -rate 5 -req 50 -email -atk mailAll`

- For production usecase run the command carefully `./suzi -url https://example.com -rate 5 -req 50 -email -atk mailAll -email /` </br>
   `-emailTo "<receiver email-id> -emailFrom "Suzi <email-addr that owns app-password>" /` </br>
   `-smtpHost "smtp.gmail.com" -smtpPort 465  -smtpTLS true` </br>
  After that receiver should get a email like this.
![WhatsApp Image 2025-09-09 at 1 31 07 PM](https://github.com/user-attachments/assets/5e93931f-4f00-44fb-8c12-8be5364f2e58)
![WhatsApp Image 2025-09-09 at 1 32 00 PM](https://github.com/user-attachments/assets/94da2a83-6c7d-4220-9e19-05db4411a15b)
![WhatsApp Image 2025-09-09 at 1 32 28 PM](https://github.com/user-attachments/assets/4697ee51-bf20-4493-b30f-7fe4cc93b0dc)

**Note: In Releases, you would get the binary for linux**


## Now The Most Fun One - Why Suzi? 🤔
- Most load testers (like JMeter or Locust) are heavy, complex, and overkill for quick checks.
- Our Suzi Offers:
     - Simple – one binary, easy flags, no config files.
     - Fast – written in Go with goroutines for concurrency.
     - Extensible – easy to add attack strategies and plotting.
     - Developer-Friendly – quick experiments, local testing, and custom benchmark.
     - Emailing - can enhance development and deployment experience.
## Contributing 
Contributions are always welcome! Feel free to open issues or submit pull requests for new features, enhancements, or bug fixes.
## License
 MIT License – free to use and modify.
