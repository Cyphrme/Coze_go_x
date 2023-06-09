package mapslice

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/cyphrme/coze"
)

var GoldenMapSLice = MapSlice{
	MapItem{Key: "abc", Value: 123},
	MapItem{Key: "def", Value: "456"},
}

func ExampleMapSlice_marshal() {
	b, err := coze.Marshal(GoldenMapSLice)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", b)
	// Output:
	// {"abc":123,"def":"456"}
}

func ExampleMapSlice_unmarshal() {
	ms := MapSlice{}
	err := json.Unmarshal([]byte(`{"abc":123,"def":"456"}`), &ms)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", ms)

	// Output:
	// [{abc 123} {def 456}]
}

func ExampleMapSlice_Keys() {
	fmt.Println(GoldenMapSLice.Keys())

	// Output:
	// [abc def]
}

func ExampleMapSlice_Values() {
	fmt.Println(GoldenMapSLice.Values())

	// Output:
	// [123 456]
}

// Demonstrates that marshaling a channel is invalid.
func TestMapSlice_Marshal_chan(t *testing.T) {
	ms := MapSlice{
		MapItem{Key: "abc", Value: make(chan int)},
	}

	e := "json: error calling MarshalJSON for type coze.MapSlice: json: unsupported type: chan int"
	_, err := json.Marshal(ms)
	if err != nil && e != err.Error() {
		t.Fatalf("expected: %s\ngot: %v", e, err)
	}
}

// See
//
//	https://go.dev/ref/mem
//	https://go.dev/doc/articles/race_detector
//
// Use:
//
//	go test -race -run TestMapSlice_race
func TestMapSlice_race(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var item MapItem
			json.Unmarshal([]byte(`{}`), &item)
		}()
	}
	wg.Wait()
}
