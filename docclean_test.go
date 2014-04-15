package main

import (
	"bytes"
	"strings"
	"testing"
)

const input = `
def foo():
    """
    """
    pass

def foo():
    '''
    '''
    pass

def foo():
    """
    Undocumented.
    """
    pass

def foo():
    """
    Undocumented. moar
    """

class Foo(Exception):
    """
    """`

const expected = `
def foo():
    pass

def foo():
    pass

def foo():
    pass

def foo():
    """
    Undocumented. moar
    """

class Foo(Exception):
    """
    """`

func TestDocClean(t *testing.T) {
	b := StripEmptyOrIrrelevantComments(bytes.NewBufferString(input))
	t.Logf("%v", b)
	out := strings.Join(b, "\n")
	if out != expected {
		t.Fatalf("expected: %v, got %v", expected, out)
	}
}
