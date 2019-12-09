package main

import (
	"testing"
)

func TestMultihostLines(t *testing.T) {
	hf := Hosts{}
	hf.loadfile("./test/data/hosts-multi")

	got := len(hf.domains)
	wantmorethan := 1

	if got <= wantmorethan {
		t.Errorf("got %d domain, want more than %d", got, wantmorethan)
	}
}

func TestEmbeddedComments(t *testing.T) {
	hf := Hosts{}
	hf.loadfile("./test/data/hosts-comments-embedded")

	got := len(hf.domains)
	want := 4

	if got != want {
		t.Errorf("got %d domains, want %d", got, want)

	}
}
