package main

import "sync"

// TargetMap maps a WorkerSpec+TargetSpec to a Worker and Target
type TargetMap struct {
	sync.Mutex
	workers map[WorkerSpec]*Worker
}

var tm TargetMap = TargetMap{
	workers: make(map[WorkerSpec]*Worker),
}

func GetTarget(ws WorkerSpec) *Target {
	// retrieve Worker
	tm.Lock()
	w, ok := tm.workers[ws]
	if !ok {
		w = NewWorker(ws)
		tm.workers[ws] = w
	}
	tm.Unlock()

	// retrieve Target
	return w.GetWorkerTarget(ws.host)
}
