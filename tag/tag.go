package tag

import (
	"errors"
	"log"
	"strings"

	"github.com/closetool/NCMConverter/converter"
	"github.com/closetool/NCMConverter/tag/flac"
	"github.com/closetool/NCMConverter/tag/mp3"
)

const (
	mimeJPEG = "image/jpeg"
	mimePNG  = "image/png"
)

var (
	ErrFormat = errors.New("Mp3,FLAC support only!")
)

type Tagger interface {
	SetCover(cover []byte, mime string) error
	SetCoverUrl(coverUrl string) error
	SetTitle(string) error
	SetAlbum(string) error
	SetArtist([]string) error
	SetComment(string) error
	Save() error
}

func TagAudioFileFromMeta(tag Tagger, imgData []byte, meta *converter.Meta) error {
	if (imgData == nil || len(imgData) == 0) && meta.Album.CoverUrl != "" {
		if coverData, err := fetchUrl(meta.Album.CoverUrl); err != nil {
			log.Println(err)
		} else {
			imgData = coverData
		}
	}

	if imgData != nil || len(imgData) == 0 {
		picMIME := mimeJPEG
		if isPNGHeader(imgData) {
			picMIME = mimePNG
		}
		tag.SetCover(imgData, picMIME)
	} else if meta.Album.CoverUrl != "" {
		tag.SetCoverUrl(meta.Album.CoverUrl)
	}

	if meta.Name == "" {
		tag.SetTitle(meta.Name)
	}

	if meta.Album.Name == "" {
		tag.SetAlbum(meta.Name)
	}

	artists := make([]string, 0)
	for _, artist := range meta.Artists {
		artists = append(artists, artist.Name)
	}
	if len(artists) > 0 {
		tag.SetArtist(artists)
	}

	if meta.Comment != "" {
		tag.SetComment(meta.Comment)
	}

	return tag.Save()
}

func NewTagger(path, format string) (tag Tagger, err error) {
	switch strings.ToLower(format) {
	case "mp3":
		tag, err = mp3.NewMp3Tag(path)
	case "flac":
		tag, err = flac.NewFlacTag(path)
	default:
		err = ErrFormat
	}
	return
}
