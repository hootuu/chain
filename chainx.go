package chainx

import (
	"fmt"
	"github.com/hootuu/domain/chain"
	"github.com/hootuu/domain/letter"
	"github.com/hootuu/utils/errors"
	"github.com/hootuu/utils/sys"
	"github.com/rs/xid"
	"sync"
	"time"
)

type Tail struct {
	block *chain.Block
	cid   chain.Cid
}

func (t *Tail) Next(data chain.Cid) *chain.Block {
	return t.block.Next(data, t.cid)
}

type Wait struct {
	signalling string
	block      *chain.Block
	timestamp  time.Time
	consensus  chan bool
}

func NewWait(block *chain.Block) *Wait {
	return &Wait{
		signalling: block.SerializeString() + xid.New().String(),
		block:      block,
		consensus:  make(chan bool, 1),
	}
}

func (w *Wait) GetBlock() *chain.Block {
	return w.block
}

func (w *Wait) GetSignalling() string {
	return w.signalling
}

func (w *Wait) SomeOneConfirmed() {
	w.consensus <- true
}

func (w *Wait) Wait() {
	count := 0
	for {
		select {
		case <-w.consensus:
			count += 1
			sys.Info("-----", count) //todo
			if count >= 3 {
				//close(w.consensus) //todo
				return
			}
		case <-time.After(5 * time.Minute):
			fmt.Println("wait.timeout")
		}
	}
}

type ChainX struct {
	chn  chain.Chain
	tail *Tail
	buf  chan chain.Cid
	wait *Wait
	lock sync.Mutex

	Link []string
}

func New(chn chain.Chain) (*ChainX, *errors.Error) {
	x := &ChainX{
		chn: chn,
		buf: make(chan chain.Cid, 1),
	}
	err := x.reload()
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (x *ChainX) reload() *errors.Error {
	//TODO add bolt
	tailCid, b, err := chain.NewHeadBlock(x.chn)
	if err != nil {
		return err
	}
	x.tail = &Tail{
		block: b,
		cid:   tailCid,
	}
	//x.tail =
	//x.wait =
	return nil
}

func (x *ChainX) GetWait() *Wait {
	return x.wait
}

func (x *ChainX) IsSameWait(signalling string, serializeStr string) bool {
	x.lock.Lock()
	defer x.lock.Unlock()
	if x.wait != nil {
		//sys.Info("x.wait.signalling == signalling : ", x.wait.signalling == signalling)
		//sys.Info("x.wait.block.SerializeString() == serializeStr : ", x.wait.block.SerializeString() == serializeStr)
		return x.wait.signalling == signalling && x.wait.block.SerializeString() == serializeStr
	} else {
		sys.Error("x.wait is nil.....")
	}
	return false
}

func (x *ChainX) Tail() *Tail {
	return x.tail
}

func (x *ChainX) Submit(data chain.Cid) *errors.Error {
	x.buf <- data
	return nil
}

func (x *ChainX) doForInput(data chain.Cid) {
	x.lock.Lock()
	x.wait = NewWait(x.tail.Next(data))
	x.lock.Unlock()
	go func() {
		x.Consensus()
	}()
	x.wait.Wait()
	x.lock.Lock()
	defer x.lock.Unlock()
	xCid, err := chain.GetStone().Inscribe(x.wait.block)
	if err != nil {
		sys.Error(err)
	}
	x.tail = &Tail{
		block: x.wait.block,
		cid:   xCid,
	}
	x.Link = append(x.Link, x.tail.cid)
	sys.Error("--->>will write to chain-->> ", x.tail.cid, "<<--", x.tail.block.Previous)
}

func (x *ChainX) StartUp() *errors.Error {
	letter.PostOffice().Register(x)
	go func() {
		//fmt.Println("wait for submit ....")
		for {
			select {
			case d := <-x.buf:
				//todo
				fmt.Println("do handle ....::", d)
				x.doForInput(d)
			}
		}
	}()
	return nil
}
