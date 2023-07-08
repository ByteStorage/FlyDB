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

func putData(c *grumble.Context) error {
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

func getData(c *grumble.Context) error {
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

func deleteKey(c *grumble.Context) error {
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
