package resource

import (
	"github.com/adrg/xdg"
	"github.com/sejunpark/headline/internal/pkg/constant"
	"path"
)

func ChromiumPath() (string, error) {
	data, err := xdg.DataFile(constant.APP_NAME)
	if err != nil {
		return "", err
	}
	return path.Join(data, "chromium"), nil
}
