package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/gnames/gnfmt"
	"github.com/gnames/gnparser"
	"github.com/gnames/gnparser/ent/nameidx"
	"github.com/gnames/gnparser/ent/parsed"
)

func getNames(
	ctx context.Context,
	f io.Reader,
) <-chan nameidx.NameIdx {
	chIn := make(chan nameidx.NameIdx)
	sc := bufio.NewScanner(f)

	go func() {
		defer close(chIn)
		var count int
		for sc.Scan() {
			nameString := sc.Text()
			select {
			case <-ctx.Done():
				return
			case chIn <- nameidx.NameIdx{Index: count, NameString: nameString}:
			}
			count++
		}
	}()
	if err := sc.Err(); err != nil {
		log.Panic(err)
	}
	return chIn
}

func parseStream(
	gnp gnparser.GNparser,
	f io.Reader,
	quiet bool,
) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	chIn := getNames(ctx, f)
	chOut := make(chan parsed.Parsed)
	var wg sync.WaitGroup
	wg.Add(1)

	if gnp.Format() == gnfmt.CSV {
		parsed.HeaderCSV()
	}

	go gnp.ParseNameStream(ctx, chIn, chOut)

	// process parsing results
	go func() {
		defer cancel()
		defer wg.Done()
		start := time.Now()
		if gnp.Format() == gnfmt.CSV {
			fmt.Println(parsed.HeaderCSV())
		}
		var count int
		for {
			count++
			if count%50_000 == 0 && !quiet {
				progressLog(start, count)
			}
			select {
			case <-ctx.Done():
				return
			case v, ok := <-chOut:
				if !ok {
					return
				}
				fmt.Println(v.Output(gnp.Format()))
			}
		}
	}()
	wg.Wait()
}
