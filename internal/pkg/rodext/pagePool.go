package rodext

import (
	"context"
	"log/slog"
)

type pagePool struct {
	ctx  context.Context
	pool chan *Page
}

// newPagePool initializes a new pagePool.
// The pool is filled with nil values.
func newPagePool(ctx context.Context, size int) *pagePool {
	p := &pagePool{
		ctx:  ctx,
		pool: make(chan *Page, size),
	}

	// You have to fill the pagePool with nil values.
	// This is to avoid the need to handle complexities when the pool is not full,
	// such as when handling waiting channel receivers and channel closing.
	// Also, if you don't fill it with nil values first, it's hard to properly clean them up.
	// There are lots of complexities to handle when a channel is not full.
	for i := 0; i < size; i++ {
		p.pool <- nil
	}
	return p
}

// cleanup closes the pool and cleans up all the Pages in the pool.
// It should be called as a result of context cancellation,
// since it doesn't call cancel to broadcast the cancellation signal to other functions or goroutines.
func (pp *pagePool) cleanup() {
	close(pp.pool)
	for page := range pp.pool {
		if page == nil {
			continue
		}
		page.cleanup()
	}
}

// Get returns a new Page from pool.
// If the pool returns nil, a new Page is created using the newPage function provided.
// Always make sure to call putPage after using the Page.
//
// The returned Page is thread-safe.
func (pp *pagePool) Get(newPage func() (*Page, func(), error)) (p *Page, putPage func(), err error) {
	if pp.ctx.Err() != nil {
		return nil, nil, pp.ctx.Err()
	}

	p = <-pp.pool

	if p == nil {
		// Don't need to handle cleanup. The page will be returned to the pool after use.
		p, _, err = newPage()
		if err != nil {
			return nil, nil, err
		}
		slog.Debug("Page: no Page in pool, created new Page", "address", p)
		return p, pp.putPageFactory(p), nil
	}

	slog.Debug("Page: got Page from pool", "address", p)
	return p, pp.putPageFactory(p), nil
}

// putPageFactory returns a function(a closure) that puts the Page back to pages.
// The user doesn't have to remember the reference to the Page, as it is already captured in the closure.
//
// putPage gracefully handles context cancellation.
func (pp *pagePool) putPageFactory(p *Page) func() {
	return func() {
		if pp.ctx.Err() != nil {
			p.cleanup()
			slog.Debug("putPage: context cancelled, cleaning up Page", "address", p)
			return
		}
		pp.pool <- p
		slog.Debug("putPage: putting back Page", "address", p)
	}
}

func (pp *pagePool) len() int {
	return len(pp.pool)
}
