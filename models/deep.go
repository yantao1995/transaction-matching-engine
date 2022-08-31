package models

type Deep struct {
	Pair          string      `json:"pair"`
	TimeUnixMilli int64       `json:"time_unix_milli"`
	Bids          [][3]string `json:"bids"`
	Asks          [][3]string `json:"asks"`
}
