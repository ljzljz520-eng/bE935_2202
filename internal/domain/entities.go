package domain

import "time"

type Case struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Client    string    `json:"client"`
	Status    string    `json:"status"`
	Archived  bool      `json:"archived"`
	CreatedAt time.Time `json:"created_at"`
}

type CaseFile struct {
	ID        string    `json:"id"`
	CaseID    string    `json:"case_id"`
	Name      string    `json:"name"`
	Kind      string    `json:"kind"`
	Content   []byte    `json:"content"`
	Version   int       `json:"version"`
	Archived  bool      `json:"archived"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Permission struct {
	CaseID      string `json:"case_id"`
	UserID      string `json:"user_id"`
	Role        string `json:"role"`
	CanDownload bool   `json:"can_download"`
	CanShare    bool   `json:"can_share"`
	CanEdit     bool   `json:"can_edit"`
}

type FileVersion struct {
	ID        string    `json:"id"`
	FileID    string    `json:"file_id"`
	Number    int       `json:"number"`
	Content   []byte    `json:"content"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type AuditEntry struct {
	ID     string    `json:"id"`
	CaseID string    `json:"case_id"`
	FileID string    `json:"file_id"`
	Actor  string    `json:"actor"`
	Action string    `json:"action"`
	Detail string    `json:"detail"`
	At     time.Time `json:"at"`
}

type ShareLink struct {
	ID        string    `json:"id"`
	FileID    string    `json:"file_id"`
	CreatedBy string    `json:"created_by"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
}

type SearchResult struct {
	File      CaseFile
	CaseTitle string
	Score     int
}

type ArchiveRequest struct {
	CaseID string
	Actor  string
	Reason string
}

func (c Case) IsOpen() bool         { return c.Status == "open" && !c.Archived }
func (c CaseFile) IsEditable() bool { return !c.Archived && c.Version > 0 }
