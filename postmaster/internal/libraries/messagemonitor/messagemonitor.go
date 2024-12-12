package messagemonitor

import (
	"database/sql"
	"os"
	"path"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	"github.com/0xdeafcafe/pillar-box/server/internal/utilities/codeextractor"
	"github.com/0xdeafcafe/pillar-box/server/internal/utilities/streamtyped"
)

type MessageMonitor struct {
	db  *sql.DB
	log *zap.Logger

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
func New(log *zap.Logger) (*MessageMonitor, error) {
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
		log:                         log,
		latestKnownRecordTimestamp:  0,
		registeredDetectionHandlers: make([]DetectionHandlerFunc, 0),
	}, nil
}

func (m *MessageMonitor) RegisterDetectionHandler(handleMessageDetection DetectionHandlerFunc) {
	m.registeredDetectionHandlers = append(m.registeredDetectionHandlers, handleMessageDetection)
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
			m.log.Error("failed to query database", zap.Error(err))
			time.Sleep(5 * time.Second)

			continue
		}

		scannedRows := make([]*ScannedRow, 0)

		for rows.Next() {
			scannedRow := &ScannedRow{}

			if err := rows.Scan(&scannedRow.GUID, &scannedRow.AttributedBody, &scannedRow.Date); err != nil {
				m.log.Error("failed to scan row", zap.Error(err))
				time.Sleep(5 * time.Second)
				continue
			}

			scannedRows = append(scannedRows, scannedRow)
		}

		for _, row := range scannedRows {
			message, err := streamtyped.ExtractMessageFromStreamTypedBuffer(row.AttributedBody)
			if err != nil {
				m.log.Error("failed to extract message from streamtyped buffer", zap.Error(err))
				continue
			}

			m.log.Info("discovered mfa code", zap.String("message", *message))

			mfaCode, err := codeextractor.ExtractMFACodeFromMessage(*message)
			if err != nil {
				m.latestKnownRecordTimestamp = row.Date
				m.log.Warn("failed to extract mfa code from message", zap.Error(err), zap.String("message", *message))

				continue
			}
			if mfaCode == "" {
				m.latestKnownRecordTimestamp = row.Date
				m.log.Info("no mfa code found in message", zap.String("message", *message))

				continue
			}

			m.latestKnownRecordTimestamp = row.Date

			for _, handler := range m.registeredDetectionHandlers {
				handler(mfaCode)
			}
		}

		time.Sleep(1 * time.Second)
	}
}
