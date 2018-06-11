package main

import (
	"github.com/abdullah2993/go-win-screencapture-example/d3dx"
	"github.com/gonutz/d3d9"
)

//Captures the screen using DirectX9 also relies on d3dx_43.dll
func captureWithDirectX(path string) {
	dx, err := d3d9.Create(d3d9.SDK_VERSION)
	check(err)

	mode, err := dx.GetAdapterDisplayMode(d3d9.ADAPTER_DEFAULT)
	check(err)

	dev, _, err := dx.CreateDevice(d3d9.ADAPTER_DEFAULT,
		d3d9.DEVTYPE_HAL,
		d3d9.HWND(0),
		d3d9.CREATE_SOFTWARE_VERTEXPROCESSING,
		d3d9.PRESENT_PARAMETERS{
			Windowed:         1,
			HDeviceWindow:    d3d9.HWND(0),
			SwapEffect:       d3d9.SWAPEFFECT_DISCARD,
			BackBufferFormat: d3d9.FMT_A8R8G8B8,
			BackBufferWidth:  mode.Width,
			BackBufferHeight: mode.Height,
			BackBufferCount:  1,
		},
	)
	check(err)

	s, err := dev.CreateOffscreenPlainSurface(uint(mode.Width), uint(mode.Height), d3d9.FMT_A8R8G8B8, d3d9.POOL_SYSTEMMEM, 0)
	check(err)

	err = dev.GetFrontBufferData(0, s)
	check(err)

	err = d3dx.SaveSurfaceToFile(path, d3dx.IFFBMP, s, nil, nil)
	check(err)
}
