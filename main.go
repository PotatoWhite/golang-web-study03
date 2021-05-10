package main

import (
	"github.com/potatowhite/web/study03/myapp"
	"net/http"
)

func main() {
	_ = http.ListenAndServe(":3000", myapp.NewHandler())
}
