package out

import (
	"io"
	"os"
)

type Generate struct {
	File *os.File
	TableName string
}

func NewGenerate(file *os.File,tableName string) *Generate {
	return &Generate{
		File: file,
	}
}

func (g *Generate)p(str string)  {
	io.WriteString(g.File,str)
}
