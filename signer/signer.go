package main

import (
	"slices"
	"strconv"
	"strings"
	"sync"
)

var (
	crc32Cache = sync.Map{}
	md5Mu      = &sync.Mutex{}
)

func crc32WithCache(data string) string {
	if crc32, ok := crc32Cache.Load(data); ok {
		return crc32.(string)
	} else {
		crc32 = DataSignerCrc32(data)
		crc32Cache.Store(data, crc32)

		return crc32.(string)
	}
}

func md5WithLock(data string) string {
	md5Mu.Lock()
	defer md5Mu.Unlock()

	return DataSignerMd5(data)
}

func SingleHash(in, out chan any) {
	defer close(out)

	wg := &sync.WaitGroup{}

	for rawData := range in {
		wg.Add(1)
		go func(rawData any, out chan any, wg *sync.WaitGroup) {
			defer wg.Done()

			dataAsInt, ok := rawData.(int)
			if !ok {
				panic("expected int data")
			}

			data := strconv.FormatInt(int64(dataAsInt), 10)
			crc32 := crc32WithCache(data)
			md5 := md5WithLock(data)
			crc32Md5 := crc32WithCache(md5)

			result := crc32 + "~" + crc32Md5

			out <- result
		}(rawData, out, wg)
	}

	wg.Wait()
}

func MultiHash(in, out chan any) {
	defer close(out)

	wgOuter := &sync.WaitGroup{}

	for rawData := range in {
		wgOuter.Add(1)

		go func(out chan any, wgOuter *sync.WaitGroup) {
			defer wgOuter.Done()

			data, ok := rawData.(string)
			if !ok {
				panic("expected string data")
			}

			resultMap := make(map[int64]string, 6)
			wg := &sync.WaitGroup{}

			for th := range int64(6) {
				wg.Add(1)

				go func(th int64, data string, resultMap map[int64]string, wg *sync.WaitGroup) {
					defer wg.Done()

					formattedTh := strconv.FormatInt(th, 10)
					resultMap[th] = crc32WithCache(formattedTh + data)
				}(th, data, resultMap, wg)
			}

			wg.Wait()

			resultBuilder := strings.Builder{}
			for th := range int64(6) {
				resultBuilder.WriteString(resultMap[th])
			}

			out <- resultBuilder.String()
		}(out, wgOuter)
	}

	wgOuter.Wait()
}

func CombineResults(in, out chan any) {
	defer close(out)

	var allData []string

	for rawData := range in {
		data, ok := rawData.(string)
		if !ok {
			panic("expected string data")
		}

		allData = append(allData, data)
	}

	slices.Sort(allData)
	result := strings.Join(allData, "_")

	out <- result
}

func ExecutePipeline(jobs ...Job) {
	var prevOut chan any
	wg := &sync.WaitGroup{}

	for i, job := range jobs {
		var out chan any
		if i < len(jobs)-1 {
			out = make(chan any)
		}

		wg.Add(1)

		go func(out, prevOut chan any, wg *sync.WaitGroup) {
			defer wg.Done()

			if prevOut == nil {
				job(nil, out)
			} else {
				job(prevOut, out)
			}
		}(out, prevOut, wg)

		prevOut = out
	}

	wg.Wait()
}
