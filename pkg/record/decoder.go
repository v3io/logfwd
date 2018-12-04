package record

import (
	"encoding/json"
	"io"

	"github.com/nuclio/logger"
	"github.com/v3io/go-errors"
)

type Decoder struct {
	log logger.Logger
}

func NewDecoder(log logger.Logger) *Decoder {
	return &Decoder{log}
}

func (d *Decoder) FromString(raw string) (*LogRecord, error) {
	return d.FromByteArray([]byte(raw))
}

func (d *Decoder) FromByteArray(raw []byte) (*LogRecord, error) {
	result := LogRecord{}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, errors.Wrapf(err, "Unable to parse '%s'", raw)
	}
	return &result, nil
}

func (d *Decoder) FromReader(reader io.Reader) (*LogRecord, error) {
	result := LogRecord{}
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&result); err != nil {
		return nil, errors.Wrapf(err, "Unable to from reader")
	}
	return &result, nil
}

type ArrayDecoder struct {
	log logger.Logger
}

func NewArrayDecoder(log logger.Logger) *ArrayDecoder {
	return &ArrayDecoder{log}
}

func (d *ArrayDecoder) FromString(raw string) (LogRecords, error) {
	return d.FromByteArray([]byte(raw))
}

func (d *ArrayDecoder) FromByteArray(raw []byte) (LogRecords, error) {
	result := LogRecords{}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, errors.Wrapf(err, "Unable to parse '%s'", raw)
	}
	return result, nil
}

func (d *ArrayDecoder) FromReader(reader io.Reader) (LogRecords, error) {
	result := LogRecords{}
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&result); err != nil {
		return nil, errors.Wrapf(err, "Unable to from reader")
	}
	return result, nil
}
