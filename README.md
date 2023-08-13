## GenPig
**GenPig** is a configuration management package based on generate.

> This package is currently in a very early stage. There may be unexpected bugs or behavior.
 
## Install
``` bash
go install github.com/ironpark/genpig/cmd/genpig@latest
```

## Design
### Make useless functions useful.
Some functions, like the `genpig.SetDefault` function, are empty functions that don't have any implementation inside them. However, they are analyzed by the AST package and referenced when generating code.

```go
func init() {
    genpig.SetConfigPaths("$HOME", ".", "./config")
    genpig.SetConfigNames("myconfig")
}
```
This is an intentional design decision. You can do the same thing with comments or struct tags, but this approach allows the programmer to autocomplete the code and get help from the compiler.

## Example
```go
package conf

import (
	"github.com/ironpark/genpig"
)

func init() {
	genpig.SetConfigPaths("$HOME", ".", "./config")
	genpig.SetConfigNames("myconfig")
}

//go:generate genpig -struct Config
type Config struct {
	Database struct {
		Host     string `env:"DB_SERVER_IP" json:"ip"`
		Port     int    `env:"DB_PORT" json:"port"`
		User     string `env:"DB_USER" json:"user"`
		Password string `env:"DB_PW" json:"password"`
	} `json:"database"`
	Server struct {
		Host string `json:"ip" env:"SERVER_IP"`
		Port int    `json:"port" env:"PORT"`
	} `json:"server"`
}
```