package client

import "math"

const (
	CommonKeyArg     = "key"
	CommonMemberArg  = "member"
	CommonMembersArg = "members"
	CommonValueArg   = "value"

	CommonDefaultEmptyString = ""

	ZSetScoreArg  = "score"
	ZSetStartArg  = "start"
	ZSetEndArg    = "end"
	ZSetIncrByArg = "incrBy"
	ZSetMinArg    = "min"
	ZSetMaxArg    = "max"

	ZSetDefaultScore      = 0
	ZSetDefaultRangeStart = 0
	ZSetDefaultRangeEnd   = math.MaxInt
	ZSetDefaultIncrBy     = 0
	ZSetDefaultMin        = 0
	ZSetDefaultMax        = math.MaxInt

	ZSetDefaultKeyHelp     = "the key of the zset"
	ZSetDefaultMemberHelp  = "the member of the zset"
	ZSetDefaultMembersHelp = "the members of the zset, separated by space, e.g. member1 member2 member3"
	ZSetDefaultValueHelp   = "the value of the zset"
	ZSetDefaultScoreHelp   = "the score of the zset"
	ZSetDefaultStartHelp   = "the start index of the zset"
	ZSetDefaultEndHelp     = "the end index of the zset"
	ZSetDefaultIncrByHelp  = "the increment value of the zset"
	ZSetDefaultMinHelp     = "the min score of the zset"
	ZSetDefaultMaxHelp     = "the max score of the zset"
)

const (
	flyDBServerPort = "8999/tcp"
)
