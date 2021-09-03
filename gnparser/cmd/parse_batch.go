package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/gnames/gnfmt"
	"github.com/gnames/gnparser"
	"github.com/gnames/gnparser/ent/parsed"
)

func parseBatch(
	gnp gnparser.GNparser,
	f io.Reader,
	quiet bool,
) {
	batch := make([]string, batchSize)
	chOut := make(chan []parsed.Parsed)
	start := time.Now()
	var wg sync.WaitGroup

	wg.Add(1)
	go processResults(chOut, &wg, gnp.Format())

	sc := bufio.NewScanner(f)
	var i, count int
	for sc.Scan() {
		batch[count] = sc.Text()
		count++
		if count == batchSize {
			i++
			if !quiet {
				progressLog(start, count*i)
			}
			chOut <- gnp.ParseNames(batch)
			batch = make([]string, batchSize)
			count = 0
		}
	}
	chOut <- gnp.ParseNames(batch[:count])
	close(chOut)
	if err := sc.Err(); err != nil {
		log.Panic(err)
	}
	wg.Wait()
}

func processResults(
	out <-chan []parsed.Parsed,
	wg *sync.WaitGroup,
	f gnfmt.Format,
) {
	defer wg.Done()

	header := parsed.HeaderCSV(f)
	if header != "" {
		fmt.Println(header)
	}

	for pr := range out {
		for i := range pr {
			fmt.Println(pr[i].Output(f))
		}
	}
}
