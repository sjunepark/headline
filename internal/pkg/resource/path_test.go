package resource

import (
	"github.com/sejunpark/headline/internal/pkg/constant"
	"github.com/stretchr/testify/assert"
	"path"
	"runtime"
	"strings"
	"testing"
)

func TestChromiumPath(t *testing.T) {
	chromiumPath, err := ChromiumPath()
	assert.NoError(t, err)

	suffix := path.Join(constant.APP_NAME, "chromium")
	if runtime.GOOS == "darwin" {
		suffix = path.Join("Application Support", suffix)
	}

	assert.True(t, strings.HasSuffix(chromiumPath, suffix))
}
