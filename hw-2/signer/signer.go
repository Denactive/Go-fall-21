package main

import (
	"fmt"
	"sort"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	wg := new(sync.WaitGroup)
	var prevChan chan interface{} = nil

	for i, task := range jobs {
		curChan := make(chan interface{}, 1)

		wg.Add(1)
		go func(wg *sync.WaitGroup, task job, i int, in, out chan interface{}) {
			defer wg.Done()
			defer close(out)

			task(in, out)
		}(wg, task, i, prevChan, curChan)

		prevChan = curChan
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	wg := new(sync.WaitGroup)
	mu := new(sync.Mutex)
	i := 0

	for input := range in {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int, data string) {
			defer wg.Done()

			mu.Lock()
			var md5 string = DataSignerMd5(data)
			mu.Unlock()

			// wait subgroup
			wsg := new(sync.WaitGroup)
			wsg.Add(2)
			crc32DataChannel := make(chan string, MaxInputDataLen)
			crc32Md5Channel := make(chan string, MaxInputDataLen)

			go func(wsg *sync.WaitGroup, i int, out chan<- string, data string) {
				defer wsg.Done()
				defer close(out)

				out <- DataSignerCrc32(data)
			}(wsg, i, crc32DataChannel, data)

			go func(wsg *sync.WaitGroup, i int, out chan<- string, md5 string) {
				defer wsg.Done()
				defer close(out)

				out <- DataSignerCrc32(md5)
			}(wsg, i, crc32Md5Channel, md5)

			wsg.Wait()
			out <- (<-crc32DataChannel + "~" + <-crc32Md5Channel)
		}(wg, i, fmt.Sprint(input.(int)))
		i++
	}
	wg.Wait()
}

const mhCalculationsNum = 6

func MultiHash(in, out chan interface{}) {
	wg := new(sync.WaitGroup)
	var i int

	for input := range in {
		wg.Add(1)

		// MultiHash subworker per data chunk
		go func(wg *sync.WaitGroup, th int, data string, out chan interface{}) {
			defer wg.Done()

			wsg := new(sync.WaitGroup)
			mu := new(sync.Mutex)
			res := make(map[int]string, mhCalculationsNum)

			// MultiHash 1/6 subworker
			for i := 0; i < mhCalculationsNum; i++ {
				wsg.Add(1)
				go func(wg *sync.WaitGroup, i int) {
					defer wg.Done()

					var crc32 string = DataSignerCrc32(fmt.Sprint(i) + data)

					mu.Lock()
					res[i] = crc32
					mu.Unlock()
				}(wsg, i)
			}
			wsg.Wait()

			// send results
			var resString string
			for i := 0; i < mhCalculationsNum; i++ {
				resString += res[i]
			}
			out <- resString
		}(wg, i, fmt.Sprint(input.(string)), out)

		i++
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	data := make([]string, 0)

	for input := range in {
		data = append(data, input.(string))
	}
	sort.Strings(data)

	var res string = data[0]
	for i := 1; i < len(data); i++ {
		res += ("_" + data[i])
	}

	out <- res
}
