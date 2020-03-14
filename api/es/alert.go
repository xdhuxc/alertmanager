package es

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/prometheus/alertmanager/types"
	"github.com/prometheus/common/model"
)

type LabelValue interface{}

type LabelSet map[model.LabelName]LabelValue

const (
	labelNameValue    model.LabelName = "value"
	LabelNameSeverity model.LabelName = "severity"
	LabelNameGroup    model.LabelName = "group"
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

func (a *Alert) ValidateLabels(labels model.LabelSet) error {
	_, ok := labels[labelNameValue]
	if !ok {
		return errors.New("the alert must have a label named value")
	}

	_, ok = labels[LabelNameSeverity]
	if !ok {
		return errors.New("the alert must have a label named severity")
	}

	_, ok = labels[LabelNameGroup]
	if !ok {
		return errors.New("the alert must have a label named group")
	}

	return nil
}

func Convert(alert *types.Alert) (*Alert, error) {
	esAlert := new(Alert)
	err := esAlert.ValidateLabels(alert.Labels)
	if err != nil {
		return nil, err
	}

	esAlert.EndsAt = alert.EndsAt
	esAlert.StartsAt = alert.StartsAt
	esAlert.Annotations = alert.Annotations
	esAlert.GeneratorURL = alert.GeneratorURL
	esAlert.Timeout = alert.Timeout
	esAlert.UpdatedAt = alert.UpdatedAt
	esls := make(LabelSet)
	// It is not suitable to convert all string values to numeric here,
	// there may be different alert using the same key, but the value types are different.
	// When creating an index in ES, an error will occur due to the type
	// make sure that the labels of an alert have the following labels: alertname（inner label of alert），severity, group, value and the value must be a number
	for k, v := range alert.Labels {
		if k == labelNameValue {
			value, err := strconv.ParseFloat(fmt.Sprintf("%s", v), 32)
			if err != nil {
				return nil, err
			} else {
				esls[k] = value
			}
		} else {
			esls[k] = v
		}
	}
	esAlert.Labels = esls

	return esAlert, nil
}
