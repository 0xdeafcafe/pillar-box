package messagemonitor

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	"github.com/0xdeafcafe/pillar-box/server/internals/broadcaster"
	"github.com/0xdeafcafe/pillar-box/server/internals/codeextractor"
	"github.com/0xdeafcafe/pillar-box/server/internals/streamtyped"
)

type MessageMonitor struct {
	db          *sql.DB
	log         *zap.Logger
	broadcaster *broadcaster.Broadcaster

	latestKnownRecordTimestamp int
}

type ScannedRow struct {
	GUID           string
	AttributedBody []byte
	Date           int
}

func New(ctx context.Context, log *zap.Logger, broadcaster *broadcaster.Broadcaster) (*MessageMonitor, error) {
	db, err := sql.Open("sqlite3", "/Users/afr/Library/Messages/chat.db")
	if err != nil {
		return nil, err
	}

	return &MessageMonitor{
		db:                         db,
		log:                        log,
		broadcaster:                broadcaster,
		latestKnownRecordTimestamp: 0,
	}, nil
}

func (m *MessageMonitor) ListenAndHandle(ctx context.Context) {
	// TODO(afr): Use FS monitoring to detect new messages instead of polling?

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
			if mfaCode == nil {
				m.latestKnownRecordTimestamp = row.Date
				m.log.Info("no mfa code found in message", zap.String("message", *message))
				continue
			}

			m.broadcaster.BroadcastMFACode(*mfaCode)
			m.latestKnownRecordTimestamp = row.Date
		}

		time.Sleep(2 * time.Second)
	}
}
