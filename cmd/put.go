package cmd

import (
	"fmt"
	"github.com/desertbit/grumble"
)

func putData(c *grumble.Context) error {
	key := c.Args.String("key")
	value := c.Args.String("value")
	if key == "" || value == "" {
		fmt.Println("key or value is empty")
		return nil
	}
	err := db.Put([]byte(key), []byte(value))
	if err != nil {
		fmt.Println("put data error: ", err)
		return err
	}
	fmt.Println("put data success")
	return nil
}
