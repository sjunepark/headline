package scraper

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func setScreenSize(rodPage *rod.Page, windowFullscreen bool) error {
	var err error

	if windowFullscreen {
		err = setScreenSizeFull(rodPage)
		if err != nil {
			return err
		}
		return nil
	}

	err = setScreenSizeDefault(rodPage)
	if err != nil {
		return err
	}
	return nil
}

func setScreenSizeFull(rodPage *rod.Page) error {
	return rodPage.SetWindow(&proto.BrowserBounds{
		WindowState: proto.BrowserWindowStateFullscreen,
	})
}

func setScreenSizeDefault(rodPage *rod.Page) error {
	return rodPage.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width:             1920,
		Height:            1080,
		DeviceScaleFactor: 1,
		Mobile:            false,
	})
}
