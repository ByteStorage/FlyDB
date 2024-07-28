package client

import (
	"fmt"
	"math"

	"github.com/desertbit/grumble"
	pbany "github.com/golang/protobuf/ptypes/any"

	"github.com/ByteStorage/FlyDB/db/grpc/client"
	"github.com/ByteStorage/FlyDB/lib/proto/gzset"
)

func checkIsEmpty(key string, value string) bool {
	if value == "" {
		fmt.Printf("%s is empty\n", key)
		return true
	}

	return false
}

func checkIsInt32(key string, value int) bool {
	if value < math.MinInt32 || value > math.MaxInt32 {
		fmt.Printf("%s is not int32\n", key)
		return true
	}

	return false
}

func string2Any(value string) *pbany.Any {
	return &pbany.Any{
		Value: []byte(value),
	}
}

func checkRange(start int, end int) bool {
	if start > end {
		fmt.Println("Start is greater than end")
		return true
	}

	return false
}

func ZSetAdd(ctx *grumble.Context) error {
	var (
		key       = ctx.Args.String(KeyArg)
		member    = ctx.Args.String(MemberArg)
		score     = ctx.Args.Int(ScoreArg)
		value     = ctx.Args.String(ValueArg)
		zSetValue = &gzset.ZSetValue{}
		err       error
	)

	if checkIsEmpty(KeyArg, key) || checkIsEmpty(MemberArg, member) || checkIsEmpty(ValueArg, value) || checkIsInt32(ScoreArg, score) {
		return nil
	}

	zSetValue.Member = member
	zSetValue.Score = int32(score)
	zSetValue.Value = string2Any(value)

	if _, err = newClient().ZAdd(key, zSetValue); err != nil {
		fmt.Println("ZAdd data error: ", err)

		return err
	}

	fmt.Println("ZAdd data success")

	return nil
}

func ZSetRem(ctx *grumble.Context) error {
	var (
		key    = ctx.Args.String(KeyArg)
		member = ctx.Args.String(MemberArg)
		err    error
	)

	if checkIsEmpty(KeyArg, key) || checkIsEmpty(MemberArg, member) {
		return nil
	}

	if _, err = newClient().ZRem(key, member); err != nil {
		fmt.Println("ZRem data error: ", err)

		return err
	}

	fmt.Println("ZRem data success")

	return nil
}

func ZSetRems(ctx *grumble.Context) error {
	var (
		key        = ctx.Args.String(KeyArg)
		memberList = ctx.Args.StringList(MembersArg)
		err        error
	)

	if checkIsEmpty(KeyArg, key) {
		return nil
	}

	if _, err = newClient().ZRems(key, memberList); err != nil {
		fmt.Println("ZRems data error: ", err)

		return err
	}

	fmt.Println("ZRems data success")

	return nil
}

func ZSetScore(ctx *grumble.Context) error {
	var (
		key      = ctx.Args.String(KeyArg)
		member   = ctx.Args.String(MemberArg)
		response *gzset.ZScoreResponse
		err      error
	)

	if checkIsEmpty(KeyArg, key) || checkIsEmpty(MemberArg, member) {
		return nil
	}

	if response, err = newClient().ZScore(key, member); err != nil {
		fmt.Println("ZScore data error: ", err)

		return err
	}

	fmt.Printf("ZScore data success, score: %d\n", response.Score)

	return nil
}

func ZSetRank(ctx *grumble.Context) error {
	var (
		key      = ctx.Args.String(KeyArg)
		member   = ctx.Args.String(MemberArg)
		response *gzset.ZRankResponse
		err      error
	)

	if checkIsEmpty(KeyArg, key) || checkIsEmpty(MemberArg, member) {
		return nil
	}

	if response, err = newClient().ZRank(key, member); err != nil {
		fmt.Println("ZRank data error: ", err)

		return err
	}

	fmt.Printf("ZRank data success, rank: %d\n", response.Rank)

	return nil
}

func ZSetRevRank(ctx *grumble.Context) error {
	var (
		key      = ctx.Args.String(KeyArg)
		member   = ctx.Args.String(MemberArg)
		response *gzset.ZRevRankResponse
		err      error
	)

	if checkIsEmpty(KeyArg, key) || checkIsEmpty(MemberArg, member) {
		return nil
	}

	if response, err = newClient().ZRevRank(key, member); err != nil {
		fmt.Println("ZRevRank data error: ", err)

		return err
	}

	fmt.Printf("ZRevRank data success, rank: %d\n", response.Rank)

	return nil
}

