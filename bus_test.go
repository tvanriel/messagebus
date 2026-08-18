package messagebus_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/tvanriel/messagebus"
)

func TestNew(t *testing.T) {
	t.Parallel()

	bus := messagebus.New()

	if bus == nil {
		t.Fail()
	}
}

func TestSubscribe(t *testing.T) {
	t.Parallel()

	bus := messagebus.New()

	if bus.Subscribe("test", func() {}) != nil {
		t.Fail()
	}

	if bus.Subscribe("test", 2) == nil {
		t.Fail()
	}
}

func TestUnsubscribe(t *testing.T) {
	t.Parallel()

	bus := messagebus.New()

	handler := func() {}

	err := bus.Subscribe("test", handler)
	if err != nil {
		t.Fatal(err)
	}

	err = bus.Unsubscribe("test", handler)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestClose(t *testing.T) {
	t.Parallel()

	bus := messagebus.New()

	handler := func() {}

	err := bus.Subscribe("test", handler)
	if err != nil {
		t.Fatal(err)
	}

	if bus.NumHandlers() == 0 {
		fmt.Println("Did not subscribed handler to topic")
		t.Fail()
	}

	bus.Close("test")

	if bus.NumHandlers() != 0 {
		fmt.Println("Did not unsubscribed handlers from topic")
		t.Fail()
	}
}

func TestPublish(t *testing.T) {
	t.Parallel()

	bus := messagebus.New()

	var wg sync.WaitGroup

	wg.Add(2)

	first := false
	second := false

	err := bus.Subscribe("topic.*", func(v bool) {
		defer wg.Done()

		first = v
	})
	if err != nil {
		t.Fatal(err)
	}

	err = bus.Subscribe("topic.*", func(v bool) {
		defer wg.Done()

		second = v
	})
	if err != nil {
		t.Fatal(err)
	}

	bus.Publish("topic.test", true)

	wg.Wait()

	if first == false || second == false {
		t.Fail()
	}
}

func TestHandleError(t *testing.T) {
	t.Parallel()

	bus := messagebus.New()

	err := bus.Subscribe("topic", func(out chan<- error) {
		out <- errors.New("I do throw error") //nolint:err113 // Used for test, not for comparison.
	})
	if err != nil {
		t.Fatal(err)
	}

	out := make(chan error)
	defer close(out)

	bus.Publish("topic", out)

	if <-out == nil {
		t.Fail()
	}
}
