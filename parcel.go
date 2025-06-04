package main

import (
	"database/sql"
	"errors"
	"time"
)

type Parcel struct {
	Number    int
	Client    int
	Status    string
	Address   string
	CreatedAt string
}

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	stmt := `INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)`
	now := time.Now().Format(time.RFC3339)
	p.CreatedAt = now
	res, err := s.db.Exec(stmt, p.Client, ParcelStatusRegistered, p.Address, now)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
	// верните идентификатор последней добавленной записи

}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка

	// заполните объект Parcel данными из таблицы
	//p := Parcel{} ?
	stmt := `SELECT number, client, status, address, created_at FROM parcel WHERE number = ?`
	var p Parcel
	err := s.db.QueryRow(stmt, number).Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	return p, err

}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	stmt := `SELECT number, client, status, address, created_at FROM parcel WHERE client = ?`
	rows, err := s.db.Query(stmt, client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		var p Parcel
		if err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, p)
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
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	var status string
	err := s.db.QueryRow(`SELECT status FROM parcel WHERE number = ?`, number).Scan(&status)
	if err != nil {
		return err
	}
	if status != ParcelStatusRegistered {
		return errors.New("cannot update address unless parcel is registered")
	}
	_, err = s.db.Exec(`UPDATE parcel SET address = ? WHERE number = ?`, address, number)
	return err
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered

	var status string
	err := s.db.QueryRow(`SELECT status FROM parcel WHERE number = ?`, number).Scan(&status)
	if err != nil {
		return err
	}
	if status != ParcelStatusRegistered {
		return errors.New("cannot delete parcel unless it is registered")
	}
	_, err = s.db.Exec(`DELETE FROM parcel WHERE number = ?`, number)
	return err
}
