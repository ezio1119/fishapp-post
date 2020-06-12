package interactor

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/disintegration/imaging"
	"github.com/ezio1119/fishapp-post/conf"
)

type resizeChan struct {
	io   io.Reader
	name string
}

// 縦,横を変えずにLanczosで、png, jpeg, gifを圧縮する
func resizeImage(r io.Reader, ch chan resizeChan, imgName string) error {

	img, t, err := image.Decode(r)
	if err != nil {
		return err
	}

	nrgba := imaging.Fit(img, conf.C.Sv.ImageWidth, conf.C.Sv.ImageHeight, imaging.Lanczos)
	buf := &bytes.Buffer{}

	switch t {
	case "jpeg":
		if err := jpeg.Encode(buf, nrgba, &jpeg.Options{Quality: jpeg.DefaultQuality}); err != nil {
			return err
		}
		imgName = imgName + ".jpg"
	case "png":
		if err := png.Encode(buf, nrgba); err != nil {
			return err
		}
		imgName = imgName + ".png"
	case "gif":
		if err := gif.Encode(buf, nrgba, nil); err != nil {
			return err
		}
		imgName = imgName + ".gif"
	}

	ch <- resizeChan{io: buf, name: imgName}

	return nil
}
