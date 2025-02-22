package main

import (
	"ContentSystem/internal/api"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	api.CmsRouters(r)
	err := r.Run()
	if err != nil {
		fmt.Printf("r run err: %v")
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
