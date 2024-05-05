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
