package ncm

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/closetool/NCMConverter/path"
)

const (
	MagicHaeder1 = 0x4e455443
	MagicHaeder2 = 0x4d414446
)

type Data struct {
	Length uint64
	Detail []byte
}

type NcmFile struct {
	Path     string
	FileDir  string
	FileName string
	fd       *os.File
	Ext      string
	valid    bool
	Key      Data
	Meta     Data
	Cover    Data
	Music    Data
}

func NewNcmFile(ncmpath string) (nf *NcmFile, err error) {
	nf = new(NcmFile)
	//if runtime.GOOS == "windows" {
	//	ncmpath = filepath.Clean(ncmpath)
	//	nf.FileName = filepath.Base(ncmpath)
	//	nf.FileDir = filepath.Dir(ncmpath)
	//	nf.Ext = filepath.Ext(ncmpath)
	//} else if (runtime.GOOS) == "linux" {
	//	ncmpath = path.Clean(ncmpath)
	//	nf.FileName = path.Base(ncmpath)
	//	nf.FileDir = path.Dir(ncmpath)
	//	nf.Ext = path.Ext(ncmpath)
	//}
	ncmpath = path.Clean(ncmpath)
	nf.FileName = path.Base(ncmpath)
	nf.FileDir = path.Dir(ncmpath)
	nf.Ext = path.Ext(ncmpath)
	nf.Path = ncmpath
	if nf.fd, err = os.Open(nf.Path); err != nil {
		return nil, err
	}
	return
}

func (nf *NcmFile) Validate() error {
	if !strings.EqualFold(nf.Ext, ".ncm") {
		return ErrExtNcm
	}

	err := nf.CheckHaeder()

	if err != nil {
		return err
	}

	nf.valid = true
	return nil
}

func (nf *NcmFile) CheckHaeder() error {
	if _, err := nf.fd.Seek(0, io.SeekStart); err != nil {
		return err
	}

	m1, err := readUint32(nf.fd)
	m2, err := readUint32(nf.fd)
	if err != nil {
		return err
	}
	if m1 != MagicHaeder1 && m2 != MagicHaeder2 {
		return ErrMagicHeader
	}
	return nil
}

func (nf *NcmFile) getData(offset int64) ([]byte, uint32, error) {
	if _, err := nf.fd.Seek(offset, io.SeekStart); err != nil {
		return nil, 0, err
	}
	length, err := readUint32(nf.fd)
	if err != nil {
		return nil, 0, err
	}

	buf := make([]byte, length)

	if _, err = nf.fd.Read(buf); err != nil {
		return nil, 0, err
	}
	return buf, length, nil
}

func (nf *NcmFile) GetKey() (err error) {
	tmp, length, err := nf.getData(4*2 + 2)
	if err != nil {
		return err
	}
	nf.Key.Length = uint64(length)
	nf.Key.Detail = tmp
	return nil
}

func (nf *NcmFile) GetMeta() (err error) {
	tmp, length, err := nf.getData(int64(4*2 + 2 + 4 + nf.Key.Length))
	if err != nil {
		return
	}
	nf.Meta.Detail = tmp
	nf.Meta.Length = uint64(length)
	return nil
}

func (nf *NcmFile) GetCover() (err error) {
	tmp, length, err := nf.getData(int64(4*2 + 2 + 4 + nf.Key.Length + 4 + nf.Meta.Length + 5 + 4))
	if err != nil {
		return
	}

	nf.Cover.Detail = tmp
	nf.Cover.Length = uint64(length)
	return nil
}

func (nf *NcmFile) GetMusicData() error {
	nf.fd.Seek(int64(4*2+2+4+nf.Key.Length+4+nf.Meta.Length+9+4+nf.Cover.Length), io.SeekStart)
	file := nf.fd
	buf := make([]byte, 1024)
	nf.Music.Detail = make([]byte, 0)
	var (
		length int
		err    error
	)
	for {
		if length, err = file.Read(buf); err != nil && err != io.EOF {
			return err
		}
		nf.Music.Detail = append(nf.Music.Detail, buf[:length]...)
		nf.Music.Length += uint64(length)
		if err == io.EOF {
			return nil
		}
	}
}

func (nf *NcmFile) Close() error {
	return nf.fd.Close()
}

func (nf *NcmFile) Parse() error {
	err := nf.Validate()
	if err != nil {
		log.Printf("Ncm magic header check failed: %v", err)
		return err
	}

	err = nf.GetKey()
	if err != nil {
		log.Printf("Get Key Failed: %v", err)
		return err
	}

	err = nf.GetMeta()
	if err != nil {
		log.Printf("Get Meta Failed: %v", err)
		return err
	}

	err = nf.GetCover()
	if err != nil {
		log.Printf("Get Cover Failed: %v", err)
		return err
	}

	err = nf.GetMusicData()
	if err != nil {
		log.Printf("Get Music Data Failed: %v", err)
		return err
	}

	return nil
}

func (nf *NcmFile) GetFDStat() (os.FileInfo, error) {
	return nf.fd.Stat()
}
