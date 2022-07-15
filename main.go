package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hmerritt/go-file-search/internal/bytebufferpool"
	"github.com/hmerritt/go-file-search/version"
	ngram "github.com/hmerritt/go-ngram"
)

func main() {
	version.PrintTitle()

	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatalln("No arguments provided")
	}

	searchFor := args[0]
	searchIn := "./"

	if len(args) > 1 {
		searchIn = args[1]
	}

	// Initialize the ngrams
	ng := ngram.NgramIndex{
		NgramMap:   make(map[string]map[int]*ngram.IndexValue),
		IndexesMap: make(map[int]*ngram.IndexValue),
		Ngram:      3,
	}

	// Walk the directory recursively
	counter := 1
	timeWalking := time.Now()
	filepath.WalkDir(searchIn, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			ng.Add(strings.ToLower(filepath.Base(s)), ngram.NewIndexValue(counter, d))
		}

		counter++
		return nil
	})

	fmt.Printf("Walkin took %d ms (%d items)\n", time.Since(timeWalking).Milliseconds(), len(ng.IndexesMap))
	timeSearching := time.Now()

	// Search through the ngram
	search := ng.Search(strings.ToLower(searchFor))

	// Print the results
	fmt.Println("Search took", time.Since(timeSearching))
	fmt.Println("\nFound", len(search), "results")
	for _, v := range search {
		data := v.Data.(fs.DirEntry) // Interface{} -> fs.DirEntry
		log_file(data)
	}
}

func log_file(f fs.DirEntry) {
	info, _ := f.Info()

	// Get new buffer
	buf := bytebufferpool.Get()

	// Format log to buffer
	_, _ = buf.WriteString(fmt.Sprintf("%60s | %7d\n",
		f.Name(),
		info.Size(),
	))

	// _, _ = buf.WriteString(fmt.Sprintf("%s |%s %3d %s| %7v | %15s |%s %-7s %s| %-"+errPaddingStr+"s %s\n",
	// 	timestamp.Load().(string),
	// 	statusColor(c.Response().StatusCode()), c.Response().StatusCode(), cReset,
	// 	stop.Sub(start).Round(time.Millisecond),
	// 	c.IP(),
	// 	methodColor(c.Method()), c.Method(), cReset,
	// 	c.Path(),
	// 	formatErr,
	// ))

	// Write buffer to output
	_, _ = os.Stdout.Write(buf.Bytes())

	// Put buffer back to pool
	bytebufferpool.Put(buf)
}
