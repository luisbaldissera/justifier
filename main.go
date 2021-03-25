package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
)

var params Params

type Params struct {
    Width       int
    MarginLeft  int
    MarginRight int
    PadLeft     int
    PadRight    int
    Align       int
}

const (
    AlignLeft   = iota
    AlignRight
    AlignCenter
    AlignJustify
)

func init() {
    flag.IntVar(&params.Width, "w", 78, "Width")
    flag.IntVar(&params.MarginLeft, "ml", 0, "Margin left")
    flag.IntVar(&params.MarginRight, "mr", 0, "Margin right")
    flag.IntVar(&params.PadLeft, "pl", 0, "Padding left")
    flag.IntVar(&params.PadRight, "pr", 0, "Padding right")
    flag.IntVar(&params.Align, "l", AlignLeft, "Align left")
    flag.IntVar(&params.Align, "r", AlignRight, "Align right")
    flag.IntVar(&params.Align, "c", AlignCenter, "Align centered")
    flag.IntVar(&params.Align, "j", AlignJustify, "Justify")
}

func main() {
    var (
        parChan = make(chan string)
        wordChan = make(chan string)
        lineChan = make(chan string)
        wg = sync.WaitGroup{}
    )
    flag.Parse()
    wg.Add(3)
    go paragrapher(&wg, parChan)
    go tokenizer(&wg, parChan, wordChan)
    go aligner(&wg, wordChan, lineChan)
    close(lineChan)
    wg.Wait()
}

func paragrapher(wg *sync.WaitGroup, par chan<- string) {
    defer wg.Done()
    s := bufio.NewScanner(os.Stdin)
    p := ""
    for s.Scan() {
        l := s.Text()
        if l == "" {
            par<- p
            p = ""
        } else {
            p = p + " " + l
        }
    }
    par<- p
    close(par)
}

func tokenizer(wg *sync.WaitGroup, par <-chan string, word chan<- string) {
    defer wg.Done()
    for p := range par {
        for _, w := range strings.Fields(p) {
            word<- w
        }
        word<- ""
    }
    close(word)
}

func aligner(wg *sync.WaitGroup, word <-chan string, line chan<- string) {
    defer wg.Done()
    c := 0
    for w := range word {
        if w == "" {
            fmt.Printf("\n\n")
            c = 0
        } else if c + len(w) <= params.Width {
            fmt.Printf("%s ", w)
            c = c + len(w) + 1
        } else {
            fmt.Printf("\n")
            c = 0
        }
    }
}

