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
			name:      "Head is under a segment of another snake",
			state:     &state{Board: board{Height: 5, Width: 5, Snakes: []snake{{Body: []coordinate{{X: 1, Y: 1}}}}}, You: snake{Head: coordinate{X: 1, Y: 0}}},
			direction: Up,
			expect:    false,
		},
	}

	for _, tc := range tt {
		if got := canMoveDirection(tc.state, tc.direction); got != tc.expect {
			t.Errorf("Test %s: got: %t, expected %t", tc.name, got, tc.expect)
		}
	}
}
