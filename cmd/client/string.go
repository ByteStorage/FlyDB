package client

import (
	"errors"
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

func stringGetType(c *grumble.Context) error {
	key := c.Args.String("key")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}
	valueType, err := newClient().Type(key)
	if err != nil {
		fmt.Println("get type error: ", err)
		return err
	}
	fmt.Println("Type:", valueType)
	return nil
}

func stringStrLen(c *grumble.Context) error {
	key := c.Args.String("key")
	strLen, err := newClient().StrLen(key)
	if err != nil {
		return nil
	}
	fmt.Println(strLen)
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
	oldValue, err := newClient().GetSet(key, value)
	if err != nil {
		return err
	}
	fmt.Println(oldValue)
	return nil
}

func stringAppend(c *grumble.Context) error {
	key := c.Args.String("key")
	value := c.Args.String("value")
	err := newClient().Append(key, value)
	if err != nil {
		fmt.Println("Append failed")
		return err
	}
	fmt.Println("Append is successful")
	return nil
}

func stringIncr(c *grumble.Context) error {
	key := c.Args.String("key")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}
	err := newClient().Incr(key)
	if err != nil {
		fmt.Println("incr operation error: ", err)
		return err
	}
	fmt.Println("Incr operation success")
	return nil
}

func stringIncrBy(c *grumble.Context) error {
	key := c.Args.String("key")
	amount := c.Args.Int64("amount")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}
	err := newClient().IncrBy(key, amount)
	if err != nil {
		fmt.Println("incrby operation error: ", err)
		return err
	}
	fmt.Println("IncrBy operation success")
	return nil
}

func stringIncrByFloat(c *grumble.Context) error {
	key := c.Args.String("key")
	amount := c.Args.Float64("amount")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}
	err := newClient().IncrByFloat(key, amount)
	if err != nil {
		fmt.Println("incrbyfloat operation error: ", err)
		return err
	}
	fmt.Println("IncrByFloat operation success")
	return nil
}

func stringDecr(c *grumble.Context) error {
	key := c.Args.String("key")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}
	err := newClient().Decr(key)
	if err != nil {
		fmt.Println("decr operation error: ", err)
		return err
	}
	fmt.Println("Decr operation success")
	return nil
}

func stringDecrBy(c *grumble.Context) error {
	key := c.Args.String("key")
	amount := c.Args.Int64("amount")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}
	err := newClient().DecrBy(key, amount)
	if err != nil {
		fmt.Println("decrby operation error: ", err)
		return err
	}
	fmt.Println("DecrBy operation success")
	return nil
}

func stringExists(c *grumble.Context) error {
	key := c.Args.String("key")
	exists, err := newClient().Exists(key)
	if err != nil {
		return err
	}
	if exists {
		fmt.Println("key is exist", key)
	}
	return nil
}

func stringExpire(c *grumble.Context) error {
	key := c.Args.String("key")
	ttl := c.Args.Int64("ttl")
	err := newClient().Expire(key, ttl)
	if err != nil {
		return err
	}
	return nil
}

func stringPersist(c *grumble.Context) error {
	key := c.Args.String("key")
	err := newClient().Persist(key)
	if err != nil {
		fmt.Println("key is not Persist")
		return err
	}
	fmt.Println("key is Persist")
	return nil
}

func stringMGet(c *grumble.Context) error {
	keys := c.Args.StringList("key")
	values, err := newClient().MGet(keys)
	if err != nil {
		fmt.Println("get data error: ", err)
		return err
	}
	fmt.Println(values)
	return nil
}

func stringMSet(c *grumble.Context) error {
	keyValuePairs := c.Args.StringList("key-value")
	if len(keyValuePairs) == 0 || len(keyValuePairs)%2 != 0 {
		return errors.New("invalid number of arguments, must provide key-value pairs")
	}

	var pairs []interface{}
	for i := 0; i < len(keyValuePairs); i += 2 {
		key := keyValuePairs[i]
		value := keyValuePairs[i+1]
		pairs = append(pairs, key, value)
	}

	err := newClient().MSet(pairs...)
	if err != nil {
		fmt.Println("set data error:", err)
		return err
	}
	fmt.Println("Data successfully set.")
	return nil
}

func stringMSetNX(c *grumble.Context) error {
	keyValuePairs := c.Args.StringList("key-value")
	if len(keyValuePairs) == 0 || len(keyValuePairs)%2 != 0 {
		return errors.New("invalid number of arguments, must provide key-value pairs")
	}

	var pairs []interface{}
	for i := 0; i < len(keyValuePairs); i += 2 {
		key := keyValuePairs[i]
		value := keyValuePairs[i+1]
		pairs = append(pairs, key, value)
	}

	err := newClient().MSetNX(pairs...)
	if err != nil {
		fmt.Println("set data error:", err)
		return err
	}
	fmt.Println("Data successfully set.")
	return nil
}
