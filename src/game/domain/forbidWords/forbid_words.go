package forbidWords

import (
	mgo "gopkg.in/mgo.v2"
	//"strings"
	"util"
)

const (
	wordFilterC = "word_filter"
)

type WordFilter struct {
	Word string `bson:"word"`
}

var forbidWords []string

func InitForbidWords() error {
	var words []*WordFilter
	err := util.WithGameCollection(wordFilterC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&words)
	})
	if err == nil {
		for _, word := range words {
			forbidWords = append(forbidWords, word.Word)
		}
	}
	return err
}

func IsForbid(s string) bool {
	return false
	/*for _, word := range forbidWords {
		if strings.Contains(s, word) {
			return true
		}
	}
	return false*/
}
