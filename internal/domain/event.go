package domain

import "time"

type Event struct {
	ID        int64             `json:"id"`
	Name      string            `json:"name"`
	Payload   map[string]string `json:"payload"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

type EventList struct {
	ResultsLength int      `json:"resultsLength"`
	Results       []*Event `json:"results"`
}
