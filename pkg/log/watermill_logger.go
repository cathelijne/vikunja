// Vikunja is a to-do list application to facilitate your life.
// Copyright 2018-2021 Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public Licensee as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public Licensee for more details.
//
// You should have received a copy of the GNU Affero General Public Licensee
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package log

import (
	"fmt"
	"strings"
	"time"

	"code.vikunja.io/api/pkg/config"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/op/go-logging"
)

const watermillFmt = `%{color}%{time:` + time.RFC3339Nano + `}: %{level}` + "\t" + `▶ [EVENTS] %{id:03x}%{color:reset} %{message}`

const watermillLogModule = `vikunja_events`

type WatermillLogger struct {
	logger *logging.Logger
}

func NewWatermillLogger() *WatermillLogger {
	lvl := strings.ToUpper(config.LogEventsLevel.GetString())
	level, err := logging.LogLevel(lvl)
	if err != nil {
		Criticalf("Error setting events log level %s: %s", lvl, err.Error())
	}

	watermillLogger := &WatermillLogger{
		logger: logging.MustGetLogger(watermillLogModule),
	}

	cf := config.LogEvents.GetString()
	var backend logging.Backend
	backend = &NoopBackend{}
	if cf != "off" && cf != "false" {
		logBackend := logging.NewLogBackend(GetLogWriter("events"), "", 0)
		backend = logging.NewBackendFormatter(logBackend, logging.MustStringFormatter(watermillFmt+"\n"))
	}

	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(level, watermillLogModule)

	watermillLogger.logger.SetBackend(backendLeveled)

	return watermillLogger
}

func concatFields(fields watermill.LogFields) string {
	full := ""

	for key, val := range fields {
		full += fmt.Sprintf("%s=%v, ", key, val)
	}

	if full != "" {
		full = full[:len(full)-2]
	}

	return full
}

func (w *WatermillLogger) Error(msg string, err error, fields watermill.LogFields) {
	w.logger.Errorf("%s: %s, %s", msg, err, concatFields(fields))
}

func (w *WatermillLogger) Info(msg string, fields watermill.LogFields) {
	w.logger.Infof("%s, %s", msg, concatFields(fields))
}

func (w *WatermillLogger) Debug(msg string, fields watermill.LogFields) {
	w.logger.Debugf("%s, %s", msg, concatFields(fields))
}

func (w *WatermillLogger) Trace(msg string, fields watermill.LogFields) {
	w.logger.Debugf("%s, %s", msg, concatFields(fields))
}

func (w *WatermillLogger) With(_ watermill.LogFields) watermill.LoggerAdapter {
	return w
}
