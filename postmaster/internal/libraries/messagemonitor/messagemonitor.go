package messagemonitor

import (
	"database/sql"
	"log"
	"os"
	"path"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/exp/rand"

	"github.com/0xdeafcafe/pillar-box/server/internal/utilities/codeextractor"
	"github.com/0xdeafcafe/pillar-box/server/internal/utilities/streamtyped"
)

type MessageMonitor struct {
	db *sql.DB

	registeredDetectionHandlers []DetectionHandlerFunc

	latestKnownRecordTimestamp int
}

type DetectionHandlerFunc func(mfaCode string)

type ScannedRow struct {
	GUID           string
	AttributedBody []byte
	Date           int
}

// New creates a new MessageMonitor instance. The MessageMonitor is responsible for
// monitoring the iMessage database for new messages and extracting MFA codes from them.
// When a new MFA code is detected, the MessageMonitor will call the provided
// HandleMessageDetectionFunc with the detected MFA code.
func New() (*MessageMonitor, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dbPath := path.Join(dirname, "Library/Messages/chat.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return &MessageMonitor{
		db:                          db,
		latestKnownRecordTimestamp:  0,
		registeredDetectionHandlers: make([]DetectionHandlerFunc, 0),
	}, nil
}

func (m *MessageMonitor) RegisterDetectionHandler(handleMessageDetection DetectionHandlerFunc) {
	m.registeredDetectionHandlers = append(m.registeredDetectionHandlers, handleMessageDetection)
}

func (m *MessageMonitor) SendMockMessage() {
	m.dispatchMFACode(generateMockMFACode())
}

func (m *MessageMonitor) ListenAndHandle() {
	for {
		var rows *sql.Rows
		var err error

		if m.latestKnownRecordTimestamp != 0 {
			rows, err = m.db.Query("SELECT guid, attributedBody, date FROM message WHERE service = 'SMS' AND date > ? ORDER BY date ASC;", m.latestKnownRecordTimestamp)
		} else {
			rows, err = m.db.Query("SELECT guid, attributedBody, date FROM message WHERE service = 'SMS' ORDER BY date DESC LIMIT 1;")
		}

		if err != nil {
			log.Printf("failed to query database: %v", err)
			time.Sleep(5 * time.Second)

			continue
		}

		scannedRows := make([]*ScannedRow, 0)

		for rows.Next() {
			scannedRow := &ScannedRow{}

			if err := rows.Scan(&scannedRow.GUID, &scannedRow.AttributedBody, &scannedRow.Date); err != nil {
				log.Printf("failed to scan row: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			scannedRows = append(scannedRows, scannedRow)
		}

		for _, row := range scannedRows {
			message, err := streamtyped.ExtractMessageFromStreamTypedBuffer(row.AttributedBody)
			if err != nil {
				log.Printf("failed to extract message from streamtyped buffer: %v", err)
				continue
			}

			log.Printf("discovered mfa code: %s", *message)

			mfaCode, err := codeextractor.ExtractMFACodeFromMessage(*message)
			if err != nil {
				m.latestKnownRecordTimestamp = row.Date
				log.Printf("failed to extract mfa code from message: %v message:%s", err, *message)

				continue
			}
			if mfaCode == "" {
				m.latestKnownRecordTimestamp = row.Date
				log.Printf("no mfa code found in message: %s", *message)

				continue
			}

			m.latestKnownRecordTimestamp = row.Date
			m.dispatchMFACode(mfaCode)
		}

		time.Sleep(1 * time.Second)
	}
}

func (m *MessageMonitor) dispatchMFACode(mfaCode string) {
	for _, handler := range m.registeredDetectionHandlers {
		handler(mfaCode)
	}
}

func generateMockMFACode() string {
	const charset = "0123456789"

	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
