package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/chai2010/webp"
	"golang.org/x/image/bmp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

func WebpEncoder(p1, p2 string, quality float32, Log bool, c chan int) (err error) {
	// if convert fails, return error; success nil
	var buf bytes.Buffer
	var img image.Image

	data, err := ioutil.ReadFile(p1)
	if err != nil {
		ChanErr(c)
		return
	}

	contentType := GetFileContentType(data[:512])
	if strings.Contains(contentType, "jpeg") {
		img, _ = jpeg.Decode(bytes.NewReader(data))
	} else if strings.Contains(contentType, "png") {
		img, _ = png.Decode(bytes.NewReader(data))
	} else if strings.Contains(contentType, "bmp") {
		img, _ = bmp.Decode(bytes.NewReader(data))
	} else if strings.Contains(contentType, "gif") {
		// TODO: need to support animated webp
		img, _ = gif.Decode(bytes.NewReader(data))
	}

	if img == nil {
		msg := "image file " + path.Base(p1) + " is corrupted or not supported"
		log.Println(msg)
		err = errors.New(msg)
		ChanErr(c)
		return
	}

	if err = webp.Encode(&buf, img, &webp.Options{Lossless: false, Quality: quality}); err != nil {
		log.Println(err)
		ChanErr(c)
		return
	}
	if err = ioutil.WriteFile(p2, buf.Bytes(), 0755); err != nil {
		log.Println(err)
		ChanErr(c)
		return
	}

	if Log {
		fmt.Printf("Save to %s ok\n", p2)
	}

	ChanErr(c)

	return nil
}