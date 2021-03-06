/*
 * Copyright (C) 2020  SuperGreenLab <towelie@supergreenlab.com>
 * Author: Constantin Clauzel <constantin.clauzel@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package db

import (
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/postgresql"
)

var (
	Sess       sqlbuilder.Database
	pgPassword = pflag.String("pgpassword", "password", "PostgreSQL password")
)

func init() {
	viper.SetDefault("PGPassword", "password")
}

func MigrateDB() {
	m, err := migrate.New(
		"file://db/migrations",
		fmt.Sprintf("postgres://postgres:%s@postgres:5432/sglapp?sslmode=disable", viper.GetString("PGPassword")))
	if err != nil {
		log.Fatalf("migrate.New() in MigrateDB failed %q\n", err)
	}
	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatalf("migrate.Up() in MigrateDB failed %q\n", err)
	}
}

func Init() {
	settings := postgresql.ConnectionURL{
		Host:     "postgres",
		Database: "sglapp",
		User:     "postgres",
		Password: viper.GetString("PGPassword"),
	}
	var err error
	Sess, err = postgresql.Open(settings)
	if err != nil {
		logrus.Fatalf("db.Open in InitDB %q\n", err)
	}
}
