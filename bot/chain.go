package bot

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"sync"
)

type prefix []string

func (p prefix) string() string {
	return strings.Join(p, " ")
}

func (p prefix) shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

type chain struct {
	suffixes  map[string][]string
	prefixLen int
	mtx       sync.Mutex
}

func newChain(prefixLen int) *chain {
	return &chain{make(map[string][]string), prefixLen, sync.Mutex{}}
}

func (c *chain) Write(b []byte) (int, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	br := bytes.NewReader(b)
	p := make(prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			left := int(br.Size()) - br.Len()
			if err == io.EOF {
				return left, nil
			}
			return left, err
		}
		key := p.string()
		c.suffixes[key] = append(c.suffixes[key], s)
		p.shift(s)
	}
}

func (c *chain) generate(n int) string {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	words := new(strings.Builder)
	p := make(prefix, c.prefixLen)
	for i := 0; i < n; i++ {
		choices := c.suffixes[p.string()]
		if len(choices) == 0 {
			break
		}
		next := choices[rand.Intn(len(choices))]
		fmt.Fprint(words, " ", next)
		p.shift(next)
	}
	fmt.Fprint(words, "\n")
	if words.Len() == 0 {
		return ""
	}
	return words.String()[1:]
}
