package chainx

import (
	"encoding/json"
	"github.com/hootuu/domain/chain"
	"github.com/hootuu/domain/letter"
	"github.com/hootuu/utils/errors"
	"github.com/hootuu/utils/logger"
	"go.uber.org/zap"
)

const (
	TellMeYourWaitTopic      = "tell.me.your.wait"
	ReplyTellMeYourWaitTopic = "reply." + TellMeYourWaitTopic
)

type TellMeYourWaitPayload struct {
	Chn        chain.Chain `bson:"c" json:"c"`
	Signalling string      `bson:"s" json:"s"`
	BlockNumb  int64       `bson:"bn" json:"bn"`
}

func TellMeYourWaitPayloadOf(str string) (*TellMeYourWaitPayload, *errors.Error) {
	var p TellMeYourWaitPayload
	nErr := json.Unmarshal([]byte(str), &p)
	if nErr != nil {
		logger.Logger.Error("TellMeYoursPayloadOf.json.Unmarshal error", zap.String("payload", str))
		return nil, errors.Sys("invalid payload")
	}
	return &p, nil
}

func NewTellMeYoursLetter(p TellMeYourWaitPayload) *letter.Letter {
	return letter.NewLetter(
		TellMeYourWaitTopic,
		p,
	)
}

type ReplyTellMeYourWaitPayload struct {
	Chn            chain.Chain `bson:"c" json:"c"`
	Signalling     string      `bson:"s" json:"s"`
	BlockNumb      int64       `bson:"bn" json:"bn"`
	BlockSerialize string      `bson:"bs" json:"bs"`
}

func ReplyTellMeYourWaitPayloadOf(str string) (*ReplyTellMeYourWaitPayload, *errors.Error) {
	var p ReplyTellMeYourWaitPayload
	nErr := json.Unmarshal([]byte(str), &p)
	if nErr != nil {
		logger.Logger.Error("ReplyTellMeYourWaitPayload.json.Unmarshal error", zap.String("payload", str))
		return nil, errors.Sys("invalid payload")
	}
	return &p, nil
}

func NewReplyTellMeYourWaitLetter(p ReplyTellMeYourWaitPayload) *letter.Letter {
	return letter.NewLetter(
		ReplyTellMeYourWaitTopic,
		p,
	)
}

func (x *ChainX) Consensus() {
	letter.PostOffice().Broadcast(NewTellMeYoursLetter(TellMeYourWaitPayload{
		Chn:        x.chn,
		Signalling: x.GetWait().GetSignalling(),
		BlockNumb:  x.GetWait().GetBlock().Numb,
	}))
}
