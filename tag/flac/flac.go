package flac

import (
	"github.com/go-flac/flacpicture"
	"github.com/go-flac/flacvorbis"
	"github.com/go-flac/go-flac"
)

type FlacTag struct {
	path string
	file *flac.File
	cmts *flacvorbis.MetaDataBlockVorbisComment
}

func NewFlacTag(path string) (*FlacTag, error) {
	// already read and closed
	f, err := flac.ParseFile(path)
	if err != nil {
		return nil, err
	}

	var cmtmeta *flac.MetaDataBlock
	for _, m := range f.Meta {
		if m.Type == flac.VorbisComment {
			cmtmeta = m
			break
		}
	}
	var cmts *flacvorbis.MetaDataBlockVorbisComment
	if cmtmeta != nil {
		cmts, err = flacvorbis.ParseFromMetaDataBlock(*cmtmeta)
		if err != nil {
			return nil, err
		}
	} else {
		cmts = flacvorbis.New()
	}

	tagger := new(FlacTag)
	tagger.file = f
	tagger.cmts = cmts
	tagger.path = path
	return tagger, nil
}

func (f *FlacTag) SetCover(buf []byte, mime string) error {
	picture, err := flacpicture.NewFromImageData(flacpicture.PictureTypeFrontCover, "Front cover", buf, mime)
	if err == nil {
		picturemeta := picture.Marshal()
		f.file.Meta = append(f.file.Meta, &picturemeta)
	}
	return err

}

func (f *FlacTag) SetCoverUrl(coverUrl string) error {
	picture := &flacpicture.MetadataBlockPicture{
		PictureType: flacpicture.PictureTypeFrontCover,
		MIME:        "-->",
		Description: "Front cover",
		ImageData:   []byte(coverUrl),
	}
	picturemeta := picture.Marshal()
	f.file.Meta = append(f.file.Meta, &picturemeta)
	return nil
}

func (f *FlacTag) SetTitle(title string) error {
	if titles, err := f.cmts.Get(flacvorbis.FIELD_TITLE); err != nil {
		return err
	} else if len(titles) == 0 {
		return f.cmts.Add(flacvorbis.FIELD_TITLE, title)
	}
	return nil
}

func (f *FlacTag) SetAlbum(album string) error {
	if albums, err := f.cmts.Get(flacvorbis.FIELD_ALBUM); err != nil {
		return err
	} else if len(albums) == 0 {
		return f.cmts.Add(flacvorbis.FIELD_ALBUM, album)
	}
	return nil
}

func (f *FlacTag) SetArtist(artists []string) error {
	if theArtists, err := f.cmts.Get(flacvorbis.FIELD_ARTIST); err != nil {
		return err
	} else if len(theArtists) == 0 {
		for _, artist := range artists {
			f.cmts.Add(flacvorbis.FIELD_ARTIST, artist)
		}
	}
	return nil
}

func (f *FlacTag) SetComment(string) error {
	// pass
	return nil
}

func (f *FlacTag) Save() error {
	res := f.cmts.Marshal()
	f.file.Meta = append(f.file.Meta, &res)
	return f.file.Save(f.path)
}
