
# uadeco

user-agent decorator, injector for `http.DefaultClient` and `http.
DefaultTransport`, it would replace the `User-Agent` header with
`serviceName-time-rev-goVersion`

## usage:

```go
package main

import "github.com/kokizzu/uadeco"

func init() {
	uadeco.InitServiceName("MyService1")
}
```