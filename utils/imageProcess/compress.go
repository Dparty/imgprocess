package imageProcess

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/nfnt/resize"
)

func CompressImg(img image.Image, im image.Config, scale float64, contentType string) []byte {
	fmt.Printf("origin size %d %d\n\n", im.Height, im.Width)
	h := float64(im.Height) * scale
	w := float64(im.Width) * scale
	fmt.Printf("reize %f %f\n\n", h, w)

	m := resize.Resize(0, uint(h), img, resize.Lanczos3)
	buf := new(bytes.Buffer)

	switch {
	case contentType == "image/jpeg":
		fmt.Println("jpeg")
		jpeg.Encode(buf, m, nil)
	case contentType == "image/png":
		fmt.Println("png")
		png.Encode(buf, m)
	case contentType == "image/jpg":
		fmt.Println("jpg")
		jpeg.Encode(buf, m, nil)
	default: //default:當前面條件都沒有滿足時將會執行此處內包含的方法
		jpeg.Encode(buf, m, nil)
	}

	send_s3 := buf.Bytes()
	return send_s3
}
