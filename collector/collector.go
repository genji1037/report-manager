package collector

import (
	"context"
	"fmt"
	"reflect"
	"report-manager/logger"
	"sync"
	"time"
)

func Collect(collectors []Collector) {
	for i := range collectors {
		typ := reflect.TypeOf(collectors[i])
		logger.Infof("%s.Collect() begin", typ)
		t0 := time.Now()
		err := collectors[i].Collect()
		if err != nil {
			logger.Errorf("%s.Collect() failed: %s", typ, err.Error())
		}
		logger.Infof("%s.Collect() cost %s", typ, time.Now().Sub(t0))
	}
}

func CollectParallel(collectors []Collector) {
	return
	wg := sync.WaitGroup{}

	ok := make([]chan struct{}, 0, len(collectors))
	for range collectors {
		ok = append(ok, make(chan struct{}))
	}

	for i := range collectors {
		wg.Add(1)
		ctx, _ := context.WithTimeout(context.Background(), time.Minute)
		go func(ctx context.Context, index int) {
			defer wg.Done()
			t0 := time.Now()
			select {
			case <-ctx.Done():
				fmt.Printf("%s.Collect() timeout\n", reflect.TypeOf(collectors[index]))
			case <-ok[i]:
			}
			fmt.Printf("%s.Collect() cost %s\n", reflect.TypeOf(collectors[index]), time.Now().Sub(t0))
		}(ctx, i)
		go func(index int) {
			err := collectors[index].Collect()
			if err != nil {
				fmt.Printf("%s.Collect() collect data failed: %s\n", reflect.TypeOf(collectors[index]), err.Error())
			}
			ok[i] <- struct{}{}
		}(i)
	}
	wg.Wait()
}
