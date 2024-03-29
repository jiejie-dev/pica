package pica

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

//Benchmark holds all the parameters required to run the benchmark and respective methods
type Benchmark struct {
	//benchStart is the time at which the benchmark is started. This value is set when you run the benchmark, i.e. `Run`
	benchStart time.Time

	//TotalRequests is the total number of requests to be fired.
	TotalRequests uint64
	//BenchDuration is the duration over which all the requests are fired, in milliseconds.
	BenchDuration uint64
	//WaitPerReq is the duration to wait before firing the consecutive request
	WaitPerReq time.Duration
	//ShowProgress is set to true if it should display current stat when the benchmark is running/in progress.
	ShowProgress bool

	//waitGroup is used for tracking go routine executions
	waitGroup sync.WaitGroup

	//successRequestTimer is cummulative time taken for all successful responses, in nano seconds
	successRequestTimer uint64

	//errorRequestTimer is cummulative time taken for all error responses, in nano seconds
	errorRequestTimer uint64

	//minReqTime is the minimum time taken for a request, in nano seconds. (i.e. fastest)
	minReqTime uint64

	//successMinReqTime is the minimum time taken for a successful request, in nano seconds. (i.e. fastest)
	successMinReqTime uint64

	//errorMinReqTime is the minimum time taken for a failed request, in nano seconds. (i.e. fastest)
	errorMinReqTime uint64

	//maxReqTime is the maximum time taken for a request, in nano seconds. (i.e. slowest)
	maxReqTime uint64

	//successMaxReqTime is the maximum time taken for a successful request, in nano seconds. (i.e. slowest)
	successMaxReqTime uint64

	//errorMaxReqTime is the maximum time taken for a failed request, in nano seconds. (i.e. slowest)
	errorMaxReqTime uint64

	//requestCounter is the no.of requests completed (fail & success) at any given point of time
	requestCounter uint64
	//errorCounter is the no.of errors encountered at any given point of time
	errorCounter uint64
	//successCounter is the no.of errors encountered at any given point of time
	successCounter uint64

	//StatReqCount is the number of requests processed after which progress statistics is printed
	StatReqCount uint64

	//errors keeps track of all the errors encountered during the benchmark
	errors map[string]errorStat

	mutex *sync.Mutex
}

type errorStat struct {
	Message string
	Count   uint64
}

type errorsList map[string]errorStat

//errStat sets the error statistics
func (bA *Benchmark) errStat(err error) {
	errMsg := err.Error()

	//Creating a checksum for the error message;
	hash := md5.Sum([]byte(errMsg))
	//errKey is used as the key in map[string]errStat
	errKey := hex.EncodeToString(hash[:])

	bA.mutex.Lock()

	if item, ok := bA.errors[errKey]; ok {
		item.Count++
		bA.errors[errKey] = item
	} else {
		bA.errors[errKey] = errorStat{
			Message: errMsg,
			Count:   1,
		}
	}

	bA.errorCounter++

	bA.mutex.Unlock()
}

//New returns a new Benchmark pointer with the default values set
func New() *Benchmark {
	bA := &Benchmark{
		//200 requests per second. i.e. 200k per minute
		TotalRequests: 200,
		//1000 milliseconds. i.e. 1 second
		BenchDuration: 1000,
		//Stat would be printed after completing every 20 requests
		StatReqCount: 200 / 10,
		//By default, show progress is set to true
		ShowProgress: true,

		//cumulative time of requests is set to 0 initially
		successRequestTimer: 0,
		errorRequestTimer:   0,

		// Setting the highest value so that the first request would replace this
		minReqTime:        1 << 63,
		successMinReqTime: 1 << 63,
		errorMinReqTime:   1 << 63,

		errors: make(map[string]errorStat),

		mutex: &sync.Mutex{},
	}

	return bA
}

//Init initializes all the fields with their respective values
func (bA *Benchmark) Init() {
	//No.of go routines to be spwaned (i.e. concurrent requests)
	bA.waitGroup.Add(int(bA.TotalRequests))

	//Progress is always showed every time it completes 1/10th of the total number of requests
	bA.StatReqCount = bA.TotalRequests / 10

	//If wait time is not provided, calculate based on benchmark duration and total no.of requests
	if bA.WaitPerReq == 0 {
		bA.WaitPerReq = time.Nanosecond * time.Duration((bA.BenchDuration*1000000)/bA.TotalRequests)
	}

	//if 1/10th of total requests is less than 1, it's set to 1
	if bA.StatReqCount == 0 {
		bA.StatReqCount = 1
	}

}

