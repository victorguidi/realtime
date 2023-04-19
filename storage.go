package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	CreateUser(*User) error
	CreateSession(*Session) error
	UpdateUser(int) error
	UpdateSession(int) error
	GetUser(int) (*User, error)
	GetSession(int) (*Session, error)
	GetAllUsers() ([]*User, error)
	GetAllSessions() ([]*Session, error)
	DeleteUser(int) error
	DeleteSession(int) error
}

type Database struct {
	db *sql.DB
}

func NewDatabase(file string) (*Database, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil

}

func (d *Database) Init() error {
	if err := d.createTables(); err != nil {
		return err
	}
	return nil
}

func (d *Database) createTables() error {

	queryUser := "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username VARCHAR(100), password VARCHAR(100), created_at DATETIME DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME DEFAULT CURRENT_TIMESTAMP)"
	querySession := "CREATE TABLE IF NOT EXISTS sessions (id INTEGER PRIMARY KEY AUTOINCREMENT, sessionName VARCHAR(100), sessionToken VARCHAR(100), created_at DATETIME DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME DEFAULT CURRENT_TIMESTAMP)"
	queryUserSession := "CREATE TABLE IF NOT EXISTS user_session (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER, session_id INTEGER, created_at DATETIME DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME DEFAULT CURRENT_TIMESTAMP)"

	if _, err := d.db.Exec(queryUser); err != nil {
		return err
	}
	if _, err := d.db.Exec(querySession); err != nil {
		return err
	}
	if _, err := d.db.Exec(queryUserSession); err != nil {
		return err
	}

	return nil
}

func (d *Database) GetAllUsers() ([]*User, error) {
	query := "SELECT * FROM users"
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user := new(User)
		err := rows.Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (d *Database) GetAllSessions() ([]*Session, error) {
	query := "SELECT * FROM sessions"
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]*Session, 0)
	for rows.Next() {
		session := new(Session)
		err := rows.Scan(&session.ID, &session.SessionName, &session.SessionToken, &session.Created_at, &session.Updated_at)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (d *Database) CreateUser(u *User) error {
	return nil
}

func (d *Database) CreateSession(s *Session) error {
	return nil
}

func (d *Database) UpdateUser(id int) error {
	return nil
}

func (d *Database) UpdateSession(id int) error {
	return nil
}

func (d *Database) GetUser(id int) (*User, error) {
	return nil, nil
}

func (d *Database) GetSession(id int) (*Session, error) {
	return nil, nil
}

func (d *Database) DeleteUser(id int) error {
	return nil
}

func (d *Database) DeleteSession(id int) error {
	return nil
}
