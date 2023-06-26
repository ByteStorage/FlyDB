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
	if db == nil {
		fmt.Println("start server first")
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

func getData(c *grumble.Context) error {
	key := c.Args.String("key")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}
	if db == nil {
		fmt.Println("start server first")
		return nil
	}
	value, err := db.Get([]byte(key))
	if err != nil {
		fmt.Println("get data error: ", err)
		return err
	}
	fmt.Println(string(value[:]))
	return nil
}

func deleteDate(c *grumble.Context) error {
	key := c.Args.String("key")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}
	if db == nil {
		fmt.Println("start server first")
		return nil
	}
	err := db.Delete([]byte(key))
	if err != nil {
		fmt.Println("get data error: ", err)
		return err
	}
	fmt.Println("delete data success")
	return nil
}

func getKeys(c *grumble.Context) error {
	if db == nil {
		fmt.Println("start server first")
		return nil
	}
	list := db.GetListKeys()
	for i, bytes := range list {
		fmt.Println(i+1, "- ", string(bytes[:]))
	}
	return nil
}
