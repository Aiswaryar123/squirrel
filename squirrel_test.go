package main

import "testing"

func TestPhi(t *testing.T) {
	c := Counts{n00: 1, n01: 1, n10: 1, n11: 1}
	if phi(c) != 0.0 {
		t.Errorf("Expected 0.0 but got %f\n", phi(c))
	}

	c = Counts{n00: 5, n01: 0, n10: 0, n11: 5}
	if phi(c) != 1.0 {
		t.Errorf("Expected 1.0 but got %f\n", phi(c))
	}

	c = Counts{n00: 0, n01: 5, n10: 5, n11: 0}
	if phi(c) != -1.0 {
		t.Errorf("Expected -1.0 but got %f\n", phi(c))
	}
}
