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
	db, err := sql.Open("sqlite3", "data/relationService.db")
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

	if err := s.createRelationTable(); err != nil {
		log.Fatal(err)
	}

	if err := s.createUserTable(); err != nil {
		log.Fatal(err)
	}
}

func (s *SqliteStorage) createRelationTable() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS relations (
            id TEXT PRIMARY KEY,
            targetId TEXT NOT NULL,
            followerId TEXT NOT NULL,

            createdAt INTEGER NOT NULL,
            updatedAt INTEGER NOT NULL,
            deletedAt INTEGER
        )`)
	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) createUserTable() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id TEXT PRIMARY KEY,
            username TEXT UNIQUE NOT NULL,
            name TEXT NOT NULL,
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

// get followers
func (s *SqliteStorage) GetFollowers(cursor int64, userId string, limit int, followers *[]User) error {
	if cursor == 0 {
		cursor = 922337203685477
	}

	stmt, err := s.db.Prepare(`
        SELECT 
            users.id, users.username, users.profile
        FROM
            relations
        INNER JOIN users ON relation.followerId = users.id
        WHERE
            relations.targetId = ?
        AND
            relations.createdAt < ?
        AND 
            relations.deletedAt IS NULL
        ORDER BY 
            relations.createdAt DESC
        LIMIT ?
    `)
	if err != nil {
		log.Println("stmt error", err)
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId, cursor, limit)
	if err != nil {
		log.Println("query error", err)
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var follower User
		err := rows.Scan(
			&follower.Id,
			&follower.Name,
			&follower.Username,
			&follower.Profile,
		)
		if err != nil {
			return err
		}
		*followers = append(*followers, follower)
	}

	return nil
}

// get following
func (s *SqliteStorage) GetFollowing(cursor int64, userId string, limit int, following *[]User) error {
	if cursor == 0 {
		cursor = 922337203685477
	}

	stmt, err := s.db.Prepare(`
        SELECT 
            users.id, users.username, users.profile
        FROM
            relations
        INNER JOIN users ON relation.targetId = users.id
        WHERE
            relations.followerId = ?
        AND 
            relations.createdAt < ?
        AND
            relations.deletedAt IS NULL
        LIMIT ?
    `)
	if err != nil {
		log.Println("stmt error", err)
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId, cursor, limit)
	if err != nil {
		log.Println("query error", err)
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var follower User
		err := rows.Scan(
			&follower.Id,
			&follower.Name,
			&follower.Username,
			&follower.Profile,
		)
		if err != nil {
			return err
		}
		*following = append(*following, follower)
	}

	return nil
}

// create relation
func (s *SqliteStorage) CreateRelation(idRelation, idTarget, idFollower string) error {
	stmt, err := s.db.Prepare(`
        INSERT INTO relations (
            id,
            targetId,
            followerId,

            createdAt,
            updatedAt
        ) VALUES (?,?,?,?,?)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	unixEpoch := time.Now().Unix()

	if _, err := stmt.Exec(idRelation, idTarget, idFollower, unixEpoch, unixEpoch); err != nil {
		return err
	}

	return nil
}

// create user
func (s *SqliteStorage) CreateUser(id, username, name, profile string, createdAt, updatedAt int64) error {
	stmt, err := s.db.Prepare(`
        INSERT INTO users (
            id,
            username,
            name,
            profile,
            createdAt,
            updatedAt
        ) VALUES (?,?,?,?,?,?)
        ON CONFLICT(id) DO UPDATE SET
            username=excluded.username,
            name=excluded.name,
            profile=excluded.profile,
            createdAt=excluded.createdAt,
            updatedAt=excluded.updatedAt,
        ;
    `)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err := stmt.Exec(
		id,
		username,
		name,
		profile,
		createdAt,
		updatedAt); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) IsFollowing(targetId, followerId string) (bool, error) {
	stmt, err := s.db.Prepare(`
        SELECT COUNT(1)
        FROM relations
        WHERE 
            targetId = ?
        AND
            followerId = ?
        AND
            deletedAt IS NULL
    `)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(targetId, followerId).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// delete relation
func (s *SqliteStorage) DeleteRelation(targetId, followerId string) error {
	stmt, err := s.db.Prepare(`
        UPDATE relations
        SET deletedAt = ?
        WHERE 
            targetId = ?
        AND
            followerId = ?
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(time.Now().Unix(), targetId, followerId); err != nil {
		return err
	}

	return nil
}

// delete user
func (s *SqliteStorage) DeleteUser(id string) error {
	// todo
	return nil
}

// update user
func (s *SqliteStorage) UpdateUser(id, name, profile string) error {
	stmt, err := s.db.Prepare(`
        UPDATE users
        SET 
            name = ?,
            profile = ?,
            updatedAt = ?
        WHERE 
            id = ?
        AND
            deletedAt IS NULL
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(name, profile, time.Now().Unix(), id); err != nil {
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
