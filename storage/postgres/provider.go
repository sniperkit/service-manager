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

package postgres

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Peripli/service-manager/env"
	"github.com/Peripli/service-manager/storage"
	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func init() {
	storage.Register("postgres", &provider{})
}

type provider struct {
	once           sync.Once
	dbStorage      storage.Storage
	provisionError error
}

func (p *provider) Provide() (storage.Storage, error) {
	p.once.Do(func() {
		closeChan := make(chan os.Signal)
		signal.Notify(closeChan, os.Interrupt, syscall.SIGTERM)

		uri, ok := env.Get("storage.uri").(string)
		if !ok {
			logrus.Panicf("Could not open connection for provided uri %s from postgres storage provider", uri)
		}
		db, err := sqlx.Open("postgres", uri)
		if err != nil {
			logrus.Panicf("Could not connect to PostgreSQL storage" + err.Error())
		}

		go func() {
			<-closeChan
			logrus.Debug("Received close for postgres storage")
			close(closeChan)
			if err := db.Close(); err != nil {
				logrus.Panic(err)
			}
		}()

		logrus.Debug("Initialized PostgreSQL storage")
		p.dbStorage, p.provisionError = newStorage(db)
	})
	return p.dbStorage, p.provisionError
}