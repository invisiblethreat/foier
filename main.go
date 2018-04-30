package main

import (
	"github.com/invisiblethreat/foier/scrape"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// CLIOptions are to set things via CLI.
type CLIOptions struct {
	Debug   bool
	URI     string
	Start   int
	End     int
	Step    int
	Suffix  string
	Workers int
}

func main() {

	// Declare the variable before we use it.
	var cli CLIOptions

	// Establish where we should hold the values from our CLI options.
	pflag.BoolVarP(&cli.Debug, "debug", "d", false, "Use debug mode")
	pflag.StringVarP(&cli.URI, "uri", "u", "", "URI to use for fetching")
	pflag.IntVarP(&cli.Start, "start", "s", 0, "Starting ID")
	pflag.IntVarP(&cli.End, "end", "e", 7000, "End ID")
	pflag.IntVarP(&cli.Step, "increment", "i", 1, "Increment between steps")
	pflag.StringVar(&cli.Suffix, "suffix", "x", "Suffix of the target files")
	pflag.IntVarP(&cli.Workers, "workers", "w", 5, "How many workers to use")

	pflag.Parse()

	// Are we going to use debug mode?
	if cli.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// Show the current log level- the default is logrus.InfoLevel.
	logrus.Infof("Logging level: %s", logrus.GetLevel().String())

	// Convert the CLI arguments into a configuration that Foier can understand.
	scrape := cli.mapToFoier()

	// We're all done with setup and now we're going to do the real work. All of
	// the items here could have been hardcoded inside the 'Scrape' method if
	// goal was the fewest lines of code possible.
	scrape.Run()
}

// mapToFoier transforms takes one type and returns another. This type of action
// is called a 'pointer receiver'. Since you're not attempting to change the
// values in the initial structure, you do not need to pass by reference.
// Self-modification of a struct requires that you pass a reference or your
// changes will not persist upon return.
func (c CLIOptions) mapToFoier() foier.Scrape {
	config := foier.Scrape{
		URI:     c.URI,
		Start:   c.Start,
		End:     c.End,
		Step:    c.Step,
		Workers: c.Workers,
		Suffix:  c.Suffix,
	}
	return config
}
