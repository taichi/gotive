/* Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package log

import (
	"fmt"
	"github.com/op/go-logging"
)

var module string = "gotive"
var logger = logging.MustGetLogger(module)

func init() {
	DisableLog()
}

func VerboseLog() {
	logging.SetLevel(logging.DEBUG, module)
}

func DisableLog() {
	logging.SetLevel(logging.WARNING, module)
}

func UseLogger(newLogger *logging.Logger) {
	if newLogger == nil {
		panic("logger is nil")
	}
	logger = newLogger
}

func IsDebugEnabled() bool {
	return logging.DEBUG <= logging.GetLevel(module)
}

func Debugf(format string, params ...interface{})    { logger.Debug(format, params...) }
func Infof(format string, params ...interface{})     { logger.Info(format, params...) }
func Warnf(format string, params ...interface{})     { logger.Warning(format, params...) }
func Errorf(format string, params ...interface{})    { logger.Error(format, params...) }
func Criticalf(format string, params ...interface{}) { logger.Critical(format, params...) }
func Panicf(format string, params ...interface{})    { logger.Panicf(format, params...) }
func Fatalf(format string, params ...interface{})    { logger.Fatalf(format, params...) }

func Debug(v ...interface{})    { logger.Debug("%s", fmt.Sprint(v...)) }
func Info(v ...interface{})     { logger.Info("%s", fmt.Sprint(v...)) }
func Warn(v ...interface{})     { logger.Warning("%s", fmt.Sprint(v...)) }
func Error(v ...interface{})    { logger.Error("%s", fmt.Sprint(v...)) }
func Critical(v ...interface{}) { logger.Critical("%s", fmt.Sprint(v...)) }
func Panic(v ...interface{})    { logger.Panic(v...) }
func Fatal(v ...interface{})    { logger.Fatal(v...) }
