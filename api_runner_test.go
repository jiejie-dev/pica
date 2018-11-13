package pica

import "testing"

func TestApiRunner_Run(t *testing.T) {
	runner := NewAPIRunnerFromFile("sample/pica.fun", nil, 0)
	runner.Run()
}
