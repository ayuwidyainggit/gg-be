package entity

type NotifyWaReq struct {
	To   string `json:"to" validate:"required"`
	Body string `json:"body" validate:"required"`
}
type NotifyCicdWaReq struct {
	To            string `json:"to" validate:"required"`
	Env           string `json:"env" validate:"required"`
	ProjectName   string `json:"project_name" validate:"required"`
	Branch        string `json:"branch" validate:"required"`
	CommitID      string `json:"commit_id" validate:"required"`
	CommitAuthor  string `json:"commit_author" validate:"required"`
	CommitMessage string `json:"commit_message" validate:"required"`
}

type NotifyWaRes struct {
	Sent    bool       `json:"sent"`
	Message MessageObj `json:"message"`
}

type MessageObj struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	ChatID    string `json:"chat_id"`
	ChatName  string `json:"chat_name"`
	From      string `json:"from"`
	FromMe    bool   `json:"from_me"`
	FromName  string `json:"from_name"`
	Source    string `json:"source"`
	Timestamp int    `json:"timestamp"`
	Status    string `json:"status"`
}

type NotifyWaErrRes struct {
	ErrRes ErrObj `json:"error"`
}

type ErrObj struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
