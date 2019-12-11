package main

import (
	"testing"
)

func TestMultihostLines(t *testing.T) {
	// testing for splitting multiple domains per line into individual lines
	hf := Hosts{}
	hf.load("./test/data/hosts-multi")

	got := len(hf.domains)
	wantmorethan := 1

	if got <= wantmorethan {
		t.Errorf("got %d domain, want more than %d", got, wantmorethan)
	}
}

func TestEmbeddedComments(t *testing.T) {
	// testing hosts lines with embedded comments
	hf := Hosts{}
	hf.load("./test/data/hosts-comments-embedded")

	got := len(hf.domains)
	want := 4

	if got != want {
		t.Errorf("got %d domains, want %d", got, want)
	}
}

func TestDuplicates(t *testing.T) {
	// testing hosts with duplicates
	hf := Hosts{}
	hf.load("./test/data/hosts-duplicates")

	got := len(hf.domains)
	want := 5
	dupesgot := len(hf.duplicates)
	dupeswant := 1

	if got != want {
		t.Errorf("got %d domains, want %d", got, want)
	}

	if dupesgot != dupeswant {
		t.Errorf("got %d duplicate domains, want %d", dupesgot, dupeswant)
	}
}

func TestJustText(t *testing.T) {
	// testing a file with just text, no hosts
	hf := Hosts{}
	hf.load("./test/data/hosts-text")

	got := len(hf.domains)
	want := 0

	if got != want {
		t.Errorf("got %d domains, want %d", got, want)
	}
}

func TestUrl(t *testing.T) {
	// testing for splitting multiple domains per line into individual lines
	hf := Hosts{}
	hf.load("https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts")

	got := len(hf.domains)
	wantmorethan := 1

	if got <= wantmorethan {
		t.Errorf("got %d domain, want more than %d", got, wantmorethan)
	}
}

func TestUrlJustText(t *testing.T) {
	// testing a file with just text, no hosts
	hf := Hosts{}
	hf.load("https://news.ycombinator.com/")

	got := len(hf.domains)
	want := 0

	if got != want {
		t.Errorf("got %d domains, want %d", got, want)
	}
}
