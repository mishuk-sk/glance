package widget

import (
	"context"
	"html/template"
	"time"

	"github.com/glanceapp/glance/internal/assets"
	"github.com/glanceapp/glance/internal/feed"
)

type TasksWidget struct {
	widgetBase `yaml:",inline"`
	Tasks      feed.Tasks `yaml:"-"`

	Limit   int    `yaml:"limit"`
	BaseURL string `yaml:"base-url"`
	Status  string `yaml:"status"`
}

func (widget *TasksWidget) Initialize() error {
	widget.withTitle("Tasks").withCacheDuration(time.Nanosecond)

	if widget.Limit <= 0 {
		widget.Limit = 20
	}

	return nil
}

func (widget *TasksWidget) Update(ctx context.Context) {
	tasks, err := feed.FetchTasks(widget.BaseURL, widget.Status)
	if len(tasks) > widget.Limit {
		tasks = tasks[:widget.Limit]
	}
	if !widget.canContinueUpdateAfterHandlingErr(err) {
		return
	}

	widget.Tasks = tasks
}

func (widget *TasksWidget) Render() template.HTML {
	return widget.render(widget, assets.TasksTemplate)
}
