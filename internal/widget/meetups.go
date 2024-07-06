package widget

import (
	"context"
	"fmt"
	"html/template"
	"slices"
	"time"

	"github.com/glanceapp/glance/internal/assets"
	"github.com/glanceapp/glance/internal/feed"
)

type Meetups struct {
	widgetBase `yaml:",inline"`
	Meetups    feed.Meetups `yaml:"-"`

	Limit         int    `yaml:"limit"`
	City          string `yaml:"city"`
	Radius        uint   `yaml:"radius"`
	TopicCategory int    `yaml:"topic-id"`
	DayOffset     int    `yaml:"day-offest"`
	DaysToFetch   int    `yaml:"days-to-fetch"`
	DisplayType   string `yaml:"display-type"`
}

func (widget *Meetups) Initialize() error {
	widget.withTitle(fmt.Sprintf("Meetups in %s", widget.City)).withCacheDuration(time.Hour)

	if widget.Limit <= 0 {
		widget.Limit = 20
	}

	if widget.Radius <= 0 {
		widget.Radius = 10
	}

	if widget.DaysToFetch <= 0 {
		widget.DaysToFetch = 7
	}
	if widget.DisplayType != "grid" {
		widget.DisplayType = "cards"
	}

	return nil
}

func (widget *Meetups) Update(ctx context.Context) {
	meetups, err := feed.FetchMeetups(widget.City, widget.Radius, widget.TopicCategory, widget.DayOffset, widget.DaysToFetch)
	if len(meetups) > widget.Limit {
		meetups = meetups[:widget.Limit]
	}
	if !widget.canContinueUpdateAfterHandlingErr(err) {
		return
	}

	slices.SortStableFunc(meetups, func(a, b feed.Meetup) int {
		if a.Time.Before(b.Time) {
			return -1
		}
		if a.Time.After(b.Time) {
			return 1
		}
		return 0
	})
	widget.Meetups = meetups
}

func (widget *Meetups) Render() template.HTML {
	if widget.DisplayType == "grid" {
		return widget.render(widget, assets.MeetupsGridTemplate)
	}
	return widget.render(widget, assets.MeetupsTemplate)
}
