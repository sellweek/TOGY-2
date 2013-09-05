package controllers

import (
	"fmt"
	"github.com/sellweek/TOGY-2/util"
)

func Home(c util.Context) error {
	fmt.Fprintf(c.W, "Hello, world!")
	return nil
}
