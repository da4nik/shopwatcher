package parsers

import (
	"fmt"
	"regexp"

	"github.com/da4nik/shopwatcher/parsers/wildberries"
	"github.com/da4nik/shopwatcher/types"
)

var parsers = make(map[string]types.Parser)

func init() {
	parsers[wildberries.ShopRegexp] = wildberries.Parse
}

// GetParser - returns parser by name
func GetParser(name string) (parser types.Parser, err error) {
	parser, ok := parsers[name]
	if !ok {
		return nil, fmt.Errorf("No such parser: %s", name)
	}

	return parser, nil
}

// Parse parses product by url
func Parse(url string) (types.Product, error) {
	var parseFunction types.Parser
	for rgex, prser := range parsers {
		if matched, _ := regexp.MatchString(rgex, url); matched {
			parseFunction = prser
			break
		}
	}
	return parseFunction(url)
}
