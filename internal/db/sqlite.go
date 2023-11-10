package db

import (
	"database/sql"
	"fmt"

	"github.com/AYM1607/ccclip/pkg/api"
	_ "github.com/mattn/go-sqlite3"
	"github.com/oklog/ulid/v2"
)

type sqliteDB struct {
	internalDB *sql.DB
}

func NewSQLiteDB(location string) DB {
	internalDb, err := sql.Open("sqlite3", fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on", location))
	if err != nil {
		panic(fmt.Sprintf("could not connect to sqlite: %s", err.Error()))
	}
	db := &sqliteDB{
		internalDB: internalDb,
	}
	if err := db.setup(); err != nil {
		panic(fmt.Sprintf("unable to initialize sqlite: %s", err.Error()))
	}
	return db
}

func (d *sqliteDB) PutUser(id string, passwordHash []byte) error {
	_, err := d.internalDB.Exec("INSERT INTO users(id, password_hash) values(?, ?)", id, passwordHash)
	return err
}

func (d *sqliteDB) GetUser(id string) (*api.User, error) {
	res := &api.User{}
	err := d.internalDB.QueryRow("SELECT id, password_hash FROM users WHERE id = ?", id).Scan(&res.ID, &res.PasswordHash)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *sqliteDB) PutDevice(pubKey []byte, userId string) (string, error) {
	id := ulid.Make().String()
	_, err := d.internalDB.Exec("INSERT INTO devices(id, public_key, user_id) values(?, ? ,?)", id, pubKey, userId)
	if err != nil {
		return "", nil
	}
	return id, nil
}

func (d *sqliteDB) GetDevice(id string) (*api.Device, error) {
	res := &api.Device{}
	err := d.internalDB.QueryRow("SELECT id, public_key FROM devices WHERE id = ?", id).Scan(&res.ID, &res.PublicKey)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *sqliteDB) GetUserDevices(userId string) ([]*api.Device, error) {
	res := []*api.Device{}
	rows, err := d.internalDB.Query("SELECT id, public_key FROM devices WHERE user_id = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		device := &api.Device{}
		err := rows.Scan(&device.ID, &device.PublicKey)
		if err != nil {
			return nil, err
		}
		res = append(res, device)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *sqliteDB) GetDeviceUser(deviceId string) (*api.User, error) {
	var userId string
	err := d.internalDB.QueryRow("SELECT user_id FROM devices WHERE id = ?", deviceId).Scan(&userId)
	if err != nil {
		return nil, err
	}

	res := &api.User{}
	err = d.internalDB.QueryRow("SELECT id, password_hash FROM users WHERE id = ?", userId).Scan(&res.ID, &res.PasswordHash)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *sqliteDB) PutClipboard(userId string, clipboard *api.Clipboard) error {
	tx, err := d.internalDB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	clipRes, err := tx.Exec("INSERT INTO clipboards(user_id, sender_public_key) values(?, ?)", userId, clipboard.SenderPublicKey)
	clipId, err := clipRes.LastInsertId()
	if err != nil {
		return nil
	}

	for deviceId, ciphertext := range clipboard.Payloads {
		_, err := tx.Exec("INSERT INTO clipboard_items(ciphertext, clipboard_id, device_id) values(?, ?, ?)", ciphertext, clipId, deviceId)
		if err != nil {
			return nil
		}
	}

	return tx.Commit()
}

func (d *sqliteDB) GetClipboard(userId string) (*api.Clipboard, error) {
	var latestClipId int
	latestClip := &api.Clipboard{}
	err := d.internalDB.QueryRow("SELECT id, sender_public_key FROM clipboards WHERE user_id = ? ORDER BY id DESC LIMIT 1", userId).Scan(&latestClipId, &latestClip.SenderPublicKey)
	if err != nil {
		return nil, err
	}

	rows, err := d.internalDB.Query("SELECT device_id, ciphertext FROM clipboard_items WHERE clipboard_id = ?", latestClipId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	latestClip.Payloads = map[string][]byte{}
	for rows.Next() {
		var deviceId string
		var cipherText []byte
		err := rows.Scan(&deviceId, &cipherText)
		if err != nil {
			return nil, err
		}
		latestClip.Payloads[deviceId] = cipherText
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return latestClip, nil
}

func (d *sqliteDB) setup() error {
	setupStm := `
	CREATE TABLE IF NOT EXISTS users(
		id TEXT PRIMARY KEY,
		password_hash BLOB
	) STRICT;
	CREATE TABLE IF NOT EXISTS devices(
		id TEXT PRIMARY KEY,
		public_key BLOB,
		user_id TEXT,
		FOREIGN KEY(user_id) REFERENCES users(id)
	) STRICT;
	CREATE TABLE IF NOT EXISTS clipboards(
		id INTEGER PRIMARY KEY,
		user_id TEXT,
		sender_public_key BLOB,
		FOREIGN KEY(user_id) REFERENCES users(id)
	) STRICT;
	CREATE TABLE IF NOT EXISTS clipboard_items(
		id INTEGER PRIMARY KEY,
		ciphertext BLOB,
		clipboard_id INTEGER,
		device_id TEXT,
		FOREIGN KEY(clipboard_id) REFERENCES clipboards(id),
		FOREIGN KEY(device_id) REFERENCES devices(id)
	) STRICT;
	`
	_, err := d.internalDB.Exec(setupStm)
	return err
}
