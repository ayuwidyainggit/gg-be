package model

type ActivitiesSummary struct {
	LastUpdate    string  `json:"last_update" db:"last_update"`
	Plan          *int    `json:"plan" db:"plan"`
	Visit         *int    `json:"visit" db:"visit"`
	ExtraCall     *int    `json:"extra_call" db:"extra_call"`
	Skip          *int    `json:"skip" db:"skip"`
	EffectiveCall *int    `json:"effective_call" db:"effective_call"`
	StartTime     *string `json:"start_time" db:"start_time"`
	EndTime       *string `json:"end_time" db:"end_time"`
	// DriveTime     string `json:"drive_time" db:"drive_time"`
	// EstTime       string `json:"est_time" db:"est_time"`
}
