package main

import (
	"flag"
	"log"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	if from == "" || to == "" {
		log.Fatal("Flags --from and --to are required")
	}

	// Ensure offset and limit are non-negative
	if offset < 0 {
		offset = 0
	}
	if limit < 0 {
		limit = 0
	}

	res := Copy(from, to, offset, limit)
	if res != nil {
		log.Fatal(res)
	}
}
