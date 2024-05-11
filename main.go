package main

import (
	"context"
	"embed"
	"github.com/sejunpark/headline/backend/app"
	"github.com/sejunpark/headline/backend/constant"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/build
var assets embed.FS

func main() {
	// Create an instance of the app structure
	scrape := app.Scrape()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "headline",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			scrape.Start(ctx)
		},
		Bind: []interface{}{
			scrape,
		},
		EnumBind: []interface{}{
			constant.Sources,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
