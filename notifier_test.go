package notifier

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func Test_Notify(t *testing.T) {
	n := NewNotifier()
	done := make(chan struct{})
	notified := false
	go func() {
		<-n.Wait()
		notified = true
		done <- struct{}{}
	}()
	time.Sleep(100 * time.Millisecond)
	n.Notify()
	<-done
	if !notified {
		t.Fatal("did not notify")
	}
}

func Test_Broadcast(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	g, _ := errgroup.WithContext(ctx)

	n := NewNotifier()
	notified := atomic.Int32{}
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			<-n.Wait()
			notified.Add(1)
			return nil
		})
	}
	time.Sleep(100 * time.Millisecond)
	n.Notify()
	if err := g.Wait(); err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}

	got := notified.Load()
	if got != 10 {
		t.Fatalf("Got %d notifications, want 10", got)
	}
}

func Test_Reuse(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	g, _ := errgroup.WithContext(ctx)

	n := NewNotifier()
	notified := atomic.Int32{}
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			<-n.Wait()
			<-n.Wait()
			notified.Add(1)
			return nil
		})
	}
	time.Sleep(100 * time.Millisecond)
	n.Notify()
	time.Sleep(100 * time.Millisecond)
	n.Notify()
	if err := g.Wait(); err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}

	got := notified.Load()
	if got != 10 {
		t.Fatalf("Got %d notifications, want 10", got)
	}
}
