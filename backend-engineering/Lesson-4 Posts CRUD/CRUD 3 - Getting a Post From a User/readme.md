Installation
Via go install (Recommended)
With go 1.25 or higher:

go install github.com/air-verse/air@latest
Via install.sh
# binary will be $(go env GOPATH)/bin/air
curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# or install it into ./bin/
curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s

air -v
# You create with
air init 


✅ **Short note for cross-platform setup (Mac & Windows):**

In your `air.toml`:

```toml
[build]
  cmd = "go build -o ./bin/main ./cmd/api"
  bin = "./bin/main"
```

Then:

* On **Mac/Linux**, run: `air` → it executes `./bin/main`
* On **Windows**, run:

  ```bash
  go build -o ./bin/main.exe ./cmd/api
  ./bin/main.exe
  ```

Or update `air.toml` to use `.exe` only on Windows.

💡 Tip: Windows needs `.exe`; Mac/Linux doesn’t.
