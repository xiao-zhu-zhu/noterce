package mode

import "time"

type File struct {
	Path  string    `json:"path,omitempty"`
	Size  int64     `json:"size,omitempty"`
	Name  string    `json:"name,omitempty"`
	IsDir bool      `json:"isDir,omitempty"`
	Mode  string    `json:"mode,omitempty"`
	Time  time.Time `json:"time,omitempty"`
}
