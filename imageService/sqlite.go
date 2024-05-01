package main
//
// import (
// 	"database/sql"
// 	"log"
// )
//
// type SqliteStorage struct {
// 	db *sql.DB
// }
//
// func NewSqliteStorage() *SqliteStorage {
//
// 	db, err := sql.Open("sqlite3", "/data/imageService.db")
// 	if err != nil {
// 		log.Panic(err)
// 	}
//
// 	if err := db.Ping(); err != nil {
// 		log.Fatal(err)
// 	}
//
// 	return &SqliteStorage{
// 		db: db,
// 	}
// }
//
//
// func (s *SqliteStorage)Init(){
//     
//     if err := createImage
// }
