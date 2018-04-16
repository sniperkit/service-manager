/*
 * Copyright 2018 The Service Manager Authors
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package rest

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// HandleError sends a JSON containing the error to the response writer
func HandleError(err error, writer http.ResponseWriter) {
	if err != nil {
		var respError *ErrorResponse
		switch t := err.(type) {
		case *ErrorResponse:
			logrus.Debug(err)
			respError = t
		case ErrorResponse:
			logrus.Debug(err)
			respError = &t
		default:
			logrus.Error(err)
			respError = &ErrorResponse{
				ErrorType:   "InternalError",
				Description: "Internal server error",
				StatusCode:  http.StatusInternalServerError,
			}
		}

		sendErr := SendJSON(writer, respError.StatusCode, respError)
		if sendErr != nil {
			logrus.Errorf("Could not write error to response: %v", sendErr)
		}
	}

}

// CreateErrorResponse create ErrorResponse object
func CreateErrorResponse(err error, statusCode int, errorType string) *ErrorResponse {
	return &ErrorResponse{
		ErrorType:   errorType,
		Description: err.Error(),
		StatusCode:  statusCode}
}
