package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"hash"
	"os"
)

type options struct {
	File      string `long:"file" description:"Path to file containing names" required:"true"`
	Tickets   uint   `long:"numbilets" description:"Total number of available tickets" required:"true"`
	Parameter int    `long:"parameter" description:"Ticket number generation parameter" required:"true"`
}

type ticketNumberGenerator struct {
	parameter int64
	hasher    hash.Hash
	n         uint
}

func newTicketNumberGenerator(parameter int, n uint) *ticketNumberGenerator {
	return &ticketNumberGenerator{
		parameter: int64(parameter),
		hasher:    sha256.New(),
		n:         n,
	}
}

func (g *ticketNumberGenerator) generateTicketNumber(name string) uint64 {
	hasher := g.hasher
	hasher.Reset()
	hasher.Write([]byte(name))
	paramBuffer := new(bytes.Buffer)
	_ = binary.Write(paramBuffer, binary.BigEndian, g.parameter)
	hasher.Write(paramBuffer.Bytes())
	resultBuf := bytes.NewReader(hasher.Sum(nil)[:8])
	result := uint64(0)
	_ = binary.Read(resultBuf, binary.BigEndian, &result)
	return result%uint64(g.n) + 1
}

func main() {
	opts := options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatalf("failed parsing command line options: %s", err)
	}
	nameFile, err := os.Open(opts.File)
	if err != nil {
		log.Fatalf("failed opening file containing names")
	}
	defer nameFile.Close()

	scanner := bufio.NewScanner(nameFile)
	generator := newTicketNumberGenerator(opts.Parameter, opts.Tickets)
	for scanner.Scan() {
		name := scanner.Text()
		fmt.Printf("%s: %d\n", name, generator.generateTicketNumber(name))
	}
	if err = scanner.Err(); err != nil {
		log.Fatalf("failed reading file containing names")
	}
}
