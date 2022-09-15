package analytic

import "fmt"

type WorkerInput []byte

type WorkerOutput struct {
	stats *ChatStatistics
	err   error
}

type WorkerPipe struct {
	input  chan<- WorkerInput
	output <-chan WorkerOutput
}

type Worker struct {
	threadCounts int
	workerPipes  <-chan WorkerPipe
}

func WorkerRunner(workerPipes chan<- WorkerPipe) {
	getData := make(chan WorkerInput)
	sendStats := make(chan WorkerOutput)
	workerPipe := WorkerPipe{
		input:  getData,
		output: sendStats,
	}
	workerPipes <- workerPipe

	for data := range getData {
		workerOutput := WorkerOutput{
			stats: nil,
			err:   nil,
		}
		data = RemoveBrokenLinesFromJsonFile(data)
		chat, err := ParseJson(data)
		if err != nil {
			workerOutput.err = fmt.Errorf("error while parse JSON file: %v\n", err)
			sendStats <- workerOutput
			continue
		}

		stats, err := GetChatStatistics(chat)
		if err != nil {
			workerOutput.err = fmt.Errorf("error while generation chat statistics: %v", err)
			sendStats <- workerOutput
			continue
		}

		workerOutput.stats = stats
		sendStats <- workerOutput
		workerPipes <- workerPipe
	}
}

func NewAnalyticWorker(threadCounts int) *Worker {
	workerPipes := make(chan WorkerPipe, threadCounts)
	for i := 0; i < threadCounts; i++ {
		go WorkerRunner(workerPipes)
	}

	return &Worker{
		threadCounts: threadCounts,
		workerPipes:  workerPipes,
	}
}

func (w *Worker) AnalyseJson(data []byte) (*ChatStatistics, error) {
	workerPipe := <-w.workerPipes
	workerPipe.input <- data
	output := <-workerPipe.output

	return output.stats, output.err
}
