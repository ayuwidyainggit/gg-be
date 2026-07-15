package model

type Status struct {
	StatusId    string `db:"status_id" json:"status_id"`
	StatusName  string `db:"status_name" json:"status_name"`
	StatusValue int    `db:"status_value" json:"status_value"`
	LangId      string `db:"lang_id" json:"lang_id"`
}

// type StatusUpdate struct {
// StatusID    string        `db:"status_id"`
// StatusName  string        `db:"status_name"`
// StatusValue sql.NullInt64 `db:"status_value"`
// LangID      string        `db:"lang_id"`
// }
