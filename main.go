package main

import (
	"flag"
	"github.com/h2non/bimg"
	"io/ioutil"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	iterations, concurrency int
)

func init() {
	flag.IntVar(&iterations, "n", 1000000, "number of iterations")
	flag.IntVar(&concurrency, "c", 20, "concurrency")
}

func main() {
	flag.Parse()

	f, err := ioutil.ReadFile("/test.jpg")
	if err != nil {
		log.Fatalf("failed to read file %v", err)
	}

	size, err := bimg.NewImage(f).Size()
	if err != nil {
		log.Fatalf("failed to get image sie %v", err)
	}

	ch := make(chan struct{}, concurrency)

	for i := 0; i < concurrency; i++ {
		ch <- struct{}{}
	}

	var wg sync.WaitGroup
	wg.Add(iterations)

	start := time.Now()
	var counter int

	defer func() {
		elapsed := time.Now().Sub(start)

		log.Printf("finished %d iterations in %v", counter, elapsed)
	}()

	for i := 0; i < iterations; i++ {

		<-ch

		counter++

		if i > 0 && i%10000 == 0 {
			log.Printf("%d iterations\n", i)
		}

		go func() {
			defer func() {
				ch <- struct{}{}
				wg.Done()
			}()

			w := rand.Intn(size.Width-10) + 10
			h := rand.Intn(size.Height-10) + 10

			o := bimg.Options{
				Quality: 85,
				Width:   w,
				Height:  h,
			}

			_, err := bimg.Resize(f, o)
			if err != nil {
				log.Fatalf("failed to resize; width: %d; height: %d; error: %v", w, h, err)
			}
		}()
	}

	wg.Wait()
}
