package handler

import "testing"

func TestRandom(t *testing.T) {

	vv := random(16)

	t.Logf("info for random %s", vv)

	vv = random(128)

	t.Logf("info for random %s", vv)

	vv = random(10)

	t.Logf("info for random %s", vv)

}
