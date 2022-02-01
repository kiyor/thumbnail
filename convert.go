package thumbnail

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	_ "image/jpeg"

	"github.com/disintegration/imaging"
)

// Convert ...
func (a *Agent) Convert(dir, jpg string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	fs, err := d.Readdir(0)
	if err != nil {
		return err
	}
	var files []string
	for _, v := range fs {
		if strings.HasSuffix(v.Name(), ".jpg") {
			files = append(files, filepath.Join(dir, v.Name()))
		}
	}
	var thumbnails []image.Image
	var width, height int
	for _, file := range files {
		img, err := imaging.Open(file)
		if err != nil {
			return err
		}
		if width == 0 {
			f, _ := os.Open(file)
			i, _, _ := image.DecodeConfig(f)
			width = i.Width
			height = i.Height
		}
		thumb := imaging.Clone(img)
		thumbnails = append(thumbnails, thumb)
	}
	var dst image.Image

	if len(thumbnails) >= 60 {
		dst = imagePaste3X20(width, height, thumbnails)
	} else {
		dst = imagePasteVertical(width, height, thumbnails)
	}

	err = imaging.Save(dst, jpg)
	if err != nil {
		return err
	}
	if a.userEnabled {
		err = os.Chown(jpg, a.user[0], a.user[1])
		if err != nil {
			return err
		}
	}
	return nil
}

func imagePasteVertical(width, height int, imgs []image.Image) image.Image {
	dst := imaging.New(width, height*len(imgs), color.NRGBA{0, 0, 0, 0})
	for i, img := range imgs {
		dst = imaging.Paste(dst, img, image.Pt(0, i*height))
	}
	return dst
}

func imagePaste3X20(width, height int, imgs []image.Image) image.Image {
	dst := imaging.New(3*width, 20*height, color.NRGBA{0, 0, 0, 0})
	id := 0
	for col := 0; col < 20; col++ {
		for row := 0; row < 3; row++ {
			dst = imaging.Paste(dst, imgs[id], image.Pt(row*width, col*height))
			id++
		}
	}
	return dst
}
