package feed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type meetupResponse struct {
	Data struct {
		RecommendedEvents struct {
			Edges []struct {
				Node struct {
					ID        string `json:"id"`
					Title     string `json:"title"`
					DateTime  string `json:"dateTime"`
					IsOnline  bool   `json:"isOnline"`
					EventUrl  string `json:"eventUrl"`
					EventType string `json:"eventType"`
					Rsvps     struct {
						TotalCount int `json:"totalCount"`
					} `json:"rsvps"`
					Venue struct {
						City string `json:"city"`
						Name string `json:"name"`
					} `json:"venue"`
					FeaturedEventPhoto struct {
						HighResUrl string `json:"highResUrl"`
					} `json:"featuredEventPhoto"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"recommendedEvents"`
	} `json:"data"`
}

func FetchMeetups(city string, radius uint, topicCategoryID int, dayOffset int, daysToFetch int) (Meetups, error) {
	endpoint := "https://api.meetup.com/gql2"
	startDate := time.Now().AddDate(0, 0, dayOffset)
	endDate := startDate.AddDate(0, 0, daysToFetch)
	query := fmt.Sprintf(`
		query {
			recommendedEvents(filter: {
				city: "%s",
				radius: %d,
				topicCategoryId: %d,
				startDateRange: "%s",
				endDateRange: "%s"
			}) {
				edges {
					node {
						id,
						title,
						dateTime,
						isOnline,
						eventUrl,
						eventType,
						venue {
							city,
							name
						},
						rsvps {
							totalCount
						},
						featuredEventPhoto {
							highResUrl
						}
					}
				}
			}
		}`, city, radius, topicCategoryID, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339))

	requestBody, err := json.Marshal(map[string]string{
		"query": query,
	})
	fmt.Println(query)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response meetupResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	var meetups Meetups

	for _, edge := range response.Data.RecommendedEvents.Edges {
		node := edge.Node
		eventTime, err := time.Parse(time.RFC3339, node.DateTime)

		if err != nil {
			return nil, err
		}

		location := "Online"
		if !node.IsOnline {
			location = node.Venue.City
		}

		meetups = append(meetups, Meetup{
			Title:        node.Title,
			ThumbnailURL: node.FeaturedEventPhoto.HighResUrl,
			Time:         eventTime,
			URL:          node.EventUrl,
			IsOnline:     node.IsOnline,
			Attendees:    node.Rsvps.TotalCount,
			Location:     location,
		})
	}

	return meetups, nil
}
