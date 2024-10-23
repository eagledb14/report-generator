package alerts

import (
	"database/sql"
	"fmt"
	// "fmt"
	"time"

	_ "modernc.org/sqlite"
)

type EventCache struct {
	db *sql.DB
}

func NewEventCache() *EventCache {

	db, err := sql.Open("sqlite", "./resources/event_cache.db")
	if err != nil {
		panic("Missing Resoruces")
	}

	newEvent := &EventCache{
		db: db,
	}

	newEvent.ensureTable()
	return newEvent
}

func (e *EventCache) ensureTable() {
	e.db.Exec(`CREATE TABLE IF NOT EXISTS events(
key INTEGER PRIMARY KEY,
ip TEXT,
port INTEGER,
trigger TEXT,
timestamp TEXT
)`)
}

func (e *EventCache) HasEventBeenSeen(event *Event) bool {
	tx, err := e.db.Begin()
	defer tx.Commit()

	rows, err := tx.Query(`SELECT key, timestamp FROM events WHERE ip = ? AND port = ? AND trigger = ?`, event.Ip, event.TriggerPort, event.Trigger)
	if err != nil {
		fmt.Println("querying", err.Error())
		return false
	}

	for rows.Next() {
		var timeString string
		var key int

		if err := rows.Scan(&key, &timeString); err != nil {
			fmt.Println("scanning", err.Error())
			continue
		}

		timestamp, err := time.Parse("02-01-2006", timeString)
		if err != nil {
			fmt.Println(err)
		}

		timeDifference := event.Timestamp.Sub(timestamp).Abs().Hours() / 24

		// reshow an alert every 2 weeks
		if timeDifference <= 14 {
			return true
		} else {
			// If I have seen the event, but it is over 2 weeks old, it should probably be reviewed
			e.deleteByKey(key, tx)
		}
	}

	return false
}

func (e *EventCache) InsertEvent(event *Event) {
	if e.HasEventBeenSeen(event) {
		return
	}

	tx, _ := e.db.Begin()
	defer tx.Commit()

	_, err := tx.Exec(`INSERT INTO events(
ip,
port,
trigger,
timestamp
) VALUES (?,?,?,?)`, event.Ip, event.TriggerPort, event.Trigger, event.Timestamp.Format("02-01-2006"))
	if err != nil {
		fmt.Println("insert", err.Error())
	}
}

func (e *EventCache) deleteByKey(key int, tx *sql.Tx) {
	tx.Exec(`DELETE FROM events WHERE key = ?`, key)
}
