package models

type RecordWordRequest struct {
	ID   string `json:"id"`
	Word string `json:"word"`
	Type string `json:"type"`
}
