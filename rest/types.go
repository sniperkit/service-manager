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

// ErrorResponse struct used to store information about error
type ErrorResponse struct {
	ErrorType   string `json:"error,omitempty"`
	Description string `json:"description"`
	StatusCode  int    `json:"-"`
}

// Error ErrorResponse should implement error
func (errorResponse ErrorResponse) Error() string {
	return errorResponse.Description
}

// Credentials credentials
type Credentials struct {
	Type     string `json:"type"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// Platform platform struct
type Platform struct {
	ID          string       `json:"id"`
	Type        string       `json:"type"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	CreatedAt   string       `json:"created_at,omitempty"`
	UpdatedAt   string       `json:"updated_at,omitempty"`
	Credentials *Credentials `json:"credentials,omitempty"`
}

// Broker broker struct
type Broker struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	URL         string                 `json:"broker_url"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
	Credentials *Credentials           `json:"credentials,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}
