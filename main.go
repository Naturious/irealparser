package main

import (
	"fmt"
	"github.com/naturious/irealparser/internal/songbook"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: irealparser <inputfile>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	book := songbook.ParseIRealInput(string(data))

	fmt.Printf("Book: %s\n", book.Name)
	for _, song := range book.Songs {
		fmt.Printf("Title: %s\n", song.Title)
		fmt.Printf("Composer: %s\n", song.Composer)
		fmt.Printf("Style: %s\n", song.Style)
		fmt.Printf("Key: %s\n", song.Key)
		fmt.Printf("TimeSig: %s\n", song.TimeSig)
		for i, m := range song.Music {
			fmt.Printf("  Measure %d: %+v\n", i+1, m)
		}
		fmt.Println()
	}
}
