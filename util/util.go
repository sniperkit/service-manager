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

package util

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// GenerateCredentials return user and password
func GenerateCredentials() (string, string, string, error) {
	plainPassword := make([]byte, 128)
	user := make([]byte, 128)

	_, err := rand.Read(user)
	if err != nil {
		return "", "", "", err
	}

	_, err = rand.Read(plainPassword)
	if err != nil {
		return "", "", "", err
	}

	encodedUser := base64.StdEncoding.EncodeToString(user)
	encodedPlainPassword := base64.StdEncoding.EncodeToString(plainPassword)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(encodedPlainPassword), 14)
	if err != nil {
		return "", "", "", err
	}

	return encodedUser, encodedPlainPassword, string(hashedPassword), nil
}

// ToRFCFormat return the time.Time object as string in RFC3339 format
func ToRFCFormat(timestamp time.Time) string {
	return timestamp.UTC().Format(time.RFC3339)
}
