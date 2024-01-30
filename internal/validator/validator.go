package validator

import (
	"regexp"
	"sync"

	"github.com/dmteterin/firefly-assignment/internal/config"
)

type Validator struct {
	minWordLength int
	validationRE  string
}

func New(config *config.Config) *Validator {
	return &Validator{
		minWordLength: config.ValidatorMinWordLength,
		validationRE:  config.ValidatorRegExp,
	}
}

func (v *Validator) Validate(words []string) *sync.Map {
	validWords := sync.Map{}
	isValidExpression := regexp.MustCompile(v.validationRE).MatchString
	for _, word := range words {
		if len(word) < v.minWordLength {
			continue
		}
		if !isValidExpression(word) {
			continue
		}
		validWords.Store(word, 0)
	}
	return &validWords
}
