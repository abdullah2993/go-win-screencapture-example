package main

import (
	"encoding/binary"
	"os"
	"unsafe"

	"github.com/gonutz/w32"
)

//Capture the screen using GDI code taken from https://github.com/gonutz/cap_screen
func captureWithGDI(path string) {
	hdcWindow := w32.GetDC(0)
	hdcMemDC := w32.CreateCompatibleDC(hdcWindow)
	if hdcMemDC == 0 {
		panic("CreateCompatibleDC has failed")
	}

	height := w32.GetSystemMetrics(w32.SM_CYVIRTUALSCREEN)
	width := w32.GetSystemMetrics(w32.SM_CXVIRTUALSCREEN)

	hbmScreen := w32.CreateCompatibleBitmap(
		hdcWindow,
		width,
		height,
	)
	if hbmScreen == 0 {
		panic("CreateCompatibleBitmap Failed")
	}

	w32.SelectObject(hdcMemDC, w32.HGDIOBJ(hbmScreen))
	if !w32.BitBlt(hdcMemDC,
		0, 0,
		width, height,
		hdcWindow,
		0, 0,
		w32.SRCCOPY) {
		panic("BitBlt has failed")
	}

	var bmpScreen w32.BITMAP
	w32.GetObject(
		w32.HGDIOBJ(hbmScreen),
		uintptr(binary.Size(bmpScreen)),
		unsafe.Pointer(&bmpScreen),
	)

	var bmfHeader w32.BITMAPFILEHEADER
	var bi w32.BITMAPINFOHEADER
	bi.BiSize = uint32(binary.Size(bi))
	bi.BiWidth = bmpScreen.BmWidth
	bi.BiHeight = bmpScreen.BmHeight
	bi.BiPlanes = 1
	bi.BiBitCount = 32
	bi.BiCompression = w32.BI_RGB
	bi.BiSizeImage = 0
	bi.BiXPelsPerMeter = 0
	bi.BiYPelsPerMeter = 0
	bi.BiClrUsed = 0
	bi.BiClrImportant = 0

	dwBmpSize := uint32(((bmpScreen.BmWidth*int32(bi.BiBitCount) + 31) / 32) * 4 * bmpScreen.BmHeight)
	colorData := make([]byte, dwBmpSize)
	w32.GetDIBits(
		hdcWindow,
		hbmScreen,
		0,
		uint(bmpScreen.BmHeight),
		unsafe.Pointer(&colorData[0]),
		(*w32.BITMAPINFO)(unsafe.Pointer(&bi)),
		w32.DIB_RGB_COLORS,
	)

	bmfHeader.BfOffBits = uint32(binary.Size(bmfHeader) + binary.Size(bi))
	bmfHeader.BfSize = dwBmpSize + uint32(binary.Size(bmfHeader)+binary.Size(bi))
	bmfHeader.BfType = 0x4D42

	f, err := os.Create(path)
	check(err)
	defer f.Close()

	check(binary.Write(f, binary.LittleEndian, bmfHeader))
	check(binary.Write(f, binary.LittleEndian, bi))
	_, err = f.Write(colorData)
	check(err)

	w32.DeleteObject(w32.HGDIOBJ(hbmScreen))
	w32.DeleteObject(w32.HGDIOBJ(hdcMemDC))
	w32.ReleaseDC(0, hdcWindow)
}
