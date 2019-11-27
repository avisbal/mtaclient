package constants

type FeedID string

const (
	L FeedID = "2"
)

func (feedID FeedID) String() string {
	return string(feedID)
}

type StopID string

const (
	BedfordAvNorthbound StopID = "L08N"
)

func (stopID StopID) String() string {
	return string(stopID)
}
