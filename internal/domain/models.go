package domain

import "time"

type User struct {
	ID     string
	Email  string
	Role   string
	Active bool
}
type Formulation struct {
	ID        string
	Name      string
	Version   int
	Chemistry string
	Approved  bool
}
type CellLot struct {
	ID                   string
	Code                 string
	FormulationID        string
	State                LotState
	Capacity             float64
	Version              int64
	CreatedAt, UpdatedAt time.Time
}
type Module struct {
	ID      string
	Serial  string
	LotID   string
	State   string
	Version int64
}
type QualificationRun struct {
	ID                    string
	LotID                 string
	ModuleID              string
	State                 RunState
	ScheduledAt           time.Time
	StartedAt, FinishedAt *time.Time
	Attempts              int
	Version               int64
	CreatedBy             string
}
type Measurement struct {
	ID, RunID, Kind, Unit, RecordedBy string
	Value                             float64
	RecordedAt                        time.Time
}
type QualityHold struct {
	ID, LotID, Reason, OpenedBy string
	State                       HoldState
	ClearedBy                   *string
	OpenedAt                    time.Time
	ClearedAt                   *time.Time
}
type AuditEvent struct {
	ID, ActorID, EntityType, EntityID, Action, RequestID string
	Payload                                              []byte
	CreatedAt                                            time.Time
}
