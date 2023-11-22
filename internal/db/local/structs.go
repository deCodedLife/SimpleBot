package local

import (
	"bot_lab/internal/model"
	"os"
)

type Database struct {
	FileName string          `json:"-"`
	File     *os.File        `json:"-"`
	Messages []model.Message `json:"messages"`
}
