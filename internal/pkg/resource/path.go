package resource

import (
	"github.com/adrg/xdg"
	"github.com/cockroachdb/errors"
	"github.com/sejunpark/headline/internal/pkg/constant"
	"path"
)

func ChromiumPath() (string, error) {
	data, err := xdg.DataFile(constant.APP_NAME)
	if err != nil {
		return "", errors.Wrapf(err, "xdg.DataFile(%s) failed", constant.APP_NAME)
	}
	return path.Join(data, "chromium"), nil
}
