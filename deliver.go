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
	DeliverTopic = "deliver"
)

type DeliverPayload struct {
	Chn  chain.Chain `bson:"c" json:"c"`
	Data chain.Cid   `bson:"d" json:"d"`
}

func DeliverPayloadOf(str string) (*DeliverPayload, *errors.Error) {
	var p DeliverPayload
	nErr := json.Unmarshal([]byte(str), &p)
	if nErr != nil {
		logger.Logger.Error("DeliverPayloadOf.json.Unmarshal error", zap.String("payload", str))
		return nil, errors.Sys("invalid payload")
	}
	return &p, nil
}

func NewDeliverLetter(p DeliverPayload) *letter.Letter {
	return letter.NewLetter(
		DeliverTopic,
		p,
	)
}

func Deliver(chn chain.Chain, data chain.Cid) {
	letter.PostOffice().Broadcast(NewDeliverLetter(DeliverPayload{
		Chn:  chn,
		Data: data,
	}))
}
