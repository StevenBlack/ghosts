package main

import (
	"testing"
)

func TestMultihostLines(t *testing.T) {
	// testing for splitting multiple Domains per line into individual lines
	hf := Hosts{}
	hf.Load("./test/hosts-multi")

	got := len(hf.Domains)
	wantmorethan := 1

	if got <= wantmorethan {
		t.Errorf("got %d domain, want more than %d", got, wantmorethan)
	}
}

func TestTLD(t *testing.T) {
	// testing for TLD Tallies
	hf := Hosts{}
	hf.Load("./test/hosts-multi")

	got := len(hf.TLDs)
	wantmorethan := 1

	if got <= wantmorethan {
		t.Errorf("got %d TLD, want more than %d", got, wantmorethan)
	}
}

func TestPreserveHeadder(t *testing.T) {
	// testing file header preservation.
	hf := Hosts{}
	hf.Load("./test/hosts-comments-embedded")

	got := len(hf.Header)
	want := 4

	if got != want {
		t.Errorf("got %d header lines, want %d", got, want)
	}
}

func TestEmbeddedComments(t *testing.T) {
	// testing hosts lines with embedded comments
	hf := Hosts{}
	hf.Load("./test/hosts-comments-embedded")

	got := len(hf.Domains)
	want := 4

	if got != want {
		t.Errorf("got %d Domains, want %d", got, want)
	}
}

func TestDuplicates(t *testing.T) {
	// testing hosts with duplicates
	hf := Hosts{}
	hf.Load("./test/hosts-duplicates")

	got := len(hf.Domains)
	want := 5
	dupesgot := len(hf.Duplicates)
	dupeswant := 3

	if got != want {
		t.Errorf("got %d Domains, want %d", got, want)
	}

	if dupesgot != dupeswant {
		t.Errorf("got %d duplicate Domains, want %d", dupesgot, dupeswant)
	}
}

func TestJustText(t *testing.T) {
	// testing a file with just text, no hosts
	hf := Hosts{}
	hf.Load("./test/hosts-text")

	got := len(hf.Domains)
	want := 0

	if got != want {
		t.Errorf("got %d domains, want %d", got, want)
	}
}

func TestUrl(t *testing.T) {
	// testing for splitting multiple Domains per line into individual lines
	hf := Hosts{}
	hf.Load("https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts")

	got := len(hf.Domains)
	wantmorethan := 1

	if got <= wantmorethan {
		t.Errorf("got %d domain, want more than %d", got, wantmorethan)
	}
}

func TestUrlJustText(t *testing.T) {
	// testing a file with just text, no hosts
	hf := Hosts{}
	hf.Load("https://news.ycombinator.com/")

	got := len(hf.Domains)
	want := 0

	if got != want {
		t.Errorf("got %d domains, want %d", got, want)
	}
}

func TestSorting(t *testing.T) {
	// testing hosts with duplicates
	hf := Hosts{}
	a := "aa.ca"
	b := "zz.aa"

	got := hf.Normalize(a) < hf.Normalize(b)
	want := true

	if got != want {
		t.Errorf("aa.ca < zz.aa")
	}

	a = "cc.ca"
	b = "aa.cc.ca"

	got = hf.Normalize(a) < hf.Normalize(b)
	want = true

	if got != want {
		t.Errorf("cc.ca < aa.cc.aa")
	}

}
