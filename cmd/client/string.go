package client

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/engine/grpc/client"
	"github.com/desertbit/grumble"
)

var Addr string

func newClient() *client.Client {
	return &client.Client{
		Addr: Addr,
	}
}

func stringPutData(c *grumble.Context) error {
	key := c.Args.String("key")
	value := c.Args.String("value")
	if key == "" || value == "" {
		fmt.Println("key or value is empty")
		return nil
	}
	err := newClient().Put(key, value)
	if err != nil {
		fmt.Println("put data error: ", err)
		return err
	}
	fmt.Println("put data success")
	return nil
}

func stringGetData(c *grumble.Context) error {
	key := c.Args.String("key")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}
	value, err := newClient().Get(key)
	if err != nil {
		fmt.Println("get data error: ", err)
		return err
	}
	fmt.Println(value)
	return nil
}

func stringDeleteKey(c *grumble.Context) error {
	key := c.Args.String("key")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}
	err := newClient().Del(key)
	if err != nil {
		fmt.Println("delete key error: ", err)
		return err
	}
	fmt.Println("delete key success")
	return nil
}

// GetSet sets the value of a key and returns its old value
func stringGetSet(c *grumble.Context) error {
	key := c.Args.String("key")
	value := c.Args.String("value")
	if key == "" || value == "" {
		fmt.Println("key or value is empty")
		return nil
	}
	oldvalue, err := newClient().Get(key)
	if err != nil {
		fmt.Println("get data error: ", err)
		return err
	}
	fmt.Println("oldvalue: ", oldvalue)
	err = newClient().Put(key, value)
	if err != nil {
		fmt.Println("set data error: ", err)
		return err
	}
	fmt.Println("set data success")
	return nil
}

func stringExists(c *grumble.Context) error {
	key := c.Args.String("key")
	_, err := newClient().Get(key)
	if err != nil {
		fmt.Println("key is not exist")
		return err
	}
	fmt.Println("key is exist")
	return nil
}

func stringMGet(c *grumble.Context) error {
	keys := c.Args.StringList("key")
	if len(keys) == 0 {
		fmt.Println("The number of input keys is empty")
		return nil
	}
	for _, key := range keys {
		value, err := newClient().Get(key)
		if err != nil {
			fmt.Println("mget error on key ", key, ": ", err)
			return err
		}
		fmt.Println("key:value -->", key, ":", value)
	}
	return nil
}

func stringMSet(c *grumble.Context) error {
	keys := c.Args.StringList("keyvalue")
	if len(keys)%2 != 0 {
		fmt.Println("Input format error")
		return nil
	}
	for i := 0; i < len(keys); i += 2 {
		key := keys[i]
		value := keys[i+1]
		err := newClient().Put(key, value)
		if err != nil {
			fmt.Println("put data error: ", err)
			return err
		}
	}
	fmt.Println("put data success")
	return nil
}
