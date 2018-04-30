package foier

// The package name does not need to match the directory that it is placed in.
// You can split a package across an arbitrary number of files as long as they
// are in the same directory. 'package' must also be on the first line of the
// file.

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"

	// This import looks like it is remote, but it exists as a path within your
	// $GOHOME, which defaults to ~/go if not set.
	"github.com/sirupsen/logrus"
)

// type can be structered the same as import, or they can be standalone
// statments- this format is editor-friendly for code folding.
type (
	// Scrape sets the items needed for enumeration. Since 'Config' is
	// capitalized, it can be accessed as a type outside of this package with
	// by importing the package and referenceing 'foier.Config'
	Scrape struct {
		// The member attributes are also public and can be set external to the
		// package
		URI     string
		Suffix  string
		Start   int
		End     int
		Step    int
		Workers int
	}

	// target is a struct for passing information about what you want to fetch
	// and save. It is lowercase and can not be accessed outside of this
	// package. Other files within the package can access it without issue.
	target struct {
		// The member attribues are lowercase which makes them private to the
		// package. This gives you the ability to change the internal workings
		// of the package/structure, and not worry about external impacts.
		uri  string
		name string
	}
)

// Run farms out the work to all of the workers
func (s Scrape) Run() {

	// initialize our channels
	targets := make(chan (target))

	// Generate all the targets in a goroutine. Pass the config through since we
	// reuse so many values. We use a goroutine so we don't block the rest of
	// execution.
	go genTargets(targets, s)

	// WaitGroup is key for not orphaning your results. In this case it is used
	// as a semaphore to track the state of goroutines that are stil running.
	// When the 'Wait()' method is called, the semaphore will block until the
	// count reaches zero, meaning that all goroutines have completed, which
	// prevents you from orphaning your results.
	var wg sync.WaitGroup

	// Spawn our workers
	for i := 0; i < s.Workers; i++ {

		// Add '1' to the WaitGroup, which is a blocking condition for the
		// 'Wait()' method.
		wg.Add(1)

		// Workers live in the background, but with the WaitGroup we know when
		// they finish all of their work. Passing 'wg' by reference, or the call
		// to 'Done()' would decrement a copy and we'd be stuck waiting forever.
		// Not using a WaitGroup would spawn off all of your workers to begin
		// execution, but not block and the function would retun almost
		// instantly and most of your results would be orphaned.
		go worker(targets, &wg)
	}

	// This call blocks further execution while the Workers are active. Without
	// this call, execution would rocket along, leaving orphaned data and
	// requests in the background. The mechanism is that a counter in the
	// WaitGroup must become zero before the call stops blocking.
	logrus.Debug("Waiting on the Workers to finish their jobs")
	wg.Wait()
}

// worker does all of the downloading. Notice that the WaitGroup is passed by
// reference so all workers are reporting back to the same semaphore.
func worker(targets <-chan target, wg *sync.WaitGroup) {

	// 'defer' is called after the function returns. 'Done()' decrements the
	// semaphore by one, getting closer to the exit condition for 'Wait()'
	defer wg.Done()

	// Range over the channel until it is empty. This call is blocking until the
	// channel is closed. That is, a channel must both be empty and closed
	// before this loop will terminate.
	for target := range targets {

		// Get the target document.
		res, err := http.Get(target.uri)

		// This is a code tick that you'll see repeatedly in Golang. Errors are
		// just values and you should handle them accordingly.
		if err != nil {
			// This is a silent failure
			logrus.Debugf("Error fetching %s: %s", target, err.Error())
			continue
		}

		// Make the file to contain the results of the document
		out, err := os.Create("data/" + target.name)
		// The tick again
		if err != nil {
			// We will see this at logrus.InfoLevel as this is probably going to
			// indicate that there is something wrong outside of this program.
			// It could also be argued to call 'panic' here, as this will likely
			// occur for each attempt of this call.
			logrus.Errorf("Error creating file %s: %s", target.name, err.Error())
			continue
		}

		// Copy the resulting request body to the file that you created.
		size, err := io.Copy(out, res.Body)
		// The tick again
		if err != nil {
			// We will see this at logrus.InfoLevel as this is probably going to
			// indicate that there is something wrong outside of this program.
			// As above, this could be a 'panic' situation and could reasonably
			// be assumed to occur on each iteration of the loop.
			logrus.Errorf("Error writing file %s: %s", target.name, err.Error())
			continue
		}

		// If we're debugging, let us know the size of the file.
		logrus.Debugf("Wrote %d bytes for file %s", size, target.name)

		// Close the response body. This can technically return an error, but if
		// it does, there's nothing that you can really do about it.
		res.Body.Close()
	}

	logrus.Debug("The channel of type Target is closed and empty- returning")
}

// genTargets generates targets and sends them to a channel
func genTargets(targets chan<- target, config Scrape) {

	// < is a less expensive operation than <=, we will add 1 to the condition.
	// <= is potentially optimized by the compiler, but this explicitly does
	// only one comparison.
	for i := config.Start; i < config.End+1; i += config.Step {

		// We must convert the 'int' to a string to build out our request or we
		// will receive a type mismatch error.
		inc := strconv.Itoa(i)

		// This is our file name
		file := inc + config.Suffix

		// This is the fully qualified path to the file
		fq := config.URI + file
		logrus.Debugf("Sending: %s", fq)

		// Send the target to the channel. Take note of the inline declaration
		// or the struc and the assignment of values.
		targets <- target{uri: fq, name: file}
	}

	logrus.Debug("Closing the target channel of type Target")
	close(targets)
}
