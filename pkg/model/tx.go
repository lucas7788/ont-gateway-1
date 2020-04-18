package model

import "time"

// TxState for tx state
type TxState int

const (
	// TxStateToPoll for topoll state
	TxStateToPoll TxState = iota
	// TxStateToNotify for tonotify state
	TxStateToNotify
	// TxStateDone for done state
	TxStateDone
)

// TxPollResult for poll result
type TxPollResult int

const (
	// TxPollResultUnknown when result unknown
	TxPollResultUnknown TxPollResult = iota
	// TxPollResultExists when result exists
	TxPollResultExists
	// TxPollResultExpired when expired
	TxPollResultExpired
)

// Tx for txn
type Tx struct {
	Hash           string       `bson:"hash" json:"hash"`
	App            int          `bson:"app" json:"app"`
	State          TxState      `bson:"state" json:"state"`
	Result         TxPollResult `bson:"result" json:"result"`
	PollErrCount   int          `bson:"poll_err_count" json:"poll_err_count"`
	PollErrMsg     string       `bson:"poll_err_msg" json:"poll_err_msg"`
	NotifyErrCount int          `bson:"notify_err_count" json:"notify_err_count"`
	NotifyErrMsg   string       `bson:"notify_err_msg" json:"notify_err_msg"`
	ExpireAt       time.Time    `bson:"expire_at" json:"expire_at"`
	UpdatedAt      time.Time    `bson:"updated_at" json:"updated_at"`
}

// IsExpired tells whether Tx is expired
func (tx *Tx) IsExpired() bool {
	return tx.ExpireAt.Before(time.Now())
}
