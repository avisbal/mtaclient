package mtaclient

import (
	"io/ioutil"
	"net/http"

	"github.com/avisbal/mtaclient/constants"
	pb "github.com/avisbal/mtaclient/generated/transit_realtime"
	"github.com/golang/protobuf/proto"
)

type MTAClient struct {
	BaseURL string
	ApiKey  string
}

func (m MTAClient) GetArrivalTimes(feedID constants.FeedID, stopID constants.StopID) ([]int64, error) {
	feed, err := m.getFeed(feedID)
	if err != nil {
		return nil, err
	}

	arrivalTimes := make([]int64, 0)
	for _, entity := range feed.GetEntity() {
		tripUpdate := entity.GetTripUpdate()
		if tripUpdate != nil {
			for _, stopTimeUpdate := range tripUpdate.GetStopTimeUpdate() {
				arrival := stopTimeUpdate.GetArrival()
				if stopID.String() == stopTimeUpdate.GetStopId() && arrival != nil {
					arrivalTime := arrival.GetTime()
					arrivalTimes = append(arrivalTimes, arrivalTime)
				}
			}
		}
	}

	return arrivalTimes, nil
}

func (m MTAClient) getFeed(feedID constants.FeedID) (*pb.FeedMessage, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", m.BaseURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("key", m.ApiKey)
	q.Add("feed_id", feedID.String())
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return unmarshalFeedMessage(body)
}

func unmarshalFeedMessage(body []byte) (*pb.FeedMessage, error) {
	message := &pb.FeedMessage{}
	err := proto.Unmarshal(body, message)
	if err != nil {
		return nil, err
	}
	return message, nil
}
