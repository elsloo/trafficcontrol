package tc

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-rfc"
)

// Alert represents an informational message, typically returned through the Traffic Ops API.
type Alert struct {
	// Text is the actual message being conveyed.
	Text string `json:"text"`
	// Level describes what kind of message is being relayed. In practice, it should be the string
	// representation of one of ErrorLevel, WarningLevel, InfoLevel or SuccessLevel.
	Level string `json:"level"`
}

// Alerts is merely a collection of arbitrary "Alert"s for ease of use in other structures, most
// notably those used in Traffic Ops API responses.
type Alerts struct {
	Alerts []Alert `json:"alerts"`
}

// CreateErrorAlerts creates and returns an Alerts structure filled with ErrorLevel-level "Alert"s
// using the errors to provide text.
func CreateErrorAlerts(errs ...error) Alerts {
	alerts := []Alert{}
	for _, err := range errs {
		if err != nil {
			alerts = append(alerts, Alert{err.Error(), ErrorLevel.String()})
		}
	}
	return Alerts{alerts}
}

// CreateAlerts creates and returns an Alerts structure filled with "Alert"s that are all of the
// provided level, each having one of messages as text in turn.
func CreateAlerts(level AlertLevel, messages ...string) Alerts {
	alerts := []Alert{}
	for _, message := range messages {
		alerts = append(alerts, Alert{message, level.String()})
	}
	return Alerts{alerts}
}

// GetHandleErrorsFunc is used to provide an error-handling function. The error handler provides a
// response to an HTTP request made to the Traffic Ops API and accepts a response code and a set of
// errors to display as alerts.
//
// Deprecated: Traffic Ops API handlers should use
// github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api.HandleErr instead.
func GetHandleErrorsFunc(w http.ResponseWriter, r *http.Request) func(status int, errs ...error) {
	return func(status int, errs ...error) {
		log.Errorf("%v %v\n", r.RemoteAddr, errs)
		errBytes, jsonErr := json.Marshal(CreateErrorAlerts(errs...))
		if jsonErr != nil {
			log.Errorf("failed to marshal error: %s\n", jsonErr)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, http.StatusText(http.StatusInternalServerError))
			return
		}
		w.Header().Set(rfc.ContentType, rfc.ApplicationJSON)

		ctx := r.Context()
		ctx = context.WithValue(ctx, StatusKey, status)
		*r = *r.WithContext(ctx)

		fmt.Fprintf(w, "%s", errBytes)
	}
}

// ToStrings converts Alerts to a slice of strings that are their messages. Note that this return
// value doesn't contain their Levels anywhere.
func (alerts *Alerts) ToStrings() []string {
	alertStrs := []string{}
	for _, alrt := range alerts.Alerts {
		at := alrt.Text
		alertStrs = append(alertStrs, at)
	}
	return alertStrs
}

// AddNewAlert constructs a new Alert with the given Level and Text and appends it to the Alerts
// structure.
func (self *Alerts) AddNewAlert(level AlertLevel, text string) {
	self.AddAlert(Alert{Level: level.String(), Text: text})
}

// AddAlert appends an alert to the Alerts structure.
func (self *Alerts) AddAlert(alert Alert) {
	self.Alerts = append(self.Alerts, alert)
}

// AddAlerts appends all of the "Alert"s in the given Alerts structure to this Alerts structure.
func (self *Alerts) AddAlerts(alerts Alerts) {
	newAlerts := make([]Alert, len(self.Alerts), len(self.Alerts)+len(alerts.Alerts))
	copy(newAlerts, self.Alerts)
	newAlerts = append(newAlerts, alerts.Alerts...)
	self.Alerts = newAlerts
}

var StatusKey = "status"
