package builder

import (
	"github.com/stretchr/testify/suite"
)

type BaseScraperBuilderSuite struct {
	suite.Suite
	Builder Builder
	cleanup func()
}

func (ts *BaseScraperBuilderSuite) SetupScraperBuilderSuite(builder Builder, cleanup func()) {
	ts.Builder = builder
	ts.cleanup = cleanup
}

func (ts *BaseScraperBuilderSuite) SetupTearDownSuite() {
	ts.cleanup()
}
