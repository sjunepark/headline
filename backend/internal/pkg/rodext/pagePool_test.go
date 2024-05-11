package rodext

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PagePoolSuite struct {
	suite.Suite
	ctx     context.Context
	pool    *pagePool
	cleanup func()
}

func (ts *PagePoolSuite) SetupTest() {
	ts.ctx = context.Background()
	ts.pool, ts.cleanup = newPagePool(ts.ctx, 16)
}

func (ts *PagePoolSuite) TearDownTest() {
	ts.cleanup()
}

func TestPagePoolSuite(t *testing.T) {
	suite.Run(t, new(PagePoolSuite))
}

func (ts *PagePoolSuite) TestPagePool_newPagePool() {
	ts.Run("All pages should be nil after initialization", func() {
		for i := 0; i < 16; i++ {
			page := <-ts.pool.pool
			ts.Nil(page)
		}
	})
}

type PagePoolCleanupSuite struct {
	suite.Suite
	ctx     context.Context
	pool    *pagePool
	cleanup func()
}

func (ts *PagePoolCleanupSuite) SetupTest() {
	ts.ctx = context.Background()
	ts.pool, ts.cleanup = newPagePool(ts.ctx, 16)
}

func TestPagePoolCleanupSuite(t *testing.T) {
	suite.Run(t, new(PagePoolCleanupSuite))
}

func (ts *PagePoolCleanupSuite) TestPagePool_cleanup() {
	ts.Run("All pages should be cleaned up after cleanup", func() {
		ts.cleanup()

		noPage, open := <-ts.pool.pool
		ts.Falsef(open, "pool should be closed")
		ts.Nilf(noPage, "pool should be empty but received %v", noPage)
	})
}
