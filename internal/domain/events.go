package domain

type Event struct {
	Name   string
	CaseID string
	FileID string
	Actor  string
	Data   string
}

func NewEvent(name, caseID, fileID, actor, data string) Event {
	return Event{Name: name, CaseID: caseID, FileID: fileID, Actor: actor, Data: data}
}

func EventNeedsAudit(name string) bool {
	switch name {
	case "upload", "edit", "download", "share", "archive", "restore":
		return true
	default:
		return false
	}
}
