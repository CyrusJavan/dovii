package main

import (
	"fmt"
	"log"

	"github.com/CyrusJavan/dovii"
	"github.com/gin-gonic/gin"
)

type env struct {
	store *dovii.KVStore
}

func main() {
	store, err := dovii.NewKVStore(dovii.BitcaskEngine)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	env := &env{store: store}

	r.GET("/:key", env.getHandler)
	r.POST("/:key/:value", env.setHandler)

	r.Run("127.0.0.1:7070")
}

func (e *env) setHandler(c *gin.Context) {
	key := c.Param("key")
	value := c.Param("value")
	err := (*e.store).Set(key, value)
	if err != nil {
		c.Status(502)
		return
	}
	c.Status(200)
}

func (e *env) getHandler(c *gin.Context) {
	key := c.Param("key")
	value, err := (*e.store).Get(key)
	if err != nil {
		c.JSON(404, gin.H{
			"error": fmt.Sprint(err),
		})
		return
	}
	c.JSON(200, gin.H{
		"value": value,
	})
}
