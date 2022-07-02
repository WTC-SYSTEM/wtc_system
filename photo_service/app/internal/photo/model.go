package photo

import (
	"fmt"
	"sync"
)

type PhotoDTO struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Bytes    []byte `json:"bytes"`
}

func (p *PhotoDTO) String() string {
	return fmt.Sprintf("%s (%d size)", p.Filename, p.Size)
}

type UploadDTO struct {
	Photos []*PhotoDTO `json:"photos"`
	Mutex  sync.Mutex
}
