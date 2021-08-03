package main

import (
	"bufio"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type Store struct {
	links []string
	path  string
}

// readPersistedData reads the previously seen links from the file at s.path.
func (s *Store) readPersistedData() error {
	file, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t := scanner.Text()
		lines = append(lines, t)
	}
	s.links = lines
	return scanner.Err()
}

// Exists checks to see if the provided URL exists in the datastore.
// Returns true if this is an existing url, the url is added to the
// datastore and returns false if it does not exist.
func (s *Store) Exists(url string) bool {
	for _, s := range s.links {
		if s == url {
			return true
		}
	}
	s.links = append(s.links, url)
	return false
}

// persistData overwrites the seen file with latest list of seen urls.
func (s *Store) persistData() error {
	log.Infof("Writing seen urls to %s", s.path)
	file, err := os.Create(s.path)
	if err != nil {
		log.Error(err)
		return err
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		log.Error(err)
		return err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		log.Error(err)
		return err
	}

	w := bufio.NewWriter(file)
	for _, line := range s.links {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

// Close persists the underlying data store to disk and before closing down.
func (s *Store) Close() {
	log.Infof("Closing seen store and saving seen urls to %s", s.path)
	if err := s.persistData(); err != nil {
		log.Error(err)
	}
}

// NewStore creates a new Store object with the data provided by
// the text file in the path parameter.
func NewStore(path string) (*Store, error) {
	store := &Store{path: path}
	if err := store.readPersistedData(); err != nil {
		return nil, err
	}
	return store, nil
}
