package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteStorage struct {
	db *sql.DB
}

func NewSqliteStorage() *SqliteStorage {

	db, err := sql.Open("sqlite3", "data/userService.db")
	if err != nil {
		log.Panic("panic:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("cannot ping:", err)
	}

	return &SqliteStorage{
		db: db,
	}
}

func (s *SqliteStorage) Init() {

	if err := s.setPragmaWal(); err != nil {
		log.Fatal(err)
	}

	if err := s.createUserTable(); err != nil {
		log.Fatal(err)
	}

}

func (s *SqliteStorage) createUserTable() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id TEXT PRIMARY KEY,
            username TEXT UNIQUE NOT NULL,
            name TEXT NOT NULL,
            hashPassword TEXT NOT NULL,
            profile TEXT NOT NULL,
            
            createdAt INTEGER NOT NULL,
            updatedAt INTEGER NOT NULL,
            deletedAt INTEGER
        )`)

	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) CreateUser(id, username, name, hashPassword, profile string, createdAt, updatedAt int64) error {

	stmt, err := s.db.Prepare(`
        INSERT INTO users (
        id,
        username,
        name,
        hashPassword,
        profile,
        createdAt,
        updatedAt
        ) VALUES (?,?,?,?,?,?,?);
        `)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err := stmt.Exec(
		id,
		username,
		name,
		hashPassword,
		profile,
		createdAt,
		updatedAt); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) UpdateUserPasswordById(newPassword, id string) error {

	stmt, err := s.db.Prepare(`
        UPDATE users
        SET hashPassword = ?
        SET updatedAt = ?
        WHERE id = ?`)

	defer stmt.Close()

	if err != nil {
		return err
	}

	unixEpoch := time.Now().Unix()

	if _, err := stmt.Exec(newPassword, unixEpoch, id); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) DeleteUserById(id string) error {

	stmt, err := s.db.Prepare(`
        UPDATE users
        SET deletedAt = ?
        WHERE id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	unixEpoch := time.Now().Unix()

	if _, err := stmt.Exec(unixEpoch, id); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) UpdateUserNameById(name string, id string) error {

	stmt, err := s.db.Prepare(`
        UPDATE users
        SET name = ?
        SET updatedAt = ?
        WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	unixEpoch := time.Now().Unix()

	if _, err := stmt.Exec(name, unixEpoch, id); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) GetUserByUsername(username string, user *ReturnUser) error {

	stmt, err := s.db.Prepare(`
        SELECT 
        id,
        username,
        name,
        hashPassword,
        profile,
        createdAt,
        deletedAt 
        FROM users WHERE username = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if err := stmt.QueryRow(username).Scan(
		&user.Id,
		&user.Username,
		&user.Name,
		&user.HashPassword,
		&user.Profile,
		&user.CreatedAt,
		&user.DeletedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) GetUserById(id string, user *ReturnUser) error {

	stmt, err := s.db.Prepare(`
        SELECT 
        id,
        username,
        name,
        hashPassword,
        profile,
        createdAt,
        deletedAt 
        FROM users WHERE id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if err := stmt.QueryRow(id).Scan(
		&user.Id,
		&user.Username,
		&user.Name,
		&user.HashPassword,
		&user.Profile,
		&user.CreatedAt,
		&user.DeletedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) setPragmaWal() error {
	_, err := s.db.Exec(`PRAGMA journal_mode=WAL;`)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