func printZSetResponse(memberList []*gzset.ZSetValue) {
	for _, member := range memberList {
		fmt.Printf("Member: %s, Score: %d, Value: %s\n", member.Member, member.Score, member.Value.Value)
	}
}

func ZSetRange(ctx *grumble.Context) error {
	var (
		key          = ctx.Args.String(KeyArg)
		start        = ctx.Args.Int(StartArg)
		end          = ctx.Args.Int(EndArg)
		requestRange = &client.Range{}
		response     *gzset.ZRangeResponse
		err          error
	)

	if checkIsEmpty(KeyArg, key) || checkIsInt32(StartArg, start) || checkIsInt32(EndArg, end) || checkRange(start, end) {
		return nil
	}

	requestRange.Start = int32(start)
	requestRange.End = int32(end)

	if response, err = newClient().ZRange(key, requestRange); err != nil {
		fmt.Println("ZRange data error: ", err)

		return err
	}

	fmt.Println("ZRange data success")
	printZSetResponse(response.Members)

	return nil
}

func ZSetCount(ctx *grumble.Context) error {
	var (
		key          = ctx.Args.String(KeyArg)
		start        = ctx.Args.Int(StartArg)
		end          = ctx.Args.Int(EndArg)
		requestRange = &client.Range{}
		response     *gzset.ZCountResponse
		err          error
	)

	if checkIsEmpty(KeyArg, key) || checkIsInt32(StartArg, start) || checkIsInt32(EndArg, end) || checkRange(start, end) {
		return nil
	}

	requestRange.Start = int32(start)
	requestRange.End = int32(end)

	if response, err = newClient().ZCount(key, requestRange); err != nil {
		fmt.Println("ZCount data error: ", err)

		return err
	}

	fmt.Printf("ZCount data success, count: %d\n", response.Count)

	return nil
}

func ZSetRevRange(ctx *grumble.Context) error {
	var (
		key          = ctx.Args.String(KeyArg)
		start        = ctx.Args.Int(StartArg)
		end          = ctx.Args.Int(EndArg)
		requestRange = &client.Range{}
		response     *gzset.ZRevRangeResponse
		err          error
	)

	if checkIsEmpty(KeyArg, key) || checkIsInt32(StartArg, start) || checkIsInt32(EndArg, end) || checkRange(start, end) {
		return nil
	}

	requestRange.Start = int32(start)
	requestRange.End = int32(end)

	if response, err = newClient().ZRevRange(key, requestRange); err != nil {
		fmt.Println("ZRevRange data error: ", err)

		return err
	}

	fmt.Println("ZRevRange data success")
	printZSetResponse(response.Members)

	return nil

}

func ZSetCard(ctx *grumble.Context) error {
	var (
		key      = ctx.Args.String(KeyArg)
		response *gzset.ZCardResponse
		err      error
	)

	if checkIsEmpty(KeyArg, key) {
		return nil
	}

	if response, err = newClient().ZCard(key); err != nil {
		fmt.Println("ZCard data error: ", err)

		return err
	}

	fmt.Printf("ZCard data success, count: %d\n", response.Count)

	return nil
}

func ZSetIncrBy(ctx *grumble.Context) error {
	var (
		key      = ctx.Args.String(KeyArg)
		member   = ctx.Args.String(MemberArg)
		incrBy   = ctx.Args.Int(IncrByArg)
		incr     = &client.Incr{}
		response *gzset.ZIncrByResponse
		err      error
	)

	if checkIsEmpty(KeyArg, key) || checkIsEmpty(MemberArg, member) || checkIsInt32(IncrByArg, incrBy) {
		return nil
	}

	incr.Member = member
	incr.Inc = int32(incrBy)

	if response, err = newClient().ZIncrBy(key, incr); err != nil {
		fmt.Println("ZIncrBy data error: ", err)

		return err
	}

	fmt.Printf("ZIncrBy data success, new score: %d\n", response.NewScore)

	return nil
}