//Done does all the operations & calculation required while ending a function call
func (bA *Benchmark) Done(doneTime time.Duration, err error) {
	timeConsumed := uint64(doneTime.Nanoseconds())

	atomic.AddUint64(&bA.requestCounter, 1)

	bA.PrintStat()

	if err != nil {
		bA.errStat(err)
		atomic.AddUint64(&bA.errorRequestTimer, timeConsumed)

		if timeConsumed > bA.errorMaxReqTime {
			atomic.StoreUint64(&bA.errorMaxReqTime, timeConsumed)
		}

		if timeConsumed < bA.errorMinReqTime {
			atomic.StoreUint64(&bA.errorMinReqTime, timeConsumed)
		}
	} else {
		atomic.AddUint64(&bA.successRequestTimer, timeConsumed)
		atomic.AddUint64(&bA.successCounter, 1)

		if timeConsumed > bA.successMaxReqTime {
			atomic.StoreUint64(&bA.successMaxReqTime, timeConsumed)
		}

		if timeConsumed < bA.successMinReqTime {
			atomic.StoreUint64(&bA.successMinReqTime, timeConsumed)
		}
	}

	bA.waitGroup.Done()
}

//Run runs the benchmark for the given function
func (bA *Benchmark) Run(fn func() error) {
	bA.benchStart = time.Now()
	fmt.Println(
		"\nDuration              :", time.Millisecond*time.Duration(bA.BenchDuration),
		"\nTotal requests        :", bA.TotalRequests,
		"\nWait time per request :", bA.WaitPerReq,
		"\nShow progess          :", bA.ShowProgress,
		", per", bA.StatReqCount, "request(s)",
		"\nStart                 :", bA.benchStart,
	)

	for i := uint64(0); i < bA.TotalRequests; i++ {
		go func() {
			startTime := time.Now()
			err := fn()
			bA.Done(time.Since(startTime), err)
		}()
		time.Sleep(bA.WaitPerReq)
	}

	bA.Final()
}

//PrintStat prints the stats available based on the given input params and the global variable values
func (bA *Benchmark) PrintStat() {
	if bA.ShowProgress == true && bA.requestCounter%bA.StatReqCount == 0 {
		fmt.Println(
			bA.requestCounter, " out of ", bA.TotalRequests, " done.",
			" Success:", bA.successCounter,
			" Errors:", bA.errorCounter)
	}
}

//Final prints all the statistics of the benchmark
//The input to Final() is the time when you start the benchmark
func (bA *Benchmark) Final() {
	//Wait for all the go routines to complete
	bA.waitGroup.Wait()

	successRatio := float64(bA.successCounter) * float64(100) / float64(bA.requestCounter)
	errorRatio := float64(bA.errorCounter) * float64(100) / float64(bA.requestCounter)

	//Req completion will be printed inside this infinite for loop, as well as the app would wait
	fmt.Println(
		"\n========================= Benchmark stats =========================\n",
		"\nDone               :", time.Now(),
		"\nTime to complete   :", time.Since(bA.benchStart),
		"\nTotal requests     :", bA.TotalRequests,
		"\nRequests completed :", bA.requestCounter,
		"\nSuccess            :", bA.successCounter, "("+strconv.FormatFloat(successRatio, 'f', -1, 64)+"%)",
		"\nErrors             :", bA.errorCounter, "("+strconv.FormatFloat(errorRatio, 'f', -1, 64)+"%)",
	)

	if bA.successCounter > 0 {
		fmt.Println(
			"\nAverage time per successful request :", time.Nanosecond*time.Duration(bA.successRequestTimer)/time.Duration(bA.successCounter),
			"\nFastest                             :", time.Duration(time.Nanosecond*time.Duration(bA.successMinReqTime)),
			"\nSlowest                             :", time.Duration(time.Nanosecond*time.Duration(bA.successMaxReqTime)),
		)
	}

	if bA.errorCounter > 0 {
		fmt.Println(
			"\nAverage time per failed request :", time.Nanosecond*time.Duration(bA.errorRequestTimer)/time.Duration(bA.errorCounter),
			"\nFastest                         :", time.Duration(time.Nanosecond*time.Duration(bA.errorMinReqTime)),
			"\nSlowest                         :", time.Duration(time.Nanosecond*time.Duration(bA.errorMaxReqTime)),
		)
	}

	if len(bA.errors) > 0 {
		println("\n\nError messages (" + strconv.Itoa(len(bA.errors)) + ")")
		idx := 1
		for _, item := range bA.errors {
			println("\n "+strconv.Itoa(idx)+".", item.Message)
			println("  Occurrences:", item.Count)
			idx++
		}
	}
}
