package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type mMarkovChainPair struct {
	first  string
	second string
}

type mMarkovChainTuple struct {
	name  string
	count uint
	prob  float64
}

type mMarkovChain struct {
	conn  map[string][]mMarkovChainTuple
	sums  map[string]uint
	pairs map[mMarkovChainPair]bool
}

func (g *mMarkovChain) randomName() string {
	var random = rand.Intn(len(g.conn))
	i := 0
	var randomname string

	for k := range g.conn {
		if i == random {
			randomname = k
			break
		}
		i++
	}

	return randomname
}

func (g *mMarkovChain) randomconnection(count int) {
	for i := 0; i < count; i++ {
		ra := g.randomName()
		rb := g.randomName()

		g.add(ra, rb)
	}
}

func (g *mMarkovChain) add(a, b string) {
	if !g.pairs[mMarkovChainPair{a, b}] {
		g.pairs[mMarkovChainPair{a, b}] = true
		g.conn[a] = append(g.conn[a], mMarkovChainTuple{b, 0, 0})
	}
	g.sums[a]++

	totalcount := g.sums[a]
	list := g.conn[a]
	var element *mMarkovChainTuple

	for i := 0; i < len(list); i++ {
		element = &list[i]

		if element.name == b {
			element.count++
		}

		element.prob = float64(element.count) / float64(totalcount)
	}
}

func (g *mMarkovChain) random(a string) (string, bool) {
	randomfloat := rand.Float64()
	totalcount := g.sums[a]
	randint := uint(randomfloat*float64(totalcount) + 0.5)
	list := g.conn[a]
	j := uint(0)

	for _, tuple := range list {
		tcount := tuple.count

		j += tcount

		if j >= randint {
			return tuple.name, g.sums[tuple.name] != 0 && !strings.HasSuffix(tuple.name, ".")
		}
	}

	lasttuple := list[len(list)-1]
	return lasttuple.name, g.sums[lasttuple.name] == 0 && !strings.HasSuffix(lasttuple.name, ".")
}

func (g *mMarkovChain) generate() string {
	sentence := "=>"
	str, end := g.random(sentence)
	sentence += " " + str

	for end {
		str, end = g.random(str)
		sentence += " " + str
	}

	return strings.TrimLeft(sentence, " ")
}

func (g *mMarkovChain) insert(longtext string) {
	list := strings.Split(longtext, " ")

	for i, current := range list {
		if i == 0 {
			continue
		}
		g.add("=>", list[i-1])
		g.add(list[i-1], current)
	}
}

func newMarkovChain() *mMarkovChain {
	g := new(mMarkovChain)
	g.conn = make(map[string][]mMarkovChainTuple)
	g.sums = make(map[string]uint)
	g.pairs = make(map[mMarkovChainPair]bool)

	rand.Seed(time.Now().UnixNano())
	return g
}

func readLargeText(path string, optimizer func(string) string) string {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		panic("Could not find file")
	}

	var largetext, line string
	var scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		line = scanner.Text()

		line = optimizer(line)

		largetext += line
	}

	return largetext
}

func main() {
	largetext := readLargeText("bulk2.txt", func(line string) string {
		return strings.ToUpper(strings.Replace(line, ".", " ", -1) + ". ")
	})

	var markovGraph = newMarkovChain()

	markovGraph.insert(largetext)
	markovGraph.randomconnection(20000)

	for i := 0; i < 50; i++ {
		fmt.Println(markovGraph.generate())
		time.Sleep(time.Second)
	}

}
