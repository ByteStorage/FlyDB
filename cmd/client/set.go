package client

import (
	"fmt"

	"github.com/desertbit/grumble"
)

func SetAdd(c *grumble.Context) error {
	key := c.Args.String("key")
	member := c.Args.String("member")
	if key == "" || member == "" {
		fmt.Println("key or value is empty")
		return nil
	}

	err := newClient().SAdd(key, member)
	if err != nil {
		fmt.Println("SAdd data error: ", err)
		return err
	}
	fmt.Println("SAdd data success")
	return nil
}

func SetAdds(c *grumble.Context) error {
	key := c.Args.String("key")
	members := c.Args.StringList("members")
	if key == "" || len(members) == 0 {
		fmt.Println("key or members are empty")
		return nil
	}

	err := newClient().SAdds(key, members)
	if err != nil {
		fmt.Println("SAdds data error: ", err)
		return err
	}
	fmt.Println("SAdds data success")
	return nil
}

func SetRem(c *grumble.Context) error {
	key := c.Args.String("key")
	member := c.Args.String("member")
	if key == "" || member == "" {
		fmt.Println("key or value is empty")
		return nil
	}

	err := newClient().SRem(key, member)
	if err != nil {
		fmt.Println("SRem data error: ", err)
		return err
	}
	fmt.Println("SRem data success")
	return nil
}

func SetRems(c *grumble.Context) error {
	key := c.Args.String("key")
	members := c.Args.StringList("members")
	if key == "" || len(members) == 0 {
		fmt.Println("key or members are empty")
		return nil
	}

	err := newClient().SRems(key, members)
	if err != nil {
		fmt.Println("SRems data error: ", err)
		return err
	}
	fmt.Println("SRems data success")
	return nil
}

func SetCard(c *grumble.Context) error {
	key := c.Args.String("key")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}

	count, err := newClient().SCard(key)
	if err != nil {
		fmt.Println("SCard data error: ", err)
		return err
	}
	fmt.Println("SCard count:", count)
	return nil
}

func SetMembers(c *grumble.Context) error {
	key := c.Args.String("key")
	if key == "" {
		fmt.Println("key is empty")
		return nil
	}

	members, err := newClient().SMembers(key)
	if err != nil {
		fmt.Println("SMembers data error: ", err)
		return err
	}
	fmt.Println("SMembers data:", members)
	return nil
}

func SetIsMember(c *grumble.Context) error {
	key := c.Args.String("key")
	member := c.Args.String("member")
	if key == "" || member == "" {
		fmt.Println("key or member is empty")
		return nil
	}

	isMember, err := newClient().SIsMember(key, member)
	if err != nil {
		fmt.Println("SIsMember data error: ", err)
		return err
	}
	fmt.Println("SIsMember result:", isMember)
	return nil
}

func SetUnion(c *grumble.Context) error {
	keys := c.Args.StringList("keys")
	if len(keys) == 0 {
		fmt.Println("keys list is empty")
		return nil
	}

	union, err := newClient().SUnion(keys)
	if err != nil {
		fmt.Println("SUnion data error: ", err)
		return err
	}
	fmt.Println("SUnion result:", union)
	return nil
}

func SetInter(c *grumble.Context) error {
	keys := c.Args.StringList("keys")
	if len(keys) == 0 {
		fmt.Println("keys list is empty")
		return nil
	}

	inter, err := newClient().SInter(keys)
	if err != nil {
		fmt.Println("SInter data error: ", err)
		return err
	}
	fmt.Println("SInter result:", inter)
	return nil
}

func SetDiff(c *grumble.Context) error {
	keys := c.Args.StringList("keys")
	if len(keys) == 0 {
		fmt.Println("keys list is empty")
		return nil
	}

	diff, err := newClient().SDiff(keys)
	if err != nil {
		fmt.Println("SDiff data error: ", err)
		return err
	}
	fmt.Println("SDiff result:", diff)
	return nil
}

func SetUnionStore(c *grumble.Context) error {
	destinationKey := c.Args.String("destinationKey")
	keys := c.Args.StringList("keys")
	if destinationKey == "" || len(keys) == 0 {
		fmt.Println("destinationKey or keys list is empty")
		return nil
	}

	err := newClient().SUnionStore(destinationKey, keys)
	if err != nil {
		fmt.Println("SUnionStore data error: ", err)
		return err
	}
	fmt.Println("SUnionStore success")
	return nil
}

func SetInterStore(c *grumble.Context) error {
	destinationKey := c.Args.String("destinationKey")
	keys := c.Args.StringList("keys")
	if destinationKey == "" || len(keys) == 0 {
		fmt.Println("destinationKey or keys list is empty")
		return nil
	}

	err := newClient().SInterStore(destinationKey, keys)
	if err != nil {
		fmt.Println("SInterStore data error: ", err)
		return err
	}
	fmt.Println("SInterStore success")
	return nil
}
