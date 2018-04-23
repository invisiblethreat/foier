package foier

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// type can be structered the same as import, or they can be standalone
// statments- this format is editor friendly for code folding.
type (
	// Config sets the items needed for enumeration
	Config struct {
		URI     string
		Suffix  string
		Start   int
		End     int
		Step    int
		Workers int
	}

	// Target is a struct for passing information
	Target struct {
		URI  string
		Name string
	}

	// Document is a response to a request that will be saved
	Document struct {
		Name string
		Body io.Closer
	}
)

// Scrape farms out the work to all of the workers
func Scrape(config Config) {

	// initialize our channels
	targets := make(chan (Target))
	docs := make(chan (Document))

	// Generate all the targets in a goroutine. Pass the config through since we
	// reuse so many values.
	go genTargets(targets, config)

	// WaitGroup is key for not orphaning your results.
	var wg sync.WaitGroup

	// Spawn our workers
	for i := 0; i < config.Workers; i++ {
		wg.Add(1)
		// Workers live in the background, but with the WaitGroup we know when
		// they finish all of their work.
		go Worker(targets, docs, &wg)
	}

	// This call blocks further execution while the Workers are active. Without
	// this call, execution would rocket along, leaving orphaned data and
	// requests in the background. The mechanism is that a counter in the
	// WaitGroup must become zero before the call stops blocking.
	logrus.Debug("Waiting on the Workers to finish their jobs")
	wg.Wait()

	// Since all of the Workers have terminated we can now close the Document
	// channel, which will allow the range to terminate.
	logrus.Debug("Closing doc channel of type Document")
	close(docs)

	// There's technically a small race condition here, from when the last
	// Worker closes and when the result is processed by the Saver in the
	// background.

	time.Sleep(1 * time.Second)
}

// Worker does all of the scraping
func Worker(targets <-chan Target, docs chan<- Document, wg *sync.WaitGroup) {
	defer wg.Done()

	// Range over the channel until it is empty. This is blocking until the
	// channel is closed.
	for target := range targets {
		res, err := http.Get(target.URI)
		if err != nil {
			logrus.Debugf("Error fetching %s: %s", target, err.Error())
			continue
		}

		out, err := os.Create("data/" + target.Name)
		if err != nil {
			logrus.Errorf("Error creating file %s: %s", target.Name, err.Error())
			continue
		}
		size, err := io.Copy(out, res.Body)
		if err != nil {
			logrus.Errorf("Error writing file %s: %s", target.Name, err.Error())
			continue
		}
		logrus.Debugf("Wrote %d bytes for file %s", size, target.Name)
		res.Body.Close()
	}

}

// genTargets generates targets and sends to a channel
func genTargets(targets chan<- Target, config Config) {
	// < is a less expensive operation than <=,  we will add 1 to the condition
	for i := config.Start; i < config.End+1; i += config.Step {
		inc := strconv.Itoa(i)
		file := inc + config.Suffix
		fq := config.URI + file
		logrus.Debugf("Sending: %s", fq)
		targets <- Target{URI: fq, Name: file}
	}

	logrus.Debug("Closing the target channel of type Target")
	close(targets)
}
