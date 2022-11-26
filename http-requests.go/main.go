package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func doRequest(i int, throttle <-chan time.Time, startChan chan string, resChan chan string, errChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	startChan <- fmt.Sprintf("%s [Worker %d] making request", fmt.Sprintf(time.Now().Format(time.RFC3339)), i)

	if resp, err := http.Get(fmt.Sprintf("http://numbersapi.com/%d", i)); err != nil {

		errChan <- fmt.Sprintf("%s [Worker %d] error:  %e", fmt.Sprintf(time.Now().Format(time.RFC3339)), i, err)
	} else {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			errChan <- fmt.Sprintf("%s [Worker %d] error:   %e", fmt.Sprintf(time.Now().Format(time.RFC3339)), i, err)
		} else {
			resChan <- fmt.Sprintf("%s [Worker %d] success:  \"%s\"", fmt.Sprintf(time.Now().Format(time.RFC3339)), i, body)
		}

	}

}

func main() {
	wg := sync.WaitGroup{}
	startChan := make(chan string)
	respChan := make(chan string)
	errChan := make(chan string)
	throttle := time.Tick(time.Second) // only make an HTTP request at most 1 per second
	limit := 50

	for i := 1; i <= limit; i++ {
		fmt.Printf("%s [Worker %d] starting\n", fmt.Sprintf(time.Now().Format(time.RFC3339)), i)
		wg.Add(1)
		go doRequest(i, throttle, startChan, respChan, errChan, &wg)

	}

	for finished := 0; finished < limit; {
		select {
		case start := <-startChan:
			fmt.Println(start)
		case resp := <-respChan:
			fmt.Println(resp)
			finished++
		case err := <-errChan:
			fmt.Println(err)
			finished++
		default:
		}
	}

	wg.Wait()

}
