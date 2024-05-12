package main

import (
	"database/sql"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	// настраиваю подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
		log.Println(err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавляем новую посылку в БД, убеждаемся в отсутствии ошибки и наличии идентификатора
	parcel.Number, err = store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, parcel.Number)

	// get
	// получаем только что добавленную посылку по идентификатору parcel.Number, убеждаемся в отсутствии ошибки
	// проверяем что значения всех полей в полученном объекте stored совпадают со значениями полей в переменной parcel,
	// то есть сравниваем их через assert.Equal
	stored, err := store.Get(parcel.Number)
	require.NoError(t, err)
	assert.Equal(t, parcel, stored)

	// delete
	// удаляем добавленную посылку по идентификатору parcel.Number, убеждаемся в отсутствии ошибки
	err = store.Delete(parcel.Number)
	require.NoError(t, err)
	// проверяем что посылку больше нельзя получить из БД
	// мы также пытаемся получить добавленную посылку по идентификатору parcel.Number через store.Get,
	// которая возвращает объект (как в примере выше stored) или err, объекта не будет,
	// соответственно мы сравним err c sql.ErrNoRows, которая есть в случае невозврата ни одной строки
	// так мы убедимся что посылку теперь не получить из БД
	_, err = store.Get(parcel.Number)
	require.Equal(t, sql.ErrNoRows, err)

}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	// настраиваю подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
		log.Println(err)
	}
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавляем новую посылку в БД, убеждаемся в отсутствии ошибки и наличии идентификатора
	parcel.Number, err = store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, parcel.Number)

	// set address
	// обновляем адрес по идентификатору parcel.Number через store.SetAddress, убеждаемся в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(parcel.Number, newAddress)
	require.NoError(t, err)

	// check
	// получаем добавленную посылку по идентификатору parcel.Number
	// проверяем на отсутствие ошибки и убеждаемся, что адрес обновился
	// через сравнение (используя assert.Equal) stored.Address и newAddress
	stored, err := store.Get(parcel.Number)
	require.NoError(t, err)
	assert.Equal(t, stored.Address, newAddress)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	// настраиваем подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
		log.Println(err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавляем новую посылку в БД, убеждаемся в отсутствии ошибки и наличии идентификатора
	parcel.Number, err = store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, parcel.Number)

	// set status
	// обновляем статус по идентификатору parcel.Number через store.SetStatus, убеждаемся в отсутствии ошибки
	newStatus := "new test status"
	err = store.SetStatus(parcel.Number, newStatus)
	require.NoError(t, err)

	// check
	// получаем добавленную посылку по идентификатору parcel.Number через store.Get,
	// убеждаемся, что статус обновился
	// через сравнение (используя assert.Equal) stored.Status и newStatus
	stored, err := store.Get(parcel.Number)
	require.NoError(t, err)
	assert.Equal(t, stored.Status, newStatus)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	// настраиваем подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		require.NoError(t, err)
		log.Println(err)
	}
	defer db.Close()
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		// добавляем новую посылку (parcels[i]) в БД через store.Add,
		// убеждаемся в отсутствии ошибки и наличии идентификатора id
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)
		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	// получаем список посылок по идентификатору клиента, сохранённого в переменной client через store.GetByClient
	// проверяем на отсутвие ошибки
	// проверяем что количество полученных посылок len(storedParcels) совпадает с количеством добавленных len(parcelMap) через assert.Equal
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Equal(t, len(parcelMap), len(storedParcels))

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		newParcel, exist := parcelMap[parcel.Number]
		// проверяем что все посылки из storedParcels есть в parcelMap через assert.True и exist
		assert.True(t, exist)
		// проверяем  что значения полей полученных посылок заполнены верно
		// через сравнение (assert.Equal) newParcel и parcel
		assert.Equal(t, newParcel, parcel)
	}
}
