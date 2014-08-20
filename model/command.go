package model

type CommandItem struct {
	UID  int
	Date string      // TODO: fix it with time.Time
	Data CommandData // TODO: fix name
}

type CommandData struct {
	Command string // need to
	Desc    string
}
