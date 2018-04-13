/*
 *    Copyright 2018 The Service Manager Authors
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

package platform

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Peripli/service-manager/rest"
	"github.com/Peripli/service-manager/storage"
	"github.com/Peripli/service-manager/util"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// Controller platform controller
type Controller struct{}

func getPlatformID(req *http.Request) string {
	vars := mux.Vars(req)
	return vars["platform_id"]
}

func getPlatformFromRequest(req *http.Request) (*rest.Platform, error) {
	var platform rest.Platform
	return &platform, rest.ReadJSONBody(req, &platform)
}

func mergePlatforms(source *rest.Platform, target *rest.Platform) {
	if source.Name != "" {
		target.Name = source.Name
	}
	if source.Description != "" {
		target.Description = source.Description
	}
	if source.Type != "" {
		target.Type = source.Type
	}
	target.UpdatedAt = time.Now().UTC()
}

func checkPlatformMandatoryProperties(platform *rest.Platform) error {
	if platform.Type == "" {
		return errors.New("Missing platform type")
	}
	if platform.Name == "" {
		return errors.New("Missing platform name")
	}
	return nil
}

func errorMissingPlatform(platformID string) error {
	return fmt.Errorf("Could not find platform with id %s", platformID)
}

// addPlatform handler for POST /v1/platforms
func (platformCtrl *Controller) addPlatform(rw http.ResponseWriter, req *http.Request) error {
	logrus.Debugf("POST to %s", req.RequestURI)

	platform, errDecode := getPlatformFromRequest(req)
	if errDecode != nil {
		return errDecode
	}
	if errMandatoryProperties := checkPlatformMandatoryProperties(platform); errMandatoryProperties != nil {
		return rest.CreateErrorResponse(errMandatoryProperties, http.StatusBadRequest, "BadRequest")
	}
	username, password := util.GenerateCredentials()
	if platform.ID == "" {
		uuid, err := uuid.NewV4()
		if err != nil {
			logrus.Error("Could not generate GUID")
			return err
		}
		platform.ID = uuid.String()
	}
	currentTime := time.Now().UTC()
	platform.CreatedAt = currentTime
	platform.UpdatedAt = currentTime

	platform.Credentials = &rest.Credentials{
		Basic: &rest.Basic{
			Username: username,
			Password: password,
		},
	}
	platformStorage := storage.Get().Platform()
	errSave := platformStorage.Create(platform)
	if errSave == storage.ErrUniqueViolation {
		return rest.CreateErrorResponse(errSave, http.StatusConflict, "Conflict")
	} else if errSave != nil {
		return errSave
	}

	return rest.SendJSON(rw, http.StatusCreated, platform)
}

// getPlatform handler for GET /v1/platforms/:platform_id
func (platformCtrl *Controller) getPlatform(rw http.ResponseWriter, req *http.Request) error {
	logrus.Debugf("GET to %s", req.RequestURI)
	platformID := getPlatformID(req)
	platformStorage := storage.Get().Platform()
	platform, err := platformStorage.Get(platformID)
	if err == storage.ErrNotFound {
		return rest.CreateErrorResponse(errorMissingPlatform(platformID), http.StatusNotFound, "NotFound")
	} else if err != nil {
		return err
	}
	return rest.SendJSON(rw, http.StatusOK, platform)
}

// getAllPlatforms handler for GET /v1/platforms
func (platformCtrl *Controller) getAllPlatforms(rw http.ResponseWriter, req *http.Request) error {
	logrus.Debugf("GET to %s", req.RequestURI)
	platformStorage := storage.Get().Platform()
	platforms, err := platformStorage.GetAll()
	if err != nil {
		return err
	}
	platformsResponse := map[string][]rest.Platform{"platforms": platforms}

	return rest.SendJSON(rw, http.StatusOK, &platformsResponse)
}

// deletePlatform handler for DELETE /v1/platforms/:platform_id
func (platformCtrl *Controller) deletePlatform(rw http.ResponseWriter, req *http.Request) error {
	logrus.Debugf("DELETE to %s", req.RequestURI)
	platformID := getPlatformID(req)

	platformStorage := storage.Get().Platform()
	errDelete := platformStorage.Delete(platformID)
	if errDelete == storage.ErrNotFound {
		return rest.CreateErrorResponse(errorMissingPlatform(platformID), http.StatusNotFound, "NotFound")
	} else if errDelete != nil {
		return errDelete
	}
	// map[string]string{} will result in empty JSON
	return rest.SendJSON(rw, http.StatusOK, map[string]string{})
}

// updatePlatform handler for PATCH /v1/platforms/:platform_id
func (platformCtrl *Controller) updatePlatform(rw http.ResponseWriter, req *http.Request) error {
	logrus.Debugf("PATCH to %s", req.RequestURI)
	platformID := getPlatformID(req)
	newPlatform, errDecode := getPlatformFromRequest(req)
	if errDecode != nil {
		return errDecode
	}
	newPlatform.ID = platformID
	newPlatform.UpdatedAt = time.Now().UTC()
	platformStorage := storage.Get().Platform()
	err := platformStorage.Update(newPlatform)
	if err != nil {
		if err == storage.ErrUniqueViolation {
			return rest.CreateErrorResponse(err, http.StatusConflict, "Conflict")
		} else if err == storage.ErrNotFound {
			return rest.CreateErrorResponse(errorMissingPlatform(platformID), http.StatusNotFound, "NotFound")
		}
		return err
	}
	platform, err := platformStorage.Get(platformID)
	if err != nil {
		return err
	}
	return rest.SendJSON(rw, http.StatusOK, platform)
}
