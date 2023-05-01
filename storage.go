package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	CreateUser(*User) error
	CreateSession(*Session) error
	UpdateUser(int) error
	UpdateSession(int) error
	GetUser(int) (*User, error)
	GetUserByName(string) (*User, error)
	GetSession(int) (*Session, error)
	GetAllUsers() ([]*User, error)
	GetAllSessions() ([]*Session, error)
	DeleteUser(int) error
	DeleteSession(int) error
	AddSessionToUser(int, int) error
	RemoveSessionFromUser(int, int) error
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
	query := "SELECT u.id, u.username, u.password, us.session_id FROM users u INNER JOIN user_session us ON u.id = us.user_id"
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type Response struct {
		userId     int
		username   string
		password   string
		sesssionId int
	}

	users := make([]*User, 0)
	for rows.Next() {
		resp := new(Response)
		err := rows.Scan(&resp.userId, &resp.username, &resp.password, &resp.sesssionId)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		user := User{
			ID:       resp.userId,
			Username: resp.username,
			Password: resp.password,
			Sessions: make([]int, 0),
		}
		user.Sessions = append(user.Sessions, resp.sesssionId)
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (d *Database) GetAllSessions() ([]*Session, error) {
	query := "SELECT id, sessionName, created_at, updated_at FROM sessions"
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]*Session, 0)
	for rows.Next() {
		session := new(Session)
		err := rows.Scan(&session.ID, &session.SessionName, &session.Created_at, &session.Updated_at)
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

func (d *Database) AddSessionToUser(userId, sessionId int) error {

	query := "INSERT INTO user_session (user_id, session_id, created_at, updated_at) SELECT ?, ?, ?, ? WHERE NOT EXISTS (SELECT id FROM user_session WHERE user_id = ? AND session_id = ?)"
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return err
	}

	current_timestamp := time.Now()

	_, err = stmt.Exec(userId, sessionId, current_timestamp, current_timestamp, userId, sessionId)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) RemoveSessionFromUser(userId, sessionId int) error {
	return nil
}

func (d *Database) CreateUser(u *User) error {
	current_timestamp := time.Now()

	query := "INSERT INTO users (username, password, created_at, updated_at) VALUES(?, ?, ?, ?)"
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(u.Username, u.Password, current_timestamp, current_timestamp)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) CreateSession(s *Session) error {
	current_timestamp := time.Now()

	query := "INSERT INTO sessions (sessionName, sessionToken, created_at, updated_at) VALUES(?, ?, ?, ?)"
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(s.SessionName, s.SessionToken, current_timestamp, current_timestamp)
	if err != nil {
		return err
	}

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
func (d *Database) GetUserByName(username string) (*User, error) {

	query := "SELECT id, username, password FROM users WHERE username = ?"
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	user := new(User)
	err = stmt.QueryRow(username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *Database) GetSession(id int) (*Session, error) {

	query := "SELECT id, sessionToken FROM sessions WHERE id = ?"
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	session := new(Session)
	err = stmt.QueryRow(id).Scan(&session.ID, &session.SessionToken)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (d *Database) DeleteUser(id int) error {
	return nil
}

func (d *Database) DeleteSession(id int) error {
	return nil
}
