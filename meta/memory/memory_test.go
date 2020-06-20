package memory_test

import (
	"testing"

	memory "github.com/antonio-alexander/go-bludgeon/meta/memory"
	tests "github.com/antonio-alexander/go-bludgeon/meta/tests"
)

func TestIntTimerReadWrite(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	json := memory.New()
	tests.TestIntTimerReadWrite(t, json)
}

func TestIntDelete(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	json := memory.New()
	tests.TestIntDelete(t, json)
}

func TestIntSliceReadWrite(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	json := memory.New()
	tests.TestIntSliceReadWrite(t, json)
}

func TestIntSliceDelete(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	json := memory.New()
	tests.TestIntSliceDelete(t, json)
}

func TestIntSliceTimer(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	json := memory.New()
	tests.TestIntSliceTimer(t, json)
}

func TestIntTimerActiveSlice(t *testing.T) {
	//Test:
	//Notes:
	//Verification:

	json := memory.New()
	tests.TestIntTimerActiveSlice(t, json)
}

//TODO: write test for deleting a timer
//TODO: write test for calculating elapsed time on
// an active time slice
