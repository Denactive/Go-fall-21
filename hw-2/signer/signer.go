package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

const pipeline_DEBUG = false
const single_hash_DEBUG = false
const multi_hash_DEBUG = false
const combine_results_DEBUG = false

const single_hash_delay = 20 * time.Millisecond

func ExecutePipeline(jobs ...job) {
	wg := new(sync.WaitGroup)
	var prev_chan chan interface{} = nil

	for i, job_i := range jobs {
		cur_chan := make(chan interface{}, 1)

		wg.Add(1)
		go func(wg *sync.WaitGroup, job_i job, i int, in, out chan interface{}) {
			defer wg.Done()

			if pipeline_DEBUG {
				fmt.Println("#", i, job_i, "started. input", in, ", output", out)
			}

			job_i(in, out)
			close(out)

			if pipeline_DEBUG {
				fmt.Println("#", i, job_i, "finished. Channel", out, "closed")
			}
		}(wg, job_i, i, prev_chan, cur_chan)

		prev_chan = cur_chan
	}
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	wg := new(sync.WaitGroup)
	i := 0

	for input := range in {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int, data string) {
			defer wg.Done()

			if single_hash_DEBUG {
				fmt.Println("#", i, "SingleHash data", data)
			}

			var md5 string = DataSignerMd5(data)

			if single_hash_DEBUG {
				fmt.Println("#", i, "SingleHash md5(data)", md5)
			}

			// wait subgroup
			wsg := new(sync.WaitGroup)
			wsg.Add(2)
			crc32_data_ch := make(chan string, MaxInputDataLen)
			crc32_md5_ch := make(chan string, MaxInputDataLen)

			go func(wsg *sync.WaitGroup, i int, out chan<- string, data string) {
				defer wsg.Done()
				var crc32 string = DataSignerCrc32(data)
				if single_hash_DEBUG {
					fmt.Println("#", i, "SingleHash crc32(data)", crc32)
				}
				out <- crc32
				close(out)
			}(wsg, i, crc32_data_ch, data)

			go func(wsg *sync.WaitGroup, i int, out chan<- string, md5 string) {
				defer wsg.Done()
				var crc32_md5 string = DataSignerCrc32(md5)
				if single_hash_DEBUG {
					fmt.Println("#", i, "SingleHash crc32(md5(data))", crc32_md5)
				}
				out <- crc32_md5
				close(out)
			}(wsg, i, crc32_md5_ch, md5)

			wsg.Wait()

			var res string = <-crc32_data_ch + "~" + <-crc32_md5_ch

			if single_hash_DEBUG {
				fmt.Println("#", i, "SingleHash result", res)
			}

			out <- res
		}(wg, i, fmt.Sprint(input.(int)))
		i++
		time.Sleep(single_hash_delay)
	}
	wg.Wait()
}

const multi_hash_num = 6

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
			res := make(map[int]string, multi_hash_num)

			// MultiHash 1/6 subworker
			for i := 0; i < multi_hash_num; i++ {
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
			var res_string string
			for i := 0; i < multi_hash_num; i++ {
				if multi_hash_DEBUG {
					fmt.Println("#", i, input, "MultiHash: crc32(th+SingleHash))", i, res[i])
				}
				res_string += res[i]
			}
			if multi_hash_DEBUG {
				fmt.Println(input, "MultiHash result:", res_string)
			}
			out <- res_string
		}(wg, i, fmt.Sprint(input.(string)), out)

		i++
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	data := make([]string, 0)

	for input := range in {
		if combine_results_DEBUG {
			fmt.Println("\tCombineResults got", input)
		}
		data = append(data, input.(string))
	}
	sort.Strings(data)
	if combine_results_DEBUG {
		fmt.Println("\tsorted:", data)
	}
	var res string = data[0]
	for i := 1; i < len(data); i++ {
		res += ("_" + data[i])
	}

	out <- res
}
