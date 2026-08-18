package messagebus_test

import (
	"strconv"
	"sync"
	"testing"

	"github.com/tvanriel/messagebus"
)

func BenchmarkPublish(b *testing.B) {
	bus := messagebus.New()

	var wg sync.WaitGroup

	wg.Add(b.N)

	_ = bus.Subscribe("topic", func() {
		wg.Done()
	})

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bus.Publish("topic")
		}
	})

	wg.Wait()
}

func BenchmarkSubscribe(b *testing.B) {
	bus := messagebus.New()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = bus.Subscribe("topic", func() {})
		}
	})
}

func benchmark(b *testing.B, subscribersCount, topicsCount int) {
	b.Helper()

	bus := messagebus.New()

	var wg sync.WaitGroup

	wg.Add(b.N * subscribersCount)

	for i := range topicsCount {
		for range subscribersCount {
			_ = bus.Subscribe(messagebus.Topic(strconv.Itoa(i)), func() {
				wg.Done()
			})
		}
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bus.Publish(messagebus.NewRouteKeyFromString(strconv.Itoa(topicsCount - 1)))
		}
	})

	wg.Wait()
}

func Benchmark1Subscriber1Topic(b *testing.B) {
	benchmark(b, 1, 1)
}

func Benchmark1Subscriber100Topics(b *testing.B) {
	benchmark(b, 1, 100)
}

func Benchmark100Subscribers1Topic(b *testing.B) {
	benchmark(b, 100, 1)
}

func Benchmark100Subscribers100Topics(b *testing.B) {
	benchmark(b, 100, 100)
}
