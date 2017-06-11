package factpacks

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"io"

	log "github.com/Sirupsen/logrus"
)

var singularSep = regexp.MustCompile(`\s=>\s`)
var pluralSep = regexp.MustCompile(`\s->\s`)
var singularVerb = regexp.MustCompile(`\sis\s`)
var pluralVerb = regexp.MustCompile(`\sare\s`)

type Fact struct {
	name     string
	value    string
	isPlural bool
}

func (f *Fact) Output() string {
	var verb string
	if f.isPlural {
		verb = "are"
	} else {
		verb = "is"
	}

	return fmt.Sprintf("%s %s %s", f.name, verb, f.value)
}

type FactStore interface {
	LoadFactPack(filename string) error
	SetFact(fact *Fact)
	GetFact(name string) *Fact
	DeleteFact(name string)
	HumanFactSet(fact string)
}

type factStore struct {
	facts    map[string]*Fact
	factsMtx sync.RWMutex
}

func (fs *factStore) SetFact(fact *Fact) {
	fs.factsMtx.Lock()
	defer fs.factsMtx.Unlock()

	fs.facts[fact.name] = fact
}

func (fs *factStore) GetFact(name string) *Fact {
	fs.factsMtx.RLock()
	defer fs.factsMtx.RUnlock()

	if val, ok := fs.facts[name]; ok {
		return val
	}

	return nil
}

func (fs *factStore) DeleteFact(name string) {
	fs.factsMtx.Lock()
	defer fs.factsMtx.Unlock()

	delete(fs.facts, name)
}

func (fs *factStore) LoadFactPack(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	buf := bufio.NewReader(f)

	done := false
	for done != true {
		line, err := buf.ReadString('\n')
		if err != nil && err == io.EOF {
			done = true
		} else if err != nil {
			break
		}

		var parts []string
		var isPlural bool
		if singularSep.MatchString(line) {
			parts = singularSep.Split(line, 2)
			isPlural = false
		} else if pluralSep.MatchString(line) {
			parts = pluralSep.Split(line, 2)
			isPlural = true
		}

		if len(parts) != 2 {
			log.Debug("Invalid fact format. Skipping.")
			continue
		}

		name := strings.TrimSpace(parts[0])
		fact := strings.TrimSpace(parts[1])

		if name == "" || fact == "" {
			log.Debug("Fact name and details can't be empty. Skipping.")
			continue
		}

		fs.SetFact(&Fact{
			name:     name,
			value:    fact,
			isPlural: isPlural,
		})
	}

	return nil
}

func (fs *factStore) HumanFactSet(fact string) {
	var parts []string
	var isPlural bool
	if singularVerb.MatchString(fact) {
		parts = singularVerb.Split(fact, 2)
		isPlural = false
	} else if pluralVerb.MatchString(fact) {
		parts = pluralVerb.Split(fact, 2)
		isPlural = true
	}

	if len(parts) != 2 {
		log.Debug("There isn't enough information to parse a fact.")
		return
	}

	fs.SetFact(&Fact{
		name:     strings.TrimSpace(parts[0]),
		value:    strings.TrimSpace(parts[1]),
		isPlural: isPlural,
	})
}

func MakeFactStore() FactStore {
	return &factStore{
		facts: make(map[string]*Fact),
	}
}
