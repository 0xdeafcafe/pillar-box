package streamtyped

import (
	"bytes"
	"errors"
	"unicode/utf8"

	"github.com/0xdeafcafe/pillar-box/server/internal/utilities/ptr"
)

var (
	// encodedMessageHeader is the binary trailer indicator for a message. Instead of
	// spending a day reversing this deprecated and undocumented format, this works fine.
	encodedMessageTrailer = []byte{0x86, 0x84, 0x02, 0x69, 0x49, 0x01}

	// streamTypedMagic is the magic bytes that identify a streamtyped buffer.
	streamTypedMagic = []byte{0x04, 0x0b}

	// streamTypedIdentifier is the identifier for a streamtyped buffer.
	streamTypedIdentifier = []byte("streamtyped")

	ErrInvalidStreamTypedBuffer    = errors.New("invalid streamtyped buffer")
	ErrStreamTypedBufferTooShort   = errors.New("streamtyped buffer too short")
	ErrStreamTypedMagicMismatch    = errors.New("streamtyped magic mismatch")
	ErrStreamTypedNoMessageTrailer = errors.New("streamtyped message has no trailer")
	ErrStreamTypedMessageInvalid   = errors.New("streamtyped message is invalid")
)

const (
	// encodedMessageOffset is the offset of the encoded buffer that the actual message
	// data (in utf-8 encoding) is located.
	encodedMessageOffset = 0x7a

	// streamTypedIdentifierOffset is the offset of the streamtyped identifier in the buffer.
	streamTypedIdentifierOffset = 0x02

	// streamTypedIdentifierLength is the length of the streamtyped identifier.
	streamTypedIdentifierLength = 0x0a
)

// ExtractMessageFromStreamTypedBuffer extracts the message from a buffer encoded in the
// streamtyped format.
func ExtractMessageFromStreamTypedBuffer(buffer []byte) (*string, error) {
	if len(buffer) < encodedMessageOffset {
		return nil, errors.Join(ErrInvalidStreamTypedBuffer, ErrStreamTypedBufferTooShort)
	}

	if !bytes.HasPrefix(buffer, streamTypedMagic) {
		return nil, errors.Join(ErrInvalidStreamTypedBuffer, ErrStreamTypedMagicMismatch)
	}

	bufferIdentifier := buffer[streamTypedIdentifierOffset : streamTypedIdentifierOffset+streamTypedIdentifierLength]
	if bytes.Compare(bufferIdentifier, streamTypedIdentifier) == 0 {
		return nil, errors.Join(ErrInvalidStreamTypedBuffer, ErrStreamTypedMagicMismatch)
	}

	removeMessagePrefix := buffer[encodedMessageOffset:]
	messageEnd := bytes.Index(removeMessagePrefix, encodedMessageTrailer)
	if messageEnd == -1 {
		return nil, errors.Join(ErrInvalidStreamTypedBuffer, ErrStreamTypedNoMessageTrailer)
	}

	if !utf8.Valid(removeMessagePrefix[:messageEnd]) {
		return nil, errors.Join(ErrInvalidStreamTypedBuffer, ErrStreamTypedMessageInvalid)
	}

	return ptr.Ptr(string(removeMessagePrefix[:messageEnd])), nil
}
