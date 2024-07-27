package client

import (
	"fmt"

	"github.com/desertbit/grumble"
	"github.com/pkg/errors"
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
func hashHExistsKey(c *grumble.Context) error {
	key := c.Args.String("key")
	field := c.Args.String("field")
	if key == "" || field == "" {
		fmt.Println("key or field is empty")
		return nil
	}
	exists, err := newClient().HExists(key, field)
	if err != nil {
		fmt.Println("HExists error: ", err)
		return err
	}
	if exists {
		fmt.Println("Field exists in the hash")
	} else {
		fmt.Println("Field does not exist in the hash")
	}
	return nil
}

func hashHLenKey(c *grumble.Context) error {
	key := c.Args.String("key")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}
	length, err := newClient().HLen(key)
	if err != nil {
		fmt.Println("HLen error: ", err)
		return err
	}
	fmt.Println("Length of the hash:", length)
	return nil
}

func hashHUpdateKey(c *grumble.Context) error {
	key := c.Args.String("key")
	field := c.Args.String("field")
	value := c.Args.String("value")
	if key == "" || field == "" {
		fmt.Println("key or field is empty")
		return nil
	}
	err := newClient().HUpdate(key, field, value)
	if err != nil {
		fmt.Println("HUpdate error: ", err)
		return err
	}
	fmt.Println("Hash updated successfully")
	return nil
}

func hashHIncrByKey(c *grumble.Context) error {
	key := c.Args.String("key")
	field := c.Args.String("field")
	value := c.Args.Int64("value")
	if key == "" || field == "" {
		fmt.Println("key or field is empty")
		return nil
	}
	result, err := newClient().HIncrBy(key, field, value)
	if err != nil {
		fmt.Println("HIncrBy error: ", err)
		return err
	}
	fmt.Println("Result:", result)
	return nil
}

func hashHIncrByFloatKey(c *grumble.Context) error {
	key := c.Args.String("key")
	field := c.Args.String("field")
	value := c.Args.Float64("value")
	if key == "" || field == "" {
		fmt.Println("key or field is empty")
		return nil
	}
	result, err := newClient().HIncrByFloat(key, field, value)
	if err != nil {
		fmt.Println("HIncrByFloat error: ", err)
		return err
	}
	fmt.Println("Result:", result)
	return nil
}

func hashHDecrByKey(c *grumble.Context) error {
	key := c.Args.String("key")
	field := c.Args.String("field")
	value := c.Args.Int64("value")
	if key == "" || field == "" {
		fmt.Println("key or field is empty")
		return nil
	}
	result, err := newClient().HDecrBy(key, field, value)
	if err != nil {
		fmt.Println("HDecrBy error: ", err)
		return err
	}
	fmt.Println("Result:", result)
	return nil
}

func hashHStrLenKey(c *grumble.Context) error {
	key := c.Args.String("key")
	field := c.Args.String("field")
	if key == "" || field == "" {
		fmt.Println("key or field is empty")
		return nil
	}
	length, err := newClient().HStrLen(key, field)
	if err != nil {
		fmt.Println("HStrLen error: ", err)
		return err
	}
	fmt.Println("Length of the string value:", length)
	return nil
}

func hashHMoveKey(c *grumble.Context) error {
	key := c.Args.String("key")
	field := c.Args.String("field")
	dest := c.Args.String("dest")
	if key == "" || field == "" || dest == "" {
		fmt.Println("key, field, or dest is empty")
		return nil
	}
	err := newClient().HMove(key, dest, field)
	if err != nil {
		fmt.Println("HMove error: ", err)
		return err
	}
	fmt.Println("Field moved successfully")
	return nil
}

func hashHSetNXKey(c *grumble.Context) error {
	key := c.Args.String("key")
	field := c.Args.String("field")
	value := c.Args.String("value")
	if key == "" || field == "" {
		fmt.Println("key or field is empty")
		return nil
	}
	err := newClient().HSetNX(key, field, value)
	if err != nil {
		fmt.Println("HSetNX error: ", err)
		return err
	}
	fmt.Println("Field set successfully")
	return nil
}

func hashHType(c *grumble.Context) error {
	key := c.Flags.String("key")
	field := c.Flags.String("field")

	if key == "" {
		return errors.New("key is empty")
	}

	if field == "" {
		return errors.New("field is empty")
	}
	hashType, err := newClient().HType(key, field)
	if err != nil {
		return errors.Wrap(err, "HType error")
	}

	fmt.Println("Type of field", field, "in hash", key, "is", hashType)
	return nil
}
