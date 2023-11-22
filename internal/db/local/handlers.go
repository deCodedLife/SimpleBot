package local

import (
	"bot_lab/internal/db"
	"bot_lab/internal/model"
	"encoding/json"
	"io"
	"log"
	"os"
)

func (db *Database) NewConnection(c db.Credentials) error {
	var err error = nil
	db.File, err = os.OpenFile(c.FileName, os.O_RDWR|os.O_CREATE, 0600)
	db.loadScheme()
	return err
}

func (db *Database) loadScheme() {
	buff, err := io.ReadAll(db.File)

	if err != nil {
		log.Fatalf("Невозможно открыть базу данных ", err.Error())
	}

	if len(buff) == 0 {
		return
	}

	err = json.Unmarshal(buff, &db)

	if err != nil {
		log.Fatalf("Невозможно прочесть данные ", err.Error())
	}

}

func (db *Database) saveDatabase() {
	_ = db.File.Truncate(0)
	_, _ = db.File.Seek(0, 0)

	buff, err := json.Marshal(db)

	if err != nil {
		log.Panicln("Невозможно сформировать записи базы данных ", err.Error())
	}

	_, err = db.File.Write(buff)

	if err != nil {
		log.Fatalf("Невозможно записать данные в базу данных ", err.Error())
	}
}

func (db *Database) DropConnection() {
	db.saveDatabase()
	err := db.File.Close()
	if err != nil {
		log.Println("Невозможно закрыть файл: ", err.Error())
	}
}

func (db *Database) SaveMessage(m model.Message) error {
	db.Messages = append(db.Messages, m)
	db.saveDatabase()
	return nil
}

func (db *Database) GetMessages() []model.Message {
	return db.Messages
}
