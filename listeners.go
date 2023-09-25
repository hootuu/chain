package chainx

import (
	"fmt"
	"github.com/hootuu/domain/letter"
	"github.com/hootuu/utils/errors"
	"github.com/hootuu/utils/logger"
	"github.com/hootuu/utils/sys"
	"go.uber.org/zap"
)

func (x *ChainX) GetTopic() []string {
	return []string{TellMeYourWaitTopic, ReplyTellMeYourWaitTopic, DeliverTopic}
}

func (x *ChainX) Deal(ltr letter.ILetter) *errors.Error {
	if ltr == nil {
		return errors.Sys("require letter")
	}
	switch ltr.GetTopic() {
	case TellMeYourWaitTopic:
		return x.doTellMeYourWaitTopic(ltr)
	case ReplyTellMeYourWaitTopic:
		return x.doReplyTellMeYourWaitTopic(ltr)
	case DeliverTopic:
		return x.doDeliverTopic(ltr)
	default:
		return errors.Sys("invalid topic:" + ltr.GetTopic())
	}
}

func (x *ChainX) doTellMeYourWaitTopic(ltr letter.ILetter) *errors.Error {
	fmt.Println("dellTellMeYourWaitTopic", ltr.GetPayload())
	tellMeYourWaitPayload, err := TellMeYourWaitPayloadOf(ltr.GetPayload())
	if err != nil {
		return err
	}
	if !x.chn.Same(&tellMeYourWaitPayload.Chn) {
		logger.Logger.Info("no need to pay attention:", zap.String("payload", ltr.GetPayload()))
		return nil
	}
	if x.GetWait() == nil {
		return nil
	}
	wb := x.GetWait().GetBlock()
	wbSerializeStr := wb.SerializeString()
	replyLetter := NewReplyTellMeYourWaitLetter(ReplyTellMeYourWaitPayload{
		Chn:            x.chn,
		Signalling:     tellMeYourWaitPayload.Signalling,
		BlockNumb:      wb.Numb,
		BlockSerialize: wbSerializeStr,
	})
	letter.PostOffice().Broadcast(replyLetter)
	return nil
}

func (x *ChainX) doReplyTellMeYourWaitTopic(ltr letter.ILetter) *errors.Error {
	//fmt.Println("doReplyTellMeYourWaitTopic", ltr.GetPayload())
	replyTellMeYourWaitPayload, err := ReplyTellMeYourWaitPayloadOf(ltr.GetPayload())
	if err != nil {
		return err
	}
	if !x.chn.Same(&replyTellMeYourWaitPayload.Chn) {
		//logger.Logger.Info("no need to pay attention:", zap.String("payload", ltr.GetPayload()))
		return nil
	}
	if !x.IsSameWait(replyTellMeYourWaitPayload.Signalling, replyTellMeYourWaitPayload.BlockSerialize) {
		logger.Logger.Info("Block Serialize not matched, could be an attack, ignore it")
		return nil
	}
	//wt := x.GetWait()
	//if wt.GetSignalling() != replyTellMeYourWaitPayload.Signalling {
	//	//logger.Logger.Info("signalling not matched, ignore it")
	//	return nil
	//}
	//wb := wt.GetBlock()
	//wbSerializeStr := wb.SerializeString()
	//if wbSerializeStr != replyTellMeYourWaitPayload.BlockSerialize {
	//	logger.Logger.Info("Block Serialize not matched, could be an attack, ignore it",
	//		zap.String("wb.data", wb.Data))
	//	return nil
	//}
	sys.Info("ok. will SomeOneConfirmed.........", replyTellMeYourWaitPayload.Signalling)
	x.GetWait().SomeOneConfirmed()
	return nil
}

func (x *ChainX) doDeliverTopic(ltr letter.ILetter) *errors.Error {
	//fmt.Println("doReplyTellMeYourWaitTopic", ltr.GetPayload())
	deliverPayload, err := DeliverPayloadOf(ltr.GetPayload())
	if err != nil {
		return err
	}
	if !x.chn.Same(&deliverPayload.Chn) {
		//logger.Logger.Info("no need to pay attention:", zap.String("payload", ltr.GetPayload()))
		return nil
	}
	err = x.Submit(deliverPayload.Data)
	if err != nil {
		return err
	}
	//if !x.IsSameWait(replyTellMeYourWaitPayload.Signalling, replyTellMeYourWaitPayload.BlockSerialize) {
	//	logger.Logger.Info("Block Serialize not matched, could be an attack, ignore it")
	//	return nil
	//}
	//wt := x.GetWait()
	//if wt.GetSignalling() != replyTellMeYourWaitPayload.Signalling {
	//	//logger.Logger.Info("signalling not matched, ignore it")
	//	return nil
	//}
	//wb := wt.GetBlock()
	//wbSerializeStr := wb.SerializeString()
	//if wbSerializeStr != replyTellMeYourWaitPayload.BlockSerialize {
	//	logger.Logger.Info("Block Serialize not matched, could be an attack, ignore it",
	//		zap.String("wb.data", wb.Data))
	//	return nil
	//}
	//sys.Info("ok. will SomeOneConfirmed.........", replyTellMeYourWaitPayload.Signalling)
	//x.GetWait().SomeOneConfirmed()
	return nil
}
