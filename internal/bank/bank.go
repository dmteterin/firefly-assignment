package bank

import (
	"bufio"
	"encoding/json"
	"os"
	"sort"
	"sync"

	"github.com/dmteterin/firefly-assignment/internal/config"
	"github.com/rs/zerolog"
)

type WordBank struct {
	ValidBank *sync.Map
	WordCount []WordCount
	logger    zerolog.Logger
	done      chan struct{}
}

type WordCount struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

type Validator interface {
	Validate(words []string) *sync.Map
}

func New(cfg *config.Config, validator Validator, done chan struct{}, logger zerolog.Logger) (*WordBank, error) {
	file, err := os.Open(cfg.WordBankFilePath)
	if err != nil {
		return nil, err
	}

	fileScanner := bufio.NewScanner(file)

	fileScanner.Split(bufio.ScanLines)

	var words []string
	for fileScanner.Scan() {
		words = append(words, fileScanner.Text())
	}

	return &WordBank{
		ValidBank: validator.Validate(words),
		done:      done,
	}, nil
}

func (w *WordBank) CountMatches(tokens chan string) {
	for token := range tokens {
		if v, ok := w.ValidBank.Load(token); ok {
			w.ValidBank.Store(token, v.(int)+1)
		}
	}
	w.ValidBank.Range(func(key, value any) bool {
		if value != 0 {
			w.WordCount = append(w.WordCount, WordCount{Word: key.(string), Count: value.(int)})
		}
		return true
	})

	sort.Slice(w.WordCount, func(i, j int) bool {
		return w.WordCount[i].Count > w.WordCount[j].Count
	})

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err := enc.Encode(w.WordCount[0:10])
	if err != nil {
		w.logger.Error().Err(err).Msg("Could not encode word list")
	}

	w.done <- struct{}{}
}
