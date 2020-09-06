package converter

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"math"
	"sync"

	"github.com/closetool/NCMconverter/ncm"
)

var (
	aesCoreKey   = []byte{0x68, 0x7A, 0x48, 0x52, 0x41, 0x6D, 0x73, 0x6F, 0x35, 0x6B, 0x49, 0x6E, 0x62, 0x61, 0x78, 0x57}
	aesModifyKey = []byte{0x23, 0x31, 0x34, 0x6C, 0x6A, 0x6B, 0x5F, 0x21, 0x5C, 0x5D, 0x26, 0x30, 0x55, 0x3C, 0x27, 0x28}
)

type Meta struct {
	Id       float64  `json:"musicId"`
	Name     string   `json:"musicName"`
	Album    *Album   `json:"-"`
	Artists  []Artist `json:"artist"`
	BitRate  float64  `json:"bitrate"`
	Duration float64  `json:"duration"`
	Format   string   `json:"format"`
	Comment  string   `json:"-"`
}

func (m *Meta) String() string {
	res, _ := json.Marshal(m)
	return string(res)
}

type Album struct {
	Id       float64 `json:"albumId"`
	Name     string  `json:"album"`
	CoverUrl string  `json:"albumPic"`
}

func (a *Album) String() string {
	res, _ := json.Marshal(a)
	return string(res)
}

type Artist struct {
	Name string
	Id   float64
}

func (a *Artist) UnmarshalJSON(data []byte) error {
	var v []interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	a.Name = v[0].(string)
	a.Id = v[1].(float64)
	return nil
}

type Converter struct {
	*ncm.NcmFile
	KeyData   []byte
	MetaData  *Meta
	MusicData []byte
}

func NewConverter(ncmFile *ncm.NcmFile) *Converter {
	return &Converter{
		NcmFile: ncmFile,
	}
}

func (c *Converter) HandleKey() error {
	tmp := make([]byte, c.Key.Length)
	for i := range c.Key.Detail {
		tmp[i] = c.Key.Detail[i] ^ 0x64
	}

	decodedData, err := decryptAes128(aesCoreKey, tmp)
	if err != nil {
		return err
	}
	// 17 = len("neteasecloudmusic")
	c.KeyData = decodedData
	return nil
}

func (c *Converter) HandleMeta() error {
	if c.Meta.Length <= 0 {
		format := "flac"
		if info, err := c.GetFDStat(); err != nil && info.Size() < int64(math.Pow(1024, 2)*16) {
			format = "map3"
		}
		c.MetaData = &Meta{
			Format: format,
		}
		return nil
	}

	tmp := make([]byte, c.Meta.Length)
	for i := range c.Meta.Detail {
		tmp[i] = c.Meta.Detail[i] ^ 0x63
	}

	// 22 = len("163 key(Don't modify):")
	decodedModifyData := make([]byte, base64.StdEncoding.DecodedLen(int(c.Meta.Length)-22))
	if _, err := base64.StdEncoding.Decode(decodedModifyData, tmp[22:]); err != nil {
		c.MetaData = &Meta{}
		return err
	}

	decodedData, err := decryptAes128(aesModifyKey, decodedModifyData)
	if err != nil {
		c.MetaData = &Meta{}
		return err
	}

	// 6 = len("music:")
	decodedData = decodedData[6:]

	var album Album
	if err := json.Unmarshal(decodedData, &album); err != nil {
		c.MetaData = &Meta{}
		return err
	}

	var meta Meta
	if err := json.Unmarshal(decodedData, &meta); err != nil {
		c.MetaData = &Meta{}
		return err
	}

	meta.Album = &album
	meta.Comment = string(tmp)
	c.MetaData = &meta
	return nil
}

func (c *Converter) HandleMusic() error {
	if c.KeyData == nil {
		var once = sync.Once{}
		once.Do(func() {
			_ = c.HandleKey()
		})
	}
	box := buildKeyBox(c.KeyData[17:])
	n := 0x8000
	var writer bytes.Buffer

	data := bytes.NewReader(c.Music.Detail)

	var tb = make([]byte, n)
	for {
		if _, err := data.Read(tb); err != nil {
			break // read EOF
		}

		for i := 0; i < n; i++ {
			j := byte((i + 1) & 0xff)
			tb[i] ^= box[(box[j]+box[(box[j]+j)&0xff])&0xff]
		}

		writer.Write(tb) // write to memory
	}

	c.MusicData = writer.Bytes()
	return nil
}

func (c *Converter) HandleAll() error {
	err := c.HandleKey()
	if err != nil {
		return err
	}

	err = c.HandleMeta()
	if err != nil {
		return err
	}

	err = c.HandleMusic()
	if err != nil {
		return err
	}

	return nil
}
