package models

import (
	"encoding/json"
	"io"
	"time"
)

type PostRequest struct {
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
	TTLSeconds int         `json:"ttl_seconds"`
}

func (v *PostRequest) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(v)
}

func (v *PostRequest) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(v)
}

type GetResponse struct {
	Key           string      `json:"key,omitempty"`
	Value         interface{} `json:"value,omitempty"`
	TimeExpiresAt time.Time   `json:"-"`
	ExpiresAt     int         `json:"expires_at,omitempty"`
}

func (v *GetResponse) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	v.TimeExpiresAt = time.Unix(int64(v.ExpiresAt), 0)
	return d.Decode(v)
}

func (v *GetResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	v.ExpiresAt = int(v.TimeExpiresAt.Unix())
	return e.Encode(v)
}

type GetAllResponse struct {
	Keys   []string      `json:"keys"`
	Values []interface{} `json:"values"`
}

func (v *GetAllResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(v)
}
