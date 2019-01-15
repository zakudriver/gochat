package gochat

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"unsafe"
)

func errHandler(s string, e error) error {
	if e == nil {
		return fmt.Errorf(" ## %s", s)
	} else {
		return fmt.Errorf(" ## %s => error: %s", s, e)
	}
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func logInfo(s interface{}) {
	log.Println("** SUCCESS **  ", s)
}

func logErr(s string) {
	log.Println("** FAIL **  ", s)
	os.Exit(1)
}

func qrcodeHandler(b []byte) error {
	img, err := jpeg.Decode(bytes.NewReader([]byte(b)))
	if err != nil {
		return err
	}

	f, err := os.Create("src/qrcode.png")
	defer f.Close()
	if err != nil {
		return err
	}
	png.Encode(f, img)
	return nil
}

func qrcodeHttp(b []byte) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(b)
	})
	http.ListenAndServe(":8004", nil)
}

func jsonFileOut(b []byte) {
	f, err := os.Create("src/file.json")
	defer f.Close()
	if err != nil {
		logErr(err.Error())
	}

	if _, err := f.Write(b); err != nil {
		logErr(err.Error())
	}
}
