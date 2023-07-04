package client

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/engine/grpc/client"
	"github.com/desertbit/grumble"
)

func newClient() *client.Client {
	return &client.Client{
		Addr: addr,
	}
}

func putData(c *grumble.Context) error {
	key := c.Args.String("key")
	value := c.Args.String("value")
	if key == "" || value == "" {
		fmt.Println("key or value is empty")
		return nil
	}
	err := newClient().Put([]byte(key), []byte(value))
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
	value, err := newClient().Get([]byte(key))
	if err != nil {
		fmt.Println("get data error: ", err)
		return err
	}
	fmt.Println(string(value))
	return nil
}

func deleteKey(c *grumble.Context) error {
	key := c.Args.String("key")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}
	err := newClient().Del([]byte(key))
	if err != nil {
		fmt.Println("delete key error: ", err)
		return err
	}
	fmt.Println("delete key success")
	return nil
}

func getKeys(c *grumble.Context) error {
	list, err := newClient().Keys()
	if err != nil {
		fmt.Println("get keys error: ", err)
		return err
	}
	fmt.Println("Total keys: ", len(list))
	for i, bytes := range list {
		fmt.Printf(string(bytes[:]) + "\t")
		if i%8 == 7 {
			fmt.Println()
		}
	}
	if len(list)%8 != 0 {
		fmt.Println()
	}
	return nil
}
