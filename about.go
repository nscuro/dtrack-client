package dtrack

import (
	"github.com/google/uuid"
)

type About struct {
	UUID        uuid.UUID      `json:"uuid"`
	SystemUUID  uuid.UUID      `json:"systemUuid"`
	Application string         `json:"application"`
	Version     string         `json:"version"`
	Timestamp   string         `json:"timestamp"`
	Framework   AboutFramework `json:"framework"`
}

type AboutFramework struct {
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	Timestamp string    `json:"timestamp"`
}
