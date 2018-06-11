package main

func main() {
	captureWithGDI("gdi.bmp")
	captureWithDirectX("directx.bmp")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
