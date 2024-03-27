package main

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var client *resty.Client

const USER_NUMBER = 3000

func main() {
	client = resty.New()
	wg := sync.WaitGroup{}
	sCh := make(chan struct{}, USER_NUMBER)
	fCh := make(chan struct{}, USER_NUMBER)

	start := time.Now()
	for i := 0; i < USER_NUMBER; i++ {
		wg.Add(1)
		go Work(&wg, sCh, fCh)
	}

	wg.Wait()
	fmt.Printf("%d 位用户秒杀到优惠券, %d 用户未抢到优惠券, 总用时 %d 毫秒。\n", len(sCh), len(fCh), time.Since(start).Milliseconds())
}

func Work(wg *sync.WaitGroup, sCh, fCh chan struct{}) {
	defer wg.Done()
	if err := Take(); err != nil {
		fCh <- struct{}{}
		return
	}
	sCh <- struct{}{}
}

func Take() error {
	resp, err := client.R().SetHeader("Accept", "application/json").Get("http://127.0.0.1:8080/take")
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.New("call take error")
	}
	return nil
}
