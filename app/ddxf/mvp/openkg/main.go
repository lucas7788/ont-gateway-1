package main

import "github.com/gin-gonic/gin"

const (
	openkgPort   = "10999"
	publishURI   = "/publish"
	buyAndUseURI = "/buyAndUse"
)

// MVP for openkg
func main() {
	r := gin.Default()
	r.POST(publishURI, publish)
	r.POST(buyAndUseURI, buyAndUse)

	r.Run(":" + openkgPort)
}
