package mp3

import (
	"github.com/bogem/id3v2"
)

type Mp3Tag struct {
	tag *id3v2.Tag
}

func NewMp3Tag(path string) (*Mp3Tag, error) {
	tag, err := id3v2.Open(path, id3v2.Options{Parse: true})
	if err != nil {
		return nil, err
	}

	mp3Tag := new(Mp3Tag)
	mp3Tag.tag = tag
	return mp3Tag, nil
}

func (m *Mp3Tag) SetCover(cover []byte, mime string) error {

	m.tag.AddAttachedPicture(id3v2.PictureFrame{
		Encoding:    id3v2.EncodingISO,
		MimeType:    mime,
		PictureType: id3v2.PTFrontCover,
		Description: "Front cover",
		Picture:     cover,
	})
	return nil
}

func (m *Mp3Tag) SetCoverUrl(coverUrl string) error {
	m.tag.AddAttachedPicture(id3v2.PictureFrame{
		Encoding:    id3v2.EncodingISO,
		MimeType:    "-->",
		PictureType: id3v2.PTFrontCover,
		Description: "Front cover",
		Picture:     []byte(coverUrl),
	})
	return nil
}

func (m *Mp3Tag) SetTitle(title string) error {

	if name := m.tag.Title(); name == "" {
		m.tag.SetTitle(title)
	}
	return nil
}

func (m *Mp3Tag) SetAlbum(album string) error {
	if name := m.tag.Album(); name == "" {
		m.tag.SetAlbum(album)
	}
	return nil
}

func (m *Mp3Tag) SetArtist(artists []string) error {
	if frames := m.tag.GetFrames(m.tag.CommonID("Artist")); len(frames) == 0 {
		for _, artist := range artists {
			m.tag.SetArtist(artist)
		}
	}
	return nil
}

func (m *Mp3Tag) SetComment(comment string) error {
	if frames := m.tag.GetFrames(m.tag.CommonID("Comments")); len(frames) == 0 {
		m.tag.AddCommentFrame(id3v2.CommentFrame{
			Encoding:    id3v2.EncodingISO,
			Language:    "XXX",
			Description: "",
			Text:        comment,
		})
	}
	return nil
}

func (m *Mp3Tag) Save() error {
	err := m.tag.Save()
	err = m.tag.Close()
	return err
}
