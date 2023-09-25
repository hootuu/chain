package main

import (
	"fmt"
	"github.com/hootuu/chainx"
	"github.com/hootuu/domain/chain"
	"github.com/hootuu/utils/sys"
	"math/rand"
	"time"
)

var xCount = 50
var xArr []*chainx.ChainX

func initPeer() {

	for i := 0; i < xCount; i++ {
		x, err := chainx.New(chain.Chain{
			VN:       "A",
			Scope:    "B",
			Category: "C",
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		xArr = append(xArr, x)
	}
}

func runPeer(id int) {
	fmt.Println("run id ", id)
	err := xArr[id].StartUp()
	if err != nil {
		fmt.Println("StartUp::", err)
		return
	}
	fmt.Println("StartUp::", id)
}

func main() {
	initPeer()
	for i := 0; i < xCount; i++ {
		go runPeer(i)
	}
	go func() {
		time.Sleep(3 * time.Second)
		for {
			chainx.Deliver(chain.Chain{
				VN:       "A",
				Scope:    "B",
				Category: "C",
			}, fmt.Sprintf("XX_%d", time.Now().UnixMilli()))
			w := rand.Intn(5000)
			time.Sleep(time.Duration(w) * time.Millisecond)
		}
	}()

	go func() {
		for {
			time.Sleep(15 * time.Second)
			for i := 0; i < xCount; i++ {
				sys.Error("[", i+1, "]", xArr[i].Link)
			}
		}

	}()
	time.Sleep(20 * time.Minute)
}
