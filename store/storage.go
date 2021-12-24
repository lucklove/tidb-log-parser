// Copyright 2021 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	Insert(fid, eid uint) error
	EventCount(eid uint) (uint, error)
	LogFragmentCount() (uint, error)
	Close() error
}

type sqliteDB struct {
	db     *sql.DB
	schema string
}

func NewSQLiteStorage(dbPath, schema string) (Storage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(fmt.Sprintf("create table if not exists %s(fid, eid, primary key(fid, eid))", schema)); err != nil {
		return nil, err
	}
	return &sqliteDB{db, schema}, nil
}

func (s *sqliteDB) Insert(fid, eid uint) error {
	if _, err := s.db.Exec(fmt.Sprintf("replace into %s(fid, eid) values(?, ?)", s.schema), fid, eid); err != nil {
		return err
	}
	return nil
}

func (s *sqliteDB) EventCount(eid uint) (uint, error) {
	rows, err := s.db.Query(fmt.Sprintf("select count(eid) from %s where eid = ?", s.schema), eid)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, rows.Err()
	}
	var count uint
	if err := rows.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *sqliteDB) LogFragmentCount() (uint, error) {
	rows, err := s.db.Query(fmt.Sprintf("select count(distinct(fid)) from %s", s.schema))
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, rows.Err()
	}
	var count uint
	if err := rows.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *sqliteDB) Close() error {
	return s.db.Close()
}
