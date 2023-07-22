package client

import (
	"fmt"
	"github.com/desertbit/grumble"
)

func stringLPushData(c *grumble.Context) error {
	key := c.Args.String("key")
	value := c.Args.String("value")
	if key == "" || value == "" {
		fmt.Println("key or value is empty")
		return nil
	}

	err := newClient().LPush(key, value)
	if err != nil {
		fmt.Println("LPush data error: ", err)
		return err
	}
	fmt.Println("LPush data success")
	return nil
}

// CLI for LPushs command
func stringLPushsData(c *grumble.Context) error {
	key := c.Args.String("key")
	values := c.Args.StringList("values")
	if key == "" || len(values) == 0 {
		fmt.Println("key or values is empty")
		return nil
	}

	interfaceValues := make([]interface{}, len(values))
	for i, v := range values {
		interfaceValues[i] = v
	}

	err := newClient().LPushs(key, interfaceValues)
	if err != nil {
		fmt.Println("LPushs data error:", err)
		return err
	}
	fmt.Println("LPushs data success")
	return nil
}

// CLI for RPush command
func stringRPushData(c *grumble.Context) error {
	key := c.Args.String("key")
	value := c.Args.String("value")
	if key == "" || value == "" {
		fmt.Println("key or value is empty")
		return nil
	}

	err := newClient().RPush(key, value)
	if err != nil {
		fmt.Println("RPush data error:", err)
		return err
	}
	fmt.Println("RPush data success")
	return nil
}

// CLI for RPushs command
func stringRPushsData(c *grumble.Context) error {
	key := c.Args.String("key")
	values := c.Args.StringList("values")
	if key == "" || len(values) == 0 {
		fmt.Println("key or values is empty")
		return nil
	}

	interfaceValues := make([]interface{}, len(values))
	for i, v := range values {
		interfaceValues[i] = v
	}
	println(interfaceValues)
	err := newClient().RPushs(key, interfaceValues)
	if err != nil {
		fmt.Println("RPushs data error:", err)
		return err
	}
	fmt.Println("RPushs data success")
	return nil
}

func stringLPopData(c *grumble.Context) error {
	key := c.Args.String("key")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}

	value, err := newClient().LPop(key)
	if err != nil {
		fmt.Println("LPop data error:", err)
		return err
	}

	fmt.Println("LPop data success:", value)
	return nil
}

// CLI for RPop command
func stringRPopData(c *grumble.Context) error {
	key := c.Args.String("key")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}

	value, err := newClient().RPop(key)
	if err != nil {
		fmt.Println("RPop data error:", err)
		return err
	}

	fmt.Println("RPop data success:", value)
	return nil
}

// CLI for LRange command
func stringLRangeData(c *grumble.Context) error {
	key := c.Args.String("key")
	start := c.Args.Int("start")
	stop := c.Args.Int("stop")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}

	values, err := newClient().LRange(key, start, stop)
	if err != nil {
		fmt.Println("LRange data error:", err)
		return err
	}

	fmt.Println("LRange data success:", values)
	return nil
}

// CLI for LLen command
func stringLLenData(c *grumble.Context) error {
	key := c.Args.String("key")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}

	length, err := newClient().LLen(key)
	if err != nil {
		fmt.Println("LLen data error:", err)
		return err
	}

	fmt.Println("LLen data success:", length)
	return nil
}

// CLI for LRem command
func stringLRemData(c *grumble.Context) error {
	key := c.Args.String("key")
	count := c.Args.Int("count")
	value := c.Args.String("value")
	if key == "" || value == "" {
		fmt.Println("key or value is empty")
		return nil
	}

	err := newClient().LRem(key, int32(count), value)
	if err != nil {
		fmt.Println("LRem data error:", err)
		return err
	}
	fmt.Println("LRem data success")
	return nil
}

// CLI for LIndex command
func stringLIndexData(c *grumble.Context) error {
	key := c.Args.String("key")
	index := c.Args.Int("index")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}

	value, err := newClient().LIndex(key, index)
	if err != nil {
		fmt.Println("LIndex data error:", err)
		return err
	}
	fmt.Println("LIndex data success:", value)
	return nil
}

// CLI for LSet command
func stringLSetData(c *grumble.Context) error {
	key := c.Args.String("key")
	index := c.Args.Int("index")
	value := c.Args.String("value")
	if key == "" || value == "" {
		fmt.Println("key or value is empty")
		return nil
	}

	err := newClient().LSet(key, index, value)
	if err != nil {
		fmt.Println("LSet data error:", err)
		return err
	}
	fmt.Println("LSet data success")
	return nil
}

// CLI for LTrim command
func stringLTrimData(c *grumble.Context) error {
	key := c.Args.String("key")
	start := c.Args.Int("start")
	stop := c.Args.Int("stop")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}

	err := newClient().LTrim(key, start, stop)
	if err != nil {
		fmt.Println("LTrim data error:", err)
		return err
	}
	fmt.Println("LTrim data success")
	return nil
}
