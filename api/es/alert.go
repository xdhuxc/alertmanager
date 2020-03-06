package es

import (
	"fmt"
	"strconv"
	"time"

	"github.com/prometheus/alertmanager/types"
	"github.com/prometheus/common/model"
)

type LabelValue interface{}

type LabelSet map[model.LabelName]LabelValue

const (
	labelNameValue model.LabelName = "value"
)

// for ES number type index
type Alert struct {
	// for ES number type index
	Labels LabelSet `json:"labels"`

	// Extra key/value information which does not define alert identity.
	Annotations model.LabelSet `json:"annotations"`

	// The known time range for this alert. Both ends are optional.
	StartsAt     time.Time `json:"startsAt,omitempty"`
	EndsAt       time.Time `json:"endsAt,omitempty"`
	GeneratorURL string    `json:"generatorURL"`

	// The authoritative timestamp.
	UpdatedAt time.Time
	Timeout   bool
}

func Convert(alert *types.Alert) *Alert {
	esAlert := new(Alert)

	esAlert.EndsAt = alert.EndsAt
	esAlert.StartsAt = alert.StartsAt
	esAlert.Annotations = alert.Annotations
	esAlert.GeneratorURL = alert.GeneratorURL
	esAlert.Timeout = alert.Timeout
	esAlert.UpdatedAt = alert.UpdatedAt
	esls := make(LabelSet)
	for k, v := range alert.Labels {
		if k == labelNameValue {
			result, err := strconv.ParseFloat(fmt.Sprintf("%s", v), 32)
			if err != nil {
				esls[k] = v
			} else {
				esls[k] = result
			}
		} else {
			esls[k] = v
		}
	}
	esAlert.Labels = esls

	return esAlert
}
