package client

import (
	"fmt"
	"github.com/desertbit/grumble"
)

func hashHSetData(c *grumble.Context) error {
	key := c.Args.String("key")
	field := c.Args.String("field")
	value := c.Args.String("value")
	if key == "" || value == "" {
		fmt.Println("key or value is empty")
		return nil
	}
	err := newClient().HSet(key, field, value)
	if err != nil {
		fmt.Println("put data error: ", err)
		return err
	}
	fmt.Println("put data success")
	return nil
}

func hashHGetData(c *grumble.Context) error {
	key := c.Args.String("key")
	field := c.Args.String("field")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}
	value, err := newClient().HGet(key, field)
	if err != nil {
		fmt.Println("get data error: ", err)
		return err
	}
	fmt.Println(value)
	return nil
}

func hashHDelKey(c *grumble.Context) error {
	key := c.Args.String("key")
	field := c.Args.String("field")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}
	err := newClient().HDel(key, field)
	if err != nil {
		fmt.Println("delete key error: ", err)
		return err
	}
	fmt.Println("delete key success")
	return nil
}
