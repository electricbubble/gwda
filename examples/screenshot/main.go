package main

import (
	"github.com/electricbubble/gwda"
	"image"
	"log"
)

func main() {
	driver, err := gwda.NewUSBDriver(nil)
	if err != nil {
		log.Fatalln(err)
	}

	screenshot, err := driver.Screenshot()
	if err != nil {
		log.Fatal(err)
	}

	img, format, err := image.Decode(screenshot)
	if err != nil {
		log.Fatal(err)
	}
	_, _ = img, format
	// userHomeDir, _ := os.UserHomeDir()
	// file, err := os.Create(userHomeDir + "/Desktop/s1." + format)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer func() { _ = file.Close() }()
	// switch format {
	// case "png":
	// 	err = png.Encode(file, img)
	// case "jpeg":
	// 	err = jpeg.Encode(file, img, nil)
	// }
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(file.Name())
}
