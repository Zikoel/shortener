package shortener

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/teris-io/shortid"
)

// Finder is able to recover our URLs from key
type Finder interface {
	Lookup(string) (string, error)
	ReadCounter(string) (uint64, error)
}

// Persister is able to store our urls
type Persister interface {
	Persist(string, interface{}) error
	Delete(string) error
	IncrementCounter(string) (uint64, error)
}

// Shortener a shortener
type Shortener struct {
	finder    Finder
	persister Persister
	sid       *shortid.Shortid
}

// DeleteURLByKey delete the URL from from the corresponding key
func (s *Shortener) DeleteURLByKey(key string) error {
	err := s.persister.Delete(key)
	if err != nil {
		return err
	}
	err = s.persister.Delete(fmt.Sprintf("%s:count", key))
	return err
}

// CountVisit count a visit for the specified key
func (s *Shortener) CountVisit(key string) (uint64, error) {
	return s.persister.IncrementCounter(key)
}

// CollectStats collect stats for the specified key
func (s *Shortener) CollectStats(key string) (uint64, error) {
	val, err := s.finder.ReadCounter(fmt.Sprintf("%s:count", key))
	if err != nil {
		return 0, err
	}
	return val, nil
}

// URLFromKey implement the logic for a lookup on storage
func (s *Shortener) URLFromKey(key string) (string, error) {
	url, err := s.finder.Lookup(key)

	if err != nil {
		return "", err
	}

	return url, nil
}

func generateHashFromURL(URL string) string {
	hasher := md5.New()
	hasher.Write([]byte(URL))
	completeHash := hasher.Sum(nil)

	encodedHash := base64.URLEncoding.EncodeToString(completeHash)

	var shordedEncodedHash strings.Builder
	shordedEncodedHash.WriteString(encodedHash[:4])
	// We skip the last two chars that can be padding, the point here is to be repeatable
	shordedEncodedHash.WriteString(encodedHash[len(encodedHash)-6:len(encodedHash)-2])

	return shordedEncodedHash.String()
}

// KeyFromURL store the url with suggested key or store by a self generated key
func (s *Shortener) KeyFromURL(URL, suggestedKey string) (string, error) {

	if !strings.HasPrefix(URL,"http://") && !strings.HasPrefix(URL,"http://") {
		return "", errors.New("URL must have the protocol part")
	}

	var key string = suggestedKey

	if key == "" {
		key = generateHashFromURL(URL)
	}

	alreadyExistURL, err := s.finder.Lookup(key)

	if err == nil && URL == alreadyExistURL {
		// we already have that key and the URL on value is the same... nothing to do
		return key, nil
	}

	if err == nil && URL != alreadyExistURL {
		// uh we detect an hash collision, very unlikely but no problem, generate a random key
		key, err = s.sid.Generate()

		if err != nil {
			return "", err
		}

		return s.KeyFromURL(URL, string(key[:8]))
	}

	err = s.persister.Persist(key, URL)
	if err != nil {
		return "", err
	}

	err = s.persister.Persist(fmt.Sprintf("%s:count", key), 0)
	if err != nil {
		return "", err
	}

	return key, nil
}

// CreateShortener create a new initilizated shortener object
func CreateShortener(finder Finder, persister Persister, seed uint64) (*Shortener, error) {
	sid, err := shortid.New(1, shortid.DefaultABC, seed)
	if err != nil {
		return nil, err
	}

	return &Shortener{
		finder:    finder,
		persister: persister,
		sid:       sid,
	}, nil
}
