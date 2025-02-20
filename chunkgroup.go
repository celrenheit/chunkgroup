package chunkgroup

import (
	"context"
	"slices"

	"golang.org/x/sync/errgroup"
)

// Group is a struct that holds a slice of items and a function to execute on them.
// The function is executed when the slice reaches its capacity.
// The function is executed concurrently with a limit of concurrency.
type Group[T any] struct {
	ctx   context.Context
	items []T
	fn    ExecFunc[T]
	eg    *errgroup.Group
}

// ExecFunc is a function that takes a slice of items and returns an error.
type ExecFunc[T any] func(context.Context, []T) error

// New creates a new Group with the given size and concurrency.
func New[T any](size, concurrency int, fn ExecFunc[T]) *Group[T] {
	cg, _ := WithContext(context.Background(), size, concurrency, fn)
	return cg
}

// WithContext creates a new Group with the given context, size, and concurrency.
// The context is used to cancel the execution of the function.
func WithContext[T any](ctx context.Context, size, concurrency int, fn ExecFunc[T]) (*Group[T], context.Context) {
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(concurrency)

	return &Group[T]{
		ctx:   ctx,
		items: make([]T, 0, size),
		fn:    fn,
		eg:    eg,
	}, ctx
}

// Add adds an item to the Group.
func (g *Group[T]) Add(item T) {
	g.items = append(g.items, item)

	if len(g.items) == cap(g.items) {
		tmp := slices.Clone(g.items)
		g.eg.Go(func() error {
			return g.fn(g.ctx, tmp)
		})

		g.items = g.items[:0]
	}
}

// Flush executes the function on the remaining items.
// It MUST be called after all items have been added.
func (g *Group[T]) Flush() error {
	if len(g.items) > 0 {
		tmp := slices.Clone(g.items)
		g.eg.Go(func() error { return g.fn(g.ctx, tmp) })

		g.items = g.items[:0]
	}

	return g.eg.Wait()
}
