package ncm

import (
	"log"
	"testing"
)

func TestGetData(t *testing.T) {
	nf, err := NewNcmFile("../Perfect.ncm")
	if err != nil {
		log.Printf("Open Ncm File failed: %v", err)
		return
	}

	err = nf.Validate()
	if err != nil {
		log.Printf("Ncm magic header check failed: %v", err)
	}

	err = nf.GetKey()
	if err != nil {
		log.Printf("Get Key Failed: %v", err)
	}

	log.Printf("Ncm File's key: %v", nf.Key.Detail)

	err = nf.GetMeta()
	if err != nil {
		log.Printf("Get Meta Failed: %v", err)
	}

	err = nf.GetCover()
	if err != nil {
		log.Printf("Get Cover Failed: %v", err)
	}

	err = nf.GetMusicData()
	if err != nil {
		log.Printf("Get Music Data Failed: %v", err)
	}
}
