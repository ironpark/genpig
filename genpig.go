package genpig

import "reflect"

type Format string

func SetConfigPaths(paths ...string) {

}

func SetConfigNames(name ...string) {

}

func SetDefault(defaultValue any) {
	if reflect.TypeOf(defaultValue).Kind() != reflect.Struct {
		panic("Default value of config can only be set to struct")
	}
}
