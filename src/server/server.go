package server

import (
	"fmt"

	"github.com/CyrusJavan/dovii/src/keyvaluestore"
	"github.com/gin-gonic/gin"
)

// StartServer starts our API server
func StartServer() {
	r := gin.Default()
	var db keyvaluestore.KeyValueStore = make(keyvaluestore.BasicMemory)

	r.GET("/:key", func(c *gin.Context) {
		key := c.Param("key")
		value, err := db.Get(key)
		if err != nil {
			c.JSON(404, gin.H{
				"error": fmt.Sprint(err),
			})
			return
		}
		c.JSON(200, gin.H{
			"value": value,
		})
	})

	r.POST("/:key/:value", func(c *gin.Context) {
		key := c.Param("key")
		value := c.Param("value")
		err := db.Set(key, value)
		if err != nil {
			c.Status(502)
			return
		}
		c.Status(200)
	})

	r.Run("127.0.0.1:7070")
}
