## GenPig
Experimental configuration package for typesafety-obsessed gopher

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
### Use your own package name 
Genpig generates configuration management logic by referencing the package name in an automatically generated location.

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

type Database struct {
	Host     string `env:"DB_SERVER_IP"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PW"`
}

type Server struct {
	Host string `json:"ip" env:"SERVER_IP"`
	Port int    `json:"port" env:"PORT"`
}

//go:generate genpig -struct Config
type Config struct {
	MainDB Database
	Server `json:"server"`
}
```