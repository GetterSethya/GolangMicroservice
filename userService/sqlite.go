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
            totalFollower INTEGER DEFAULT 0 NOT NULL,
            totalFollowing INTEGER DEFAULT 0 NOT NULL,
            
            createdAt INTEGER NOT NULL,
            updatedAt INTEGER NOT NULL,
            deletedAt INTEGER
        )`)

	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) UpdateProfileById(profileUrl, id string) error {

	stmt, err := s.db.Prepare(`
        UPDATE users
        SET
            profile = ?
        WHERE
            id = ?
        `)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(profileUrl, id); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) IncrementFollowerById(id string) error {

	stmt, err := s.db.Prepare(`
        UPDATE users
        SET
            totalFollower = totalFollower + 1
        WHERE
            id = ?
        `)

	if err != nil {
		return err
	}

	if _, err := stmt.Exec(id); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) DecrementFollowerById(id string) error {

	stmt, err := s.db.Prepare(`
        UPDATE users
        SET
            totalFollower = totalFollower - 1
        WHERE
            id = ?
        `)

	if err != nil {
		return err
	}

	if _, err := stmt.Exec(id); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) IncrementFollowingById(id string) error {
	stmt, err := s.db.Prepare(`
        UPDATE users
        SET
            totalFollowing = totalFollowing + 1
        WHERE
            id = ?
        `)

	if err != nil {
		return err
	}

	if _, err := stmt.Exec(id); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) DecrementFollowingById(id string) error {
	stmt, err := s.db.Prepare(`
        UPDATE users
        SET
            totalFollowing = totalFollowing - 1
        WHERE
            id = ?
        `)

	if err != nil {
		return err
	}

	if _, err := stmt.Exec(id); err != nil {
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

func (s *SqliteStorage) GetUserPasswordByUsername(username string, user *User) error {

	stmt, err := s.db.Prepare(`
        SELECT 
        id,
        username,
        hashPassword,
        createdAt,
        updatedAt 
        FROM users WHERE username = ? AND deletedAt IS NULL`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if err := stmt.QueryRow(username).Scan(
		&user.Id,
		&user.Username,
		&user.HashPassword,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return err
	}
	return nil
}

func (s *SqliteStorage) GetUserPasswordById(id string, user *User) error {

	stmt, err := s.db.Prepare(`
        SELECT 
        id,
        hashPassword,
        createdAt,
        updatedAt 
        FROM users WHERE id = ? AND deletedAt IS NULL`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if err := stmt.QueryRow(id).Scan(
		&user.Id,
		&user.HashPassword,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) UpdateUserPasswordById(newPassword, id string) error {

	stmt, err := s.db.Prepare(`
        UPDATE users
        SET 
            hashPassword = ?,
            updatedAt = ?
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

func (s *SqliteStorage) UpdateUserNameAndProfile(name, profile, id string) error {

	stmt, err := s.db.Prepare(`
        UPDATE users
        SET 
            name = ?,
            profile = ?,
            updatedAt = ?
        WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	unixEpoch := time.Now().Unix()

	if _, err := stmt.Exec(name, profile, unixEpoch, id); err != nil {
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
        profile,
        createdAt,
        updatedAt 
        FROM users WHERE username = ? AND deletedAt IS NULL`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if err := stmt.QueryRow(username).Scan(
		&user.Id,
		&user.Username,
		&user.Name,
		&user.Profile,
		&user.CreatedAt,
		&user.UpdatedAt,
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
        profile,
        createdAt,
        updatedAt
        FROM users WHERE id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if err := stmt.QueryRow(id).Scan(
		&user.Id,
		&user.Username,
		&user.Name,
		&user.Profile,
		&user.CreatedAt,
		&user.UpdatedAt,
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
