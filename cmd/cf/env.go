/*
 * Copyright 2018 The Service Manager Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cf

import (
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/Peripli/service-manager/server"
	"fmt"
)

func newEnv(delegate server.Environment) *cfEnvironment {
	return &cfEnvironment{delegate: delegate}
}

type cfEnvironment struct {
	cfEnv    *cfenv.App
	delegate server.Environment
}

func (e *cfEnvironment) Load() error {
	var err error
	if err = e.delegate.Load(); err != nil {
		return err
	}
	e.cfEnv, err = cfenv.Current()
	return err
}

func (e *cfEnvironment) Get(key string) interface{} {
	value, exists := cfenv.CurrentEnv()[key]
	if !exists {
		return e.delegate.Get(key)
	}
	return value
}

func (e *cfEnvironment) Unmarshal(value interface{}) error {
	err := e.delegate.Unmarshal(value)
	cfg, ok := value.(*server.Settings)
	if ok {
		storageName := "postgres"
		service, err := e.cfEnv.Services.WithName(storageName)
		if err != nil {
			return fmt.Errorf("No service with name %s is bound to the application", storageName)
		}
		uri := service.Credentials["uri"]
		if ok {
			cfg.Db = &server.DbSettings{
				URI: uri.(string),
			}
		}
		cfg.Server = &server.AppSettings{
			Port: e.cfEnv.Port,
		}
	}
	return err
}
