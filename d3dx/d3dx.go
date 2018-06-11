package d3dx

import (
	"errors"
	"syscall"
	"unsafe"

	"github.com/gonutz/d3d9"
)

// ImageFileFormat Describes the supported image file formats
type ImageFileFormat int

const (
	// IFFBMP Windows bitmap (BMP) file format.
	IFFBMP ImageFileFormat = iota
	// IFFJPG Joint Photographics Experts Group (JPEG) compressed file format.
	IFFJPG
	// IFFTGA Truevision (Targa, or TGA) image file format.
	IFFTGA
	// IFFPNG Portable Network Graphics (PNG) file format.
	IFFPNG
	// IFFDDS DirectDraw surface (DDS) file format.
	IFFDDS
	// IFFPPM Portable pixmap (PPM) file format.
	IFFPPM
	// IFFDIB Windows device-independent bitmap (DIB) file format.
	IFFDIB
	// IFFHDR High dynamic range (HDR) file format.
	IFFHDR
	// IFFPFM Portable float map file format.
	IFFPFM
)

var (
	dll = syscall.NewLazyDLL("d3dx9_43.dll")
)

// SaveSurfaceToFile Saves a surface to a file.
// https://msdn.microsoft.com/en-us/library/windows/desktop/bb205431(v=vs.85).aspx
func SaveSurfaceToFile(destFile string,
	destFormat ImageFileFormat,
	srcSurface *d3d9.Surface,
	srcPalette *d3d9.PALETTEENTRY,
	srcRect *d3d9.RECT) error {

	lp := dll.NewProc("D3DXSaveSurfaceToFileW")
	ret, _, _ := lp.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(destFile))),
		uintptr(destFormat),
		uintptr(unsafe.Pointer(srcSurface)),
		uintptr(unsafe.Pointer(srcPalette)),
		uintptr(unsafe.Pointer(srcRect)))
	if int(ret) == d3d9.ERR_INVALIDCALL {
		return errors.New("d3dx: invalid call")
	}
	return nil
}
