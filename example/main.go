package main

import (
	"fmt"
	"sync"

	"github.com/sergei-svistunov/local_ctx"
)

func ext_func() {
	data, err := local_ctx.Data()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", data.(string))
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			local_ctx.Call(fmt.Sprintf("goroutine %d", i), func() {
				ext_func()

				var wg1 sync.WaitGroup
				wg1.Add(1)

				local_ctx.Go(func() {
					ext_func()
					wg1.Done()
				})
				wg1.Wait()
			})
		}(i)
	}

	wg.Wait()
}
