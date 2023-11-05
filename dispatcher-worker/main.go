package main

import (
	"fmt"

	"github.com/taxintt/go-playgrounds/dispatcher-worker/component"
)

var (
	MaxWorker = 2
	MaxQueue  = 5
)

func main() {
	dispatcher := component.NewDispatcher(MaxWorker)
	dispatcher.Run()

	for i := 0; i < MaxQueue; i++ {
		dispatcher.Add(component.Payload{Message: "Hello World! " + fmt.Sprint(i) + ""})
	}

	dispatcher.Stop()
}
