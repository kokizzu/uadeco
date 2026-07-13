
# uadeco

User-agent decorator for `http.DefaultClient` and `http.DefaultTransport`.
It replaces the `User-Agent` header with `serviceName-time-rev-goVersion`
so outbound HTTP calls can be traced back to a service and build.

## usage:

```go
package main

import "github.com/kokizzu/uadeco"

func init() {
	uadeco.InitServiceName("MyService1")
}
```

## Maintenance checklist

- [x] Go runtime updated to 1.26.5 and module metadata added.
- [x] Removed unnecessary third-party test assertions; the module is dependency-free.
- [x] `make test` uses `httptest` instead of a fixed local port or external network.
- [x] `make verify-dependency-security` and `make vulncheck` check dependency security.
