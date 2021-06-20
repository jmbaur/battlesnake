package game

import "testing"

func TestCanMoveDirections(t *testing.T) {
	tt := []struct {
		name      string
		state     *state
		direction string
		expect    bool
	}{
		{
			name:      "Empty board",
			state:     &state{Board: board{Height: 5, Width: 5}},
			direction: Up,
			expect:    true,
		},
		{
			name:      "0x0 sized board",
			state:     &state{Board: board{Height: 0, Width: 0}},
			direction: Down,
			expect:    false,
		},
		{
			name:      "Head is on the left wall",
			state:     &state{Board: board{Height: 5, Width: 5}, You: snake{Head: coordinate{X: 0, Y: 0}}},
			direction: Left,
			expect:    false,
		},
		{
			name:      "Head is under a segment of another snake",
			state:     &state{Board: board{Height: 5, Width: 5, Snakes: []snake{{Body: []coordinate{{X: 1, Y: 1}}}}}, You: snake{Head: coordinate{X: 1, Y: 0}}},
			direction: Up,
			expect:    false,
		},
		{
			name:      "Head-to-head with another snake",
			state:     &state{Board: board{Height: 5, Width: 5, Snakes: []snake{{Health: 1, Head: coordinate{X: 0, Y: 1}}}}, You: snake{Health: 2, Head: coordinate{X: 0, Y: 0}}},
			direction: Up,
			expect:    true,
		},
	}

	for _, tc := range tt {
		if got := canMoveDirection(tc.state, tc.direction); got != tc.expect {
			t.Errorf("Test %s: got: %t, expected %t\n", tc.name, got, tc.expect)
		}
	}
}

func TestClosestFood(t *testing.T) {
	tt := []struct {
		name   string
		input  *state
		expect *coordinate
	}{
		{
			name:   "One food available",
			input:  &state{Board: board{Height: 5, Width: 5, Food: []coordinate{{X: 1, Y: 1}}}, You: snake{Head: coordinate{X: 0, Y: 0}}},
			expect: &coordinate{X: 1, Y: 1},
		},
		{
			name:   "Multiple food available",
			input:  &state{Board: board{Height: 5, Width: 5, Food: []coordinate{{X: 1, Y: 1}, {X: 4, Y: 1}}}, You: snake{Head: coordinate{X: 4, Y: 2}}},
			expect: &coordinate{X: 4, Y: 1},
		},
	}

	for _, tc := range tt {
		coord, _ := closestFood(tc.input)
		if *coord != *tc.expect {
			t.Errorf("Test %s: got %+v, expected %+v\n", tc.name, *coord, *tc.expect)
		}
	}
}

func TestDesiredDirection(t *testing.T) {
	tt := []struct {
		name   string
		a      coordinate
		b      coordinate
		expect string
	}{}

	for _, tc := range tt {
		if got := desiredDirection(tc.a, tc.b); got != tc.expect {
			t.Errorf("Test %s: got %+v, expected %+v\n", tc.name, got, tc.expect)
		}
	}
}
