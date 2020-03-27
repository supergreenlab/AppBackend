/*
 * Copyright (C) 2019  SuperGreenLab <towelie@supergreenlab.com>
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

package main

import (
	"fmt"

	"github.com/SuperGreenLab/AppBackend/internal/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	pgPassword = pflag.String("pgpassword", "password", "PostgreSQL password")
)

func initDB() {
	m, err := migrate.New(
		"file://db/migrations",
		fmt.Sprintf("postgres://postgres:%s@postgres:5432/sglapp?sslmode=disable", viper.GetString("PGPassword")))
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}
}

func main() {
	viper.SetConfigName("appbackend")
	viper.AddConfigPath("/etc/appbackend")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	viper.SetEnvPrefix("APPBACKEND")
	viper.AutomaticEnv()

	viper.SetDefault("PGPassword", "password")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	initDB()

	server.Start()

	log.Info("AppBackend started")

	select {}
}
