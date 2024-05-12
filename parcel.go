package main

import (
	"database/sql"
	"log"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуем добавление строки в таблицу parcel, используя данные из переменной p
	// в res (объект sql.Result) через s.db.Exec инсертнем данные из p, кроме number (сам присвоится)
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :createdAt)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}
	// возвращаем идентификатор последней добавленной записи через LastInsertId
	// помним, что LastInsertId возвращает int64, err, а у нас (int, error)
	// так что приведем int64 к int
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуем чтение строки по заданному number через s.db.QueryRow
	// здесь из таблицы должна вернуться только одна строка так что пользоваться будем QueryRow
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number", sql.Named("number", number))
	// заполняем объект Parcel данными из таблицы
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуем чтение строк из таблицы parcel по заданному client через s.db.Query
	// здесь из таблицы может вернуться несколько строк так что пользоваться будем Query
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	// заполняем срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
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
	// реализуем обновление статуса в таблице parcel через s.db.Exec
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуем обновление адреса в таблице parcel через s.db.Exec
	// менять адрес можно только если значение статуса registered
	// const из main.go тут не имеют власти, так что добавим переменную ParcelStatusRegistered
	ParcelStatusRegistered := "registered"
	_, err := s.db.Exec("UPDATE parcel SET address = :address  WHERE number = :number AND status = :status",
		sql.Named("address ", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуем удаление строки из таблицы parcel через s.db.Exec
	// удалять строку можно только если значение статуса registered
	// const из main.go тут не имеют власти, так что добавим переменную ParcelStatusRegistered
	ParcelStatusRegistered := "registered"
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number, status = :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
