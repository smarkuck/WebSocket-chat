package bot

import (
	"io"
	"time"
)

var Chain = newChain(2)

type bot struct {
	io.ReadCloser
	out io.Writer
}

func New() io.ReadWriteCloser {
	r, out := io.Pipe()
	return bot{r, out}
}

func (b bot) Write(buf []byte) (int, error) {
	go b.speak()
	return len(buf), nil
}

func (b bot) speak() {
	time.Sleep(time.Second)
	msg := Chain.generate(10)
	b.out.Write([]byte(msg))
}
