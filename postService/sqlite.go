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

	db, err := sql.Open("sqlite3", "data/postService.db")
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

	if err := s.createPostTable(); err != nil {
		log.Fatal(err)
	}

}

func (s *SqliteStorage) createPostTable() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS posts (
            id TEXT PRIMARY KEY,
            image TEXT,
            body TEXT NOT NULL,
            idUser TEXT NOT NULL,
            username TEXT NOT NULL,
            name TEXT NOT NULL,
            profile TEXT NOT NULL,
            totalLikes INTEGER DEFAULT 0 NOT NULL,
            totalReplies INTEGER DEFAULT 0 NOT NULL,
            
            createdAt INTEGER NOT NULL,
            updatedAt INTEGER NOT NULL,
            deletedAt INTEGER
        )`)

	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) UpdatePostBody(id, body, userid string) error {

	stmt, err := s.db.Prepare(`
        UPDATE posts
        SET
            body = ?,
            updatedAt = ?
        WHERE 
            id = ?
            AND idUser = ?
            AND deletedAt IS NULL
        `)
	if err != nil {
		return err
	}

	unixEpoch := time.Now().Unix()

	_, err = stmt.Exec(body, unixEpoch, id, userid)
	if err != nil {
		return err
	}

	return nil
}

// listPostByUser --> nampilin list post yang dibuat oleh user
func (s *SqliteStorage) ListPostByUser(cursor int64, userId string, posts *[]Post) error {
	stmt, err := s.listPostByUserStatement(cursor, userId)
	if err != nil {
		log.Println("Stmt error:", err)
		return err
	}

	defer stmt.Close()

	rows, err := stmt.Query(userId, cursor)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.Id,
			&post.Image,
			&post.Body,
			&post.IdUser,
			&post.Username,
			&post.Name,
			&post.TotalLikes,
			&post.TotalReplies,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return err
		}
		*posts = append(*posts, post)

	}

	return nil

}

// listPosts --> nampilin list post
func (s *SqliteStorage) ListPost(cursor int64, posts *[]Post) error {
	stmt, err := s.listPostStatement(cursor)
	defer stmt.Close()
	if err != nil {
		log.Println("Error when creating stmt in listPosts", err)
	}

	rows, err := stmt.Query(cursor)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.Id,
			&post.Image,
			&post.Body,
			&post.IdUser,
			&post.Username,
			&post.Name,
			&post.TotalLikes,
			&post.TotalReplies,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return err
		}
		*posts = append(*posts, post)

	}

	return nil
}

// getPostById --> nampilin satu post
func (s *SqliteStorage) GetPostById(id string, post *Post) error {
	stmt, err := s.db.Prepare(`
        SELECT
            id,
            image,
            body,
            idUser,
            username,
            name,
            totalLikes,
            totalReplies,
            createdAt,
            updatedAt
        FROM
            posts 
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
		&post.Id,
		&post.Image,
		&post.Body,
		&post.IdUser,
		&post.Username,
		&post.Name,
		&post.TotalLikes,
		&post.TotalReplies,
		&post.CreatedAt,
		&post.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) DeletePostById(id, userId string) error {

	stmt, err := s.db.Prepare(`
        UPDATE 
            posts
        SET 
            deletedAt = ?
        WHERE 
            id = ?
            AND idUser = ?
            AND deletedAt IS NULL
        `)
	if err != nil {
		return err
	}

	defer stmt.Close()

	unixEpoch := time.Now().Unix()

	if _, err := stmt.Exec(unixEpoch, id, userId); err != nil {
		return err
	}

	return nil
}

// func (s *SqliteStorage) UpdatePostUsername(id, username string) error {
// 	stmt, err := s.db.Prepare(`
//         UPDATE posts
//         SET
//             username = ?,
//             updatedAt = ?
//         WHERE
//             id = ?
//             AND deletedAt IS NULL
//         `)
//
// 	defer stmt.Close()
//
// 	if err != nil {
// 		return err
// 	}
//
// 	unixEpoch := time.Now().Unix()
//
// 	if _, err := stmt.Exec(username, unixEpoch, id); err != nil {
// 		return err
// 	}
//
// 	return nil
// }

func (s *SqliteStorage) UpdateUserDetail(idUser, profile, name string) error {
	stmt, err := s.db.Prepare(`
        UPDATE posts
        SET
            name = ?,
            profile = ?,
            updatedAt = ?
        WHERE 
            idUser = ?
            AND deletedAt IS NULL
        `)

	defer stmt.Close()

	if err != nil {
		return err
	}

	unixEpoch := time.Now().Unix()

	if _, err := stmt.Exec(name, profile, unixEpoch, idUser); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) CreatePost(id, image, body, idUser, username, name, profile string) error {
	stmt, err := s.db.Prepare(`
        INSERT INTO posts (
            id,
            image,
            body,
            idUser,
            username,
            name,
            profile,
            
            createdAt,
            updatedAt
        ) VALUES (?,?,?,?,?,?,?,?,?)
        `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	unixEpoch := time.Now().Unix()

	if _, err := stmt.Exec(
		id,
		image,
		body,
		idUser,
		username,
		name,
		profile,
		unixEpoch,
		unixEpoch); err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) listPostByUserStatement(cursor int64, idUser string) (*sql.Stmt, error) {
	var queryStr string
	if cursor == 0 {
		queryStr = `
        SELECT
            id,
            image,
            body,
            idUser,
            username,
            name,
            totalLikes,
            totalReplies,
            createdAt,
            updatedAt
        FROM
            posts 
        WHERE
            idUser = ?
            AND deletedAt IS NULL
        ORDER BY
            createdAt DESC
        LIMIT 10`

		return s.db.Prepare(queryStr)
	} else {
		queryStr = `
        SELECT
            id,
            image,
            body,
            idUser,
            username,
            name,
            totalLikes,
            totalReplies,
            createdAt,
            updatedAt
        FROM
            posts 
        WHERE
            idUser = ?
            AND deletedAt IS NULL
            AND createdAt < ?
        ORDER BY
            createdAt DESC
        LIMIT 10`

		return s.db.Prepare(queryStr)
	}
}

func (s *SqliteStorage) listPostStatement(cursor int64) (*sql.Stmt, error) {
	var queryStr string
	if cursor == 0 {
		queryStr = `
        SELECT
            id,
            image,
            body,
            idUser,
            username,
            name,
            totalLikes,
            totalReplies,
            createdAt,
            updatedAt
        FROM
            posts 
        WHERE
            deletedAt IS NULL
        ORDER BY
            createdAt DESC
        LIMIT 10
        `
		return s.db.Prepare(queryStr)
	} else {
		queryStr = `
        SELECT
            id,
            image,
            body,
            idUser,
            username,
            name,
            totalLikes,
            totalReplies,
            createdAt,
            updatedAt
        FROM
            posts 
        WHERE
            deletedAt IS NULL
            AND createdAt < ?
        ORDER BY
            createdAt DESC
        LIMIT 10
        `
		return s.db.Prepare(queryStr)
	}
}

func (s *SqliteStorage) setPragmaWal() error {
	_, err := s.db.Exec(`PRAGMA journal_mode=WAL;`)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
