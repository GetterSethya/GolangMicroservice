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
	db, err := sql.Open("sqlite3", "data/replyService.db")
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

	if err := s.createReplyTable(); err != nil {
		log.Fatal(err)
	}
}

func (s *SqliteStorage) createReplyTable() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS replies (
            id TEXT PRIMARY KEY,
            body TEXT NOT NULL,
            idUser TEXT NOT NULL,
            username TEXT NOT NULL,
            name TEXT NOT NULL,
            profile TEXT NOT NULL,
            totalChild INTEGER DEFAULT 0 NOT NULL ,
            idPost TEXT NOT NULL,
            parentId TEXT,
            
            createdAt INTEGER NOT NULL,
            updatedAt INTEGER NOT NULL,
            deletedAt INTEGER
        )`)
	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) DecrementReplyChildCount(parentId string) error {
	query := `
    UPDATE replies
    SET totalChild = totalChild - 1
    WHERE parentId = ?
    `
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err := stmt.Exec(parentId); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) IncrementReplyChildCount(parentId string) error {
	query := `
    UPDATE replies
    SET totalChild = totalChild + 1
    WHERE parentId = ?
    `
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err := stmt.Exec(parentId); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) CreateReply(id, body, idUser, username, name, profile, idPost string, parentId interface{}) error {
	query := `
    INSERT INTO replies (
        id,
        body,
        idUser,
        username,
        name,
        profile,
        idPost,
        parentId,
        createdAt,
        updatedAt
    ) VALUES (?,?,?,?,?,?,?,?,?,?);
    `
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	unixEpoch := time.Now().Unix()

	defer stmt.Close()
	if _, err := stmt.Exec(
		id, body, idUser, username, name, profile, idPost, parentId, unixEpoch, unixEpoch,
	); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) GetReplyByParentId(parentId string, cursor int64, limit int32, replies *[]Reply) error {
	stmt, err := s.db.Prepare(`
        SELECT
            id,
            body,
            idUser,
            username,
            name,
            profile,
            idPost,
            parentId,
            createdAt,
            updatedAt
        FROM
            replies
        WHERE
            deletedAt IS NULL
            AND parentId = ?
            AND createdAt < ?
        ORDER BY
            createdAT DESC
        LIMIT ?
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if cursor == 0 {
		cursor = 922337203685477
	}

	rows, err := stmt.Query(parentId, cursor, limit)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var reply Reply
		err := rows.Scan(
			&reply.Id,
			&reply.Body,
			&reply.IdUser,
			&reply.Username,
			&reply.Name,
			&reply.Profile,
			&reply.IdPost,
			&reply.ParentId,
			&reply.CreatedAt,
			&reply.UpdatedAt,
		)
		if err != nil {
			return err
		}

		*replies = append(*replies, reply)
	}
	return nil
}

func (s *SqliteStorage) GetReplyByPostId(postId string, cursor int64, limit int32, replies *[]Reply) error {
	stmt, err := s.db.Prepare(`
        SELECT
            id,
            body,
            idUser,
            username,
            name,
            profile,
            idPost,
            parentId,
            createdAt,
            updatedAt
        FROM
            replies
        WHERE
            deletedAt IS NULL
            AND idPost = ?
            AND createdAt < ?
        ORDER BY
            createdAT DESC
        LIMIT ?
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if cursor == 0 {
		cursor = 922337203685477
	}

	rows, err := stmt.Query(postId, cursor, limit)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var reply Reply
		err := rows.Scan(
			&reply.Id,
			&reply.Body,
			&reply.IdUser,
			&reply.Username,
			&reply.Name,
			&reply.Profile,
			&reply.IdPost,
			&reply.ParentId,
			&reply.CreatedAt,
			&reply.UpdatedAt,
		)
		if err != nil {
			return err
		}

		*replies = append(*replies, reply)
	}

	return nil
}

func (s *SqliteStorage) GetReplyById(id string, reply *Reply) error {
	stmt, err := s.db.Prepare(`
        SELECT
            id,
            body,
            idUser,
            username,
            name,
            profile,
            idPost,
            parentId,
            createdAt,
            updatedAt
        FROM
            replies
        WHERE
            id = ?
            AND deletedAt IS NULL
        LIMIT 1
    `)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if err := stmt.QueryRow(id).Scan(
		&reply.Id,
		&reply.Body,
		&reply.IdUser,
		&reply.Username,
		&reply.Name,
		&reply.Profile,
		&reply.IdPost,
		&reply.ParentId,
		&reply.CreatedAt,
		&reply.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) DeleteReply(id string) error {
	stmt, err := s.db.Prepare(`
        UPDATE replies
        SET
            deletedAt = ?
        WHERE
            id = ?
            AND deletedAt is NULL
    `)
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

func (s *SqliteStorage) UpdateReply(id, newReplyBody string) error {
	stmt, err := s.db.Prepare(`
        UPDATE replies
        SET
            body = ?
            updatedAt = ?
        WHERE
            id = ?
            AND deletedAt IS NULL
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	unixEpoch := time.Now().Unix()

	if _, err := stmt.Exec(newReplyBody, unixEpoch, id); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) UpdateUserDetail(idUser, profile, name string) error {
	stmt, err := s.db.Prepare(`
        UPDATE replies
        SET
            name = ?,
            profile = ?,
            updatedAt = ?
        WHERE 
            idUser = ?
            AND deletedAt IS NULL
    `)
	if err != nil {
		return err
	}

	defer stmt.Close()

	unixEpoch := time.Now().Unix()

	if _, err := stmt.Exec(name, profile, unixEpoch, idUser); err != nil {
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
