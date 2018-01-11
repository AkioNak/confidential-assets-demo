// Copyright (c) 2017 DG Lab
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/*
Package lib provides worker model functions.

There is a very simple paralll traffic control framework.
*/
package lib

import (
	"math/rand"
	"sync"
	"time"
)

const (
	workerThreds = 3
	// sleepAvarage = 3000
)

type CallbackFunc func([]interface{}) []interface{}

type RequestItem struct {
	Callback CallbackFunc
	Params   []interface{}
	Response chan []interface{}
}

type Dispatcher struct {
	WaitGroup    sync.WaitGroup
	RequestQueue chan RequestItem
}

var once sync.Once
var instance *Dispatcher

// GetDispatcherInstance returns dispatcher instance.
func GetDispatcherInstance() *Dispatcher {
	once.Do(func() {
		rand.Seed(time.Now().UnixNano())
		instance = &Dispatcher{RequestQueue: make(chan RequestItem, 10)}
		instance.startWorker()
	})
	return instance
}

func (d *Dispatcher) Enqueue(callback CallbackFunc, params []interface{}) chan []interface{} {
	ch := make(chan []interface{})
	req := RequestItem{Callback: callback, Params: params, Response: ch}
	d.RequestQueue <- req
	return ch
}

func (d *Dispatcher) Stop() {
	close(d.RequestQueue)
	d.WaitGroup.Wait()
}

func (d *Dispatcher) worker(no int) {
	defer d.WaitGroup.Done()
	for {
		req, ok := <-d.RequestQueue // 'ok' will be false when q is closed
		if !ok {
			return
		}

		res := req.Callback(req.Params)
		req.Response <- res

		// time.Sleep(time.Duration(rand.Int63n(sleepAvarage*2)) * time.Microsecond)
	}
}

func (d *Dispatcher) startWorker() {
	for i := 0; i < workerThreds; i++ {
		d.WaitGroup.Add(1)
		go d.worker(i)
	}
}
