package main

import (
	"database/sql"
	"errors"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	stmt := `INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)`
	res, err := s.db.Exec(stmt, p.Client, ParcelStatusRegistered, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	stmt := `SELECT number, client, status, address, created_at FROM parcel WHERE number = ?`
	var p Parcel
	err := s.db.QueryRow(stmt, number).Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	stmt := `SELECT number, client, status, address, created_at FROM parcel WHERE client = ?`
	rows, err := s.db.Query(stmt, client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Parcel
	for rows.Next() {
		var p Parcel
		if err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel

	stmt := `UPDATE parcel SET status = ? WHERE number = ?`
	_, err := s.db.Exec(stmt, status, number)
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	res, err := s.db.Exec(`
	UPDATE parcel SET address = ?
	WHERE number = ? AND status = ?`, address, number, ParcelStatusRegistered)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("cannot update address unless parcel is registered")
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	res, err := s.db.Exec(`DELETE FROM parcel WHERE number = ? AND status = ?`, number, ParcelStatusRegistered)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("cannot delete parcel unless it is registered")
	}
	return nil
}
