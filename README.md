# go-queue
A Multi-safe Delay queue by Golang.

#### How To Use

```go
type Task struct {
	ID int
}

// create a delay queue with cap 3
q := NewDelayQueue(3)
// add a Task to queue which should run after 10 seconds later
q.Offer(Task{ID: 1}, time.Now().Add(10*time.Second))

// run a goroutine to deal timer check
go q.Poll()

// get ready element from delay queue
ele := <-q.C
fmt.Println(ele) // Task{ID: 1}

// stop go q.Poll(), you can run queue by call q.Poll() then.
q.Stop()

```
