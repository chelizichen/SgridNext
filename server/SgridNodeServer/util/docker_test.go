package util

import (
	"testing"
)

func TestDockerGetAlive(t *testing.T) {
	alive, err := DockerGetAlive("sgridnode")
	t.Log("alive", alive)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(alive)
}
