
<div align="center">

[![contributions - welcome](https://img.shields.io/badge/contributions-welcome-blue/green)](/CONTRIBUTING.md "Go to contributions doc")
[![GitHub License](https://img.shields.io/github/license/stfsy/go-api-kit.svg)](https://github.com/stfsy/go-api-kit/blob/master/LICENSE)
<br/>
[![Go Report Card](https://goreportcard.com/badge/github.com/stfsy/go-api-kit)](https://goreportcard.com/report/github.com/stfsy/go-api-kit)
[![Go](https://img.shields.io/github/go-mod/go-version/stfsy/go-api-kit
)](https://go.dev/ "Go to golang homepage")
</div>

<br/>

# go-api-kit

Kickstarts your API by providing out-of-the-box implementations for must-have modules and components for a successful API.

# ⭐ Use Case
Never has it been easier and cheaper to build and provide services, web applications to customers. Current cloud providers like Microsoft, Google, Amazon,
provide a high abstraction level and consumption-based pricing, which allow even individuals to compete with well-established service providers.

Although, the barrier for entering the market has been lowered, the demand for security has steadily increased. It is clear why: The internet has not only connected businesses and customers, or families and friends, but also hackers, and criminals. Thus, the need for security applications, infrastructure and service is as high as it has ever been.

This API starter kit provides a solid foundation for SaaS, Client/Server and API products. It provides out-of-the-box mitigations for the [10 OWASP risks for APIs](https://owasp.org/API-Security/editions/2023/en/0x11-t10/):

- [API1:2023](https://owasp.org/API-Security/editions/2023/en/0xa1-broken-object-level-authorization/) - **Broken Object Level Authorization**	APIs tend to expose endpoints that handle object identifiers, creating a wide attack surface of Object Level Access Control issues. Object level authorization checks should be considered in every function that accesses a data source using an ID from the user.
- [API2:2023](https://owasp.org/API-Security/editions/2023/en/0xa2-broken-authentication/) - **Broken Authentication**	Authentication mechanisms are often implemented incorrectly, allowing attackers to compromise authentication tokens or to exploit implementation flaws to assume other user's identities temporarily or permanently. Compromising a system's ability to identify the client/user, compromises API security overall.
- [API3:2023](https://owasp.org/API-Security/editions/2023/en/0xa3-broken-object-property-level-authorization/) - **Broken Object Property Level Authorization**	This category combines API3:2019 Excessive Data Exposure and API6:2019 - Mass Assignment, focusing on the root cause: the lack of or improper authorization validation at the object property level. This leads to information exposure or manipulation by unauthorized parties.
- [API4:2023](https://owasp.org/API-Security/editions/2023/en/0xa4-unrestricted-resource-consumption/) - **Unrestricted Resource Consumption**	Satisfying API requests requires resources such as network bandwidth, CPU, memory, and storage. Other resources such as emails/SMS/phone calls or biometrics validation are made available by service providers via API integrations, and paid for per request. Successful attacks can lead to Denial of Service or an increase of operational costs.
- [API5:2023](https://owasp.org/API-Security/editions/2023/en/0xa5-broken-function-level-authorization/) - **Broken Function Level Authorization**	Complex access control policies with different hierarchies, groups, and roles, and an unclear separation between administrative and regular functions, tend to lead to authorization flaws. By exploiting these issues, attackers can gain access to other users’ resources and/or administrative functions.
- [API6:2023](https://owasp.org/API-Security/editions/2023/en/0xa6-unrestricted-access-to-sensitive-business-flows/) - **Unrestricted Access to Sensitive Business Flows**	APIs vulnerable to this risk expose a business flow - such as buying a ticket, or posting a comment - without compensating for how the functionality could harm the business if used excessively in an automated manner. This doesn't necessarily come from implementation bugs.
- [API7:2023](https://owasp.org/API-Security/editions/2023/en/0xa7-server-side-request-forgery/) - **Server Side Request Forgery**	Server-Side Request Forgery (SSRF) flaws can occur when an API is fetching a remote resource without validating the user-supplied URI. This enables an attacker to coerce the application to send a crafted request to an unexpected destination, even when protected by a firewall or a VPN.
- [API8:2023](https://owasp.org/API-Security/editions/2023/en/0xa8-security-misconfiguration/) - **Security Misconfiguration**	APIs and the systems supporting them typically contain complex configurations, meant to make the APIs more customizable. Software and DevOps engineers can miss these configurations, or don't follow security best practices when it comes to configuration, opening the door for different types of attacks.
- [API9:2023](https://owasp.org/API-Security/editions/2023/en/0xa9-improper-inventory-management/) - **Improper Inventory Management**	APIs tend to expose more endpoints than traditional web applications, making proper and updated documentation highly important. A proper inventory of hosts and deployed API versions also are important to mitigate issues such as deprecated API versions and exposed debug endpoints.
- [API10:2023](https://owasp.org/API-Security/editions/2023/en/0xaa-unsafe-consumption-of-apis/) - **Unsafe Consumption** of APIs	Developers tend to trust data received from third-party APIs more than user input, and so tend to adopt weaker security standards. In order to compromise APIs, attackers go after integrated third-party services instead of trying to compromise the target API directly.

## 📦 Installation

```bash
go get https://github.com/stfsy/go-api-kit
```

## 🚀 Usage
### main.go
```golang
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/stfsy/go-api-kit/server"
)

var s *server.Server

func main() {
	startServerNonBlocking()
	stopServerAfterSignal()
}

func stopServerAfterSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	s.Stop()

	fmt.Println("Graceful shutdown complete.")
}

func startServerNonBlocking() {
	s = server.NewServer(&server.ServerConfig{
		MuxCallback: func(*http.ServeMux) {
			// add your endpoints and middlewares here
		},
		ListenCallback: func() {
			// do sth just after listen was called on the server instance and
			// just before the server starts serving requests
		},
		// port override is optional but can be used if you want to
		// define the port manually. If empty the value of env.PORT is used.
		PortOverride: "8080",
	})
	go func() {
		err := s.Start()
		if err != nil {
			panic(fmt.Errorf("unable to start server %w", err))
		}
	}()
}
```

### Configuration
This module will read the following environment variables.

#### Env Vars
- `API_KIT_ENV`: default=production
- `API_KIT_MAX_BODY_SIZE`: default=10485760 (bytes) = 10 MB
- `API_KIT_READ_TIMEOUT`: default=10 (seconds)
- `API_KIT_WRITE_TIMEOUT`: default=10 (seconds)
- `API_KIT_IDLE_TIMEOUT`: default=620 (seconds)
#### Standard Env Vars
- `PORT`: default=8080

### Middlewares
By default, four cross-cutting middlewares are configured. These build the foundation for a secure API.

- ✅ accessLog: Logs each incoming request to give insights about usage and response times
- ✅ contentType: Validates the incoming content type, if the request method implies a state change e.g. POST
- ✅ securityHeaders: Adds additional security headers to the response to prevent common attacks and protect users and their data
- ✅ upstreamCacheControl: Instructs proxy servers between the client and the API to not cache responses

## 🧪 Running Tests
To run tests, run the following command

```bash
./test.sh
```

## 📄 License
[MIT](https://choosealicense.com/licenses/mit/)
