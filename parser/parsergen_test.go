package parser

import "testing"

func TestEqualLr0Sets(t *testing.T) {

	set1 := []lr0Item{
		{1, 3},
		{1, 4},
		{1, 5},
		{2, 7},
		{2, 6},
		{2, 0},
	}

	set2 := []lr0Item{
		{2, 0},
		{1, 4},
		{2, 6},
		{1, 3},
		{2, 7},
		{1, 5},
	}

	set3 := []lr0Item{
		{2, 0},
		{2, 6},
		{2, 7},
		{1, 5},
		{1, 4},
		{1, 3},
	}

	set4 := []lr0Item{
		{2, 0},
		{2, 6},
		{2, 7},
		{1, 5},
		{1, 4},
		{1, 9},
	}

	set5 := []lr0Item{
		{0, 6},
		{0, 5},
		{0, 4},
		{0, 3},
		{0, 2},
		{0, 1},
	}

	set6 := []lr0Item{
		{0, 3},
		{0, 2},
		{0, 1},
		{0, 5},
		{0, 4},
		{0, 6},
	}

	if !checkEqualLr0ItemSets(set1, set2) {
		t.Fatalf("Two sets should be equal!")
	}

	if !checkEqualLr0ItemSets(set1, set3) {
		t.Fatalf("Two sets should be equal!")
	}

	if !checkEqualLr0ItemSets(set2, set3) {
		t.Fatalf("Two sets should be equal!")
	}

	if checkEqualLr0ItemSets(set1, set4) {
		t.Fatalf("Two sets shouldn't be equal!")
	}

	if checkEqualLr0ItemSets(set2, set4) {
		t.Fatalf("Two sets shouldn't be equal!")
	}

	if checkEqualLr0ItemSets(set3, set4) {
		t.Fatalf("Two sets shouldn't be equal!")
	}

	if !checkEqualLr0ItemSets(set5, set6) {
		t.Fatalf("Two sets shouldn't be equal!")
	}

}
