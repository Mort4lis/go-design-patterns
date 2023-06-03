package concurrency

// Job is a function type that is used to construct stages in Pipeline.
type Job func(in <-chan any, out chan<- any)

// Pipeline pattern is a series of stages connected by channels, where
// output channel of the stage is an input another stage. Each stage is a
// function of Job type, running at the separate goroutine. Pipeline returns
// two channels, where first one is for sending value to Pipeline, another one is
// for receiving a result through all the stages of Pipeline. If the sending channel
// is closed, it'll close other channels between stages, and Pipeline will be finished.
//
// By using a Pipeline, it separates the concerns of each stage,
// which provides numerous benefits such as:
//   - modify stages independent of one another.
//   - mix and match how stages are combined independently of modifying the stage.
func Pipeline(jobs ...Job) (chan<- any, <-chan any) {
	var out chan any

	in := make(chan any)
	resIn := in

	for _, j := range jobs {
		job := j
		out = make(chan any)

		go func(in <-chan any, out chan<- any) {
			job(in, out)
			close(out)
		}(in, out)

		in = out
	}

	return resIn, out
}
