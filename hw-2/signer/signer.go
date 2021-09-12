package main

import (
	"fmt"
	"sync"
	"time"
)

const PIPELINE_DEBUG = true

func ExecutePipeline(jobs ...job) {
	wg := new(sync.WaitGroup)
	var prev_chan chan interface{} = nil

	for i, job_i := range jobs {
		cur_chan := make(chan interface{}, 1)

		wg.Add(1)
		go func(wg *sync.WaitGroup, job_i job, i int, in, out chan interface{}) {
			defer wg.Done()

			if PIPELINE_DEBUG {
				fmt.Println("#", i, job_i, "started. input", in, ", output", out)
			}

			job_i(in, out)
			time.Sleep(time.Second)
			close(out)

			if PIPELINE_DEBUG {
				fmt.Println("#", i, job_i, "finished. Channel", out, "closed")
			}
		}(wg, job_i, i, prev_chan, cur_chan)

		prev_chan = cur_chan
	}
	wg.Wait()
}
