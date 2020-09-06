package converter

import (
	"log"
	"os"
	"testing"

	"github.com/closetool/NCMConverter/ncm"
)

var nf *ncm.NcmFile
var cv *Converter

func init() {
	nf, _ = ncm.NewNcmFile("../Perfect.ncm")
	_ = nf.Parse()
	cv = NewConverter(nf)
}

func TestParseKey(t *testing.T) {
	except := []byte{110, 101, 116, 101, 97, 115, 101, 99, 108, 111, 117, 100, 109, 117, 115, 105, 99, 50, 51, 49, 56, 50,
		51, 51, 52, 56, 50, 48, 48, 52, 57, 69, 55, 102, 84, 52, 57, 120, 55, 100, 111, 102, 57, 79, 75, 67, 103,
		103, 57, 99, 100, 118, 104, 69, 117, 101, 122, 121, 51, 105, 90, 67, 76, 49, 110, 70, 118, 66, 70,
		100, 49, 84, 52, 117, 83, 107, 116, 65, 74, 75, 109, 119, 90, 88, 115, 105, 106, 80, 98, 105, 106, 108,
		105, 105, 111, 110, 86, 85, 88, 88, 103, 57, 112, 108, 84, 98, 88, 69, 99, 108, 65, 69, 57, 76, 98}
	except = except[17:]
	err := cv.HandleKey()
	if err != nil {
		log.Printf("an error occured: %v", err)
		return
	}

	for i, b := range cv.KeyData {
		if except[i] != b {
			panic("result does not match excepted!")
		}
	}
	log.Println(string(cv.KeyData))
}

func TestParseMeta(t *testing.T) {
	err := cv.HandleMeta()
	if err != nil {
		log.Printf("an error occured: %v", err)
	}
	//	metaJSON := `{
	//	"musicId": 35539156,
	//	"musicName": "Perfect",
	//	"artist": [
	//		["One Direction", 98351]
	//	],
	//	"bitrate": 320000,
	//	"mp3DocId": "b902149837040f19c58fb57afc52fcba",
	//	"duration": 230333,
	//	"mvId": 493092,
	//	"alias": [],
	//	"transNames": [],
	//	"format": "mp3"
	//}`
	//	albumJSON := `{
	//	"albumId": 3317771,
	//	"album": "Perfect",
	//	"albumPic": "https://p3.music.126.net/jEhH4HtXX4Bkg4NGPURFaw==/3367804117152954.jpg",
	//}`
	log.Printf("Handled meta=%s", cv.MetaData)
	log.Printf("Handled MusicInfo=%s", cv.MetaData.Album)
}

func TestParseMusic(t *testing.T) {
	err := cv.HandleMusic()
	if err != nil {
		log.Printf("An error occured: %v", err)
		return
	}
	fd, err := os.OpenFile("../output.mp3", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return
	}
	fd.Write(cv.MusicData)
}
