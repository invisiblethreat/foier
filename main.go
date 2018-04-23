package main

import (
	"github.com/invisiblethreat/foier/scrape"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// CLIOptions are to set things via CLI
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

	var cli CLIOptions
	pflag.BoolVarP(&cli.Debug, "debug", "d", false, "Use debug mode")
	pflag.StringVarP(&cli.URI, "uri", "u",
		"https://foipop.novascotia.ca/foia/views/_AttachmentDownload.jsp?attachmentRSN=",
		"URI to use for fetching")
	pflag.IntVarP(&cli.Start, "start", "s", 0, "Starting ID")
	pflag.IntVarP(&cli.End, "end", "e", 7000, "End ID")
	pflag.IntVarP(&cli.Step, "increment", "i", 1, "Increment between steps")
	pflag.StringVar(&cli.Suffix, "suffix", "", "Suffix of the file you're downloading")
	pflag.IntVarP(&cli.Workers, "workers", "w", 5, "How many workers to use")

	pflag.Parse()

	if cli.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Debugf("Logging level: %s", logrus.GetLevel().String())

	config := foier.Config{
		URI:     cli.URI,
		Start:   cli.Start,
		End:     cli.End,
		Step:    cli.Step,
		Workers: cli.Workers,
		Suffix:  cli.Suffix,
	}

	foier.Scrape(config)
}
