package client

import (
	"fmt"
	"github.com/ByteStorage/FlyDB/lib/proto/gzset"
	"github.com/ByteStorage/FlyDB/structure"
	"github.com/desertbit/grumble"
	"strconv"
)

func ZSetAdd(c *grumble.Context) error {
	key := c.Args.String("key")
	score := c.Args.Int("score")
	member := c.Args.String("member")
	value := c.Args.String("value")
	if key == "" || member == "" {
		fmt.Println("key or value is empty")
		return nil
	}
	err := newClient().ZAdd(key, score, member, value)
	if err != nil {
		fmt.Println("SAdd data error: ", err)
		return err
	}
	fmt.Println("SAdd data success")
	return nil
}

func ZSetAdds(c *grumble.Context) error {
	key := c.Args.String("key")
	zsetvalues := c.Args.StringList("members")
	if key == "" {
		fmt.Println("key or value is empty")
		return nil
	}
	var zsetItems []structure.ZSetValue
	// Ensure that the length of zsetvalues is a multiple of 3 to avoid index out of range error
	if len(zsetvalues)%3 != 0 {
		fmt.Println("Invalid number of elements in zsetvalues")
		return nil
	}

	for i := 0; i < len(zsetvalues); i += 3 {
		score := zsetvalues[i]
		member := zsetvalues[i+1]
		value := zsetvalues[i+2]
		scoreInt, err := strconv.Atoi(score)
		if err != nil {
			fmt.Printf("Error converting score '%s' to int: %v\n", score, err)
			return err
		}

		zsetItem := structure.ZSetValue{
			Score:  scoreInt,
			Member: member,
			Value:  value,
		}

		zsetItems = append(zsetItems, zsetItem)
	}
	err := newClient().ZAdds(key, zsetItems...)
	if err != nil {
		fmt.Println("SAdds data error: ", err)
		return err
	}
	fmt.Println("SAdds data success")
	return nil
}

func ZSetRem(c *grumble.Context) error {
	key := c.Args.String("key")
	member := c.Args.String("member")
	if key == "" || member == "" {
		fmt.Println("key or value is empty")
		return nil
	}

	err := newClient().ZRem(key, member)
	if err != nil {
		fmt.Println("SRem data error: ", err)
		return err
	}
	fmt.Println("SRem data success")
	return nil
}

func ZSetRems(c *grumble.Context) error {
	key := c.Args.String("key")
	members := c.Args.StringList("members")
	if key == "" || len(members) == 0 {
		fmt.Println("key or members are empty")
		return nil
	}

	err := newClient().ZRems(key, members)
	if err != nil {
		fmt.Println("ZRems data error: ", err)
		return err
	}
	fmt.Println("ZRems data success")
	return nil
}

func ZSetScore(c *grumble.Context) error {
	key := c.Args.String("key")
	member := c.Args.String("member")
	if key == "" || member == "" {
		fmt.Println("key or member is empty")
		return nil
	}

	score, err := newClient().ZScore(key, member)
	if err != nil {
		fmt.Println("ZScore data error: ", err)
		return err
	}
	fmt.Println("ZScore:", score)
	return nil
}

func ZSetRank(c *grumble.Context) error {
	key := c.Args.String("key")
	member := c.Args.String("member")
	if key == "" || member == "" {
		fmt.Println("key or member is empty")
		return nil
	}

	rank, err := newClient().ZRank(key, member)
	if err != nil {
		fmt.Println("ZRank data error: ", err)
		return err
	}
	fmt.Println("ZRank:", rank)
	return nil
}

func ZSetRevRank(c *grumble.Context) error {
	key := c.Args.String("key")
	member := c.Args.String("member")
	if key == "" || member == "" {
		fmt.Println("key or member is empty")
		return nil
	}

	rank, err := newClient().ZRevRank(key, member)
	if err != nil {
		fmt.Println("ZRevRank data error: ", err)
		return err
	}
	fmt.Println("ZRevRank:", rank)
	return nil
}

func ZSetRange(c *grumble.Context) error {
	key := c.Args.String("key")
	start := c.Args.Int("start")
	stop := c.Args.Int("stop")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}

	values, err := newClient().ZRange(key, int32(start), int32(stop))
	if err != nil {
		fmt.Println("ZRange data error: ", err)
		return err
	}
	for _, zsetValue := range values {
		if value, ok := zsetValue.Value.(*gzset.ZSetValue_StringValue); ok {
			fmt.Printf("aaa Score: %d, Member: %s, Value: %s\n", zsetValue.Score, zsetValue.Member, value)
		} else {
			fmt.Printf("bbb Score: %d, Member: %s, Value: %v\n", zsetValue.Score, zsetValue.Member, zsetValue.Value)
		}
	}
	return nil
}

func ZSetCount(c *grumble.Context) error {
	key := c.Args.String("key")
	min := c.Args.Int("min")
	max := c.Args.Int("max")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}

	count, err := newClient().ZCount(key, int32(min), int32(max))
	if err != nil {
		fmt.Println("ZCount data error: ", err)
		return err
	}
	fmt.Println("ZCount:", count)
	return nil
}

func ZSetRevRange(c *grumble.Context) error {
	key := c.Args.String("key")
	start := c.Args.Int("start")
	stop := c.Args.Int("stop")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}

	values, err := newClient().ZRevRange(key, int32(start), int32(stop))
	if err != nil {
		fmt.Println("ZRevRange data error: ", err)
		return err
	}
	fmt.Println("ZRevRange:", values)
	return nil
}

func ZSetCard(c *grumble.Context) error {
	key := c.Args.String("key")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}

	count, err := newClient().ZCard(key)
	if err != nil {
		fmt.Println("ZCard data error: ", err)
		return err
	}
	fmt.Println("ZCard:", count)
	return nil
}

func ZSetIncrBy(c *grumble.Context) error {
	key := c.Args.String("key")
	member := c.Args.String("member")
	increment := c.Args.Int("increment")
	if key == "" || member == "" {
		fmt.Println("key or member is empty")
		return nil
	}

	err := newClient().ZIncrBy(key, member, int32(increment))
	if err != nil {
		fmt.Println("ZIncrBy data error: ", err)
		return err
	}
	return nil
}

func interfaceToBytes(value interface{}) []byte {
	switch value := value.(type) {

	case string:
		return []byte(value)
	case int:
		return []byte(strconv.Itoa(value))
	case int64:
		return []byte(strconv.FormatInt(value, 10))
	case float64:
		return []byte(strconv.FormatFloat(value, 'f', -1, 64))
	case bool:
		return []byte(strconv.FormatBool(value))
	case []byte:
		return value
	default:
		return nil
	}
}
