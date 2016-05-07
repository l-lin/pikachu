package service

import (
	"github.com/l-lin/pikachu/db"
	"log"
	"database/sql"
)

// An instance of a service
type Instance struct {
	InstanceId     int    `json:"instanceId"`
	ServiceId      int    `json:"serviceId"`
	Name           string `json:"name"`
	UrlHealthCheck string `json:"urlHealthCheck"`
	Status         string `json:"status"`
}

func NewInstance() *Instance {
	return &Instance{}
}

// Get the instances of a service
func GetInstances(serviceId int) []*Instance {
	instances := make([]*Instance, 0)
	database := db.Connect()
	defer database.Close()

	rows, err := database.Query(`
		SELECT instance_id, service_id, name, url_health_check, status
		FROM pika_instance
		WHERE service_id = $1
	`, serviceId)
	if err != nil {
		log.Printf("[x] Error when getting the list of instances. Reason: %s", err.Error())
		return instances
	}
	for rows.Next() {
		instance := toInstance(rows)
		instances = append(instances, instance)
	}
	if err := rows.Err(); err != nil {
		log.Printf("[x] Error when getting the list of instances. Reason: %s", err.Error())
	}
	return instances
}

// Get the instance from a given instanceId
func GetInstance(instanceId int) *Instance {
	database := db.Connect()
	defer database.Close()

	row := database.QueryRow("SELECT instance_id, service_id, name, url_health_check, status FROM pika_instance WHERE instance_id = $1", instanceId)
	return toInstance(row)
}

// Save the instance in the db
func (i *Instance) Save() {
	database := db.Connect()
	defer database.Close()
	tx, err := database.Begin()
	if err != nil {
		log.Printf("[x] Could not start the transaction. Reason: %s", err.Error())
	}
	row := tx.QueryRow(`
		INSERT INTO pika_instance (service_id, name, url_health_check, status) VALUES ($1, $2, $3, $4) RETURNING instance_id
		`, i.ServiceId, i.Name, i.UrlHealthCheck, i.Status)
	var lastId int
	if err := row.Scan(&lastId); err != nil {
		tx.Rollback()
		log.Printf("[x] Could not fetch the instance_id of the newly created instance. Reason: %s", err.Error())
	}
	i.InstanceId = lastId
	if err := tx.Commit(); err != nil {
		log.Printf("[x] Could not commit the transaction. Reason: %s", err.Error())
	}
}

// Update the instance
func (i *Instance) Update() {
	database := db.Connect()
	defer database.Close()
	tx, err := database.Begin()
	if err != nil {
		log.Printf("[x] Could not start the transaction. Reason: %s", err.Error())
		return
	}
	_, err = tx.Exec(`
		UPDATE pika_instance SET name = $1, url_health_check = $2, status = $3 WHERE service_id = $4
	`, i.Name, i.UrlHealthCheck, i.Status, i.InstanceId)
	if err != nil {
		tx.Rollback()
		log.Printf("[x] Could not update the instane. Reason: %s", err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[x] Could not commit the transaction. Reason: %s", err.Error())
	}
}

// Delete an instance
func (i *Instance) Delete() {
	database := db.Connect()
	defer database.Close()
	tx, err := database.Begin()
	if err != nil {
		log.Fatalf("[x] Could not start the transaction. Reason: %s", err.Error())
	}

	_, err = tx.Exec("DELETE FROM pika_instance WHERE instance_id = $1", i.InstanceId)
	if err != nil {
		tx.Rollback()
		log.Printf("[x] Could not delete the instance. Reason: %s", err.Error())
		return
	}
	if err := tx.Commit(); err != nil {
		log.Printf("[x] Could not commit the transaction. Reason: %s", err.Error())
	}
}

// Delete the instance of a service
func DeleteInstances(serviceId int) {
	database := db.Connect()
	defer database.Close()
	tx, err := database.Begin()
	if err != nil {
		log.Fatalf("[x] Could not start the transaction. Reason: %s", err.Error())
	}

	_, err = tx.Exec("DELETE FROM pika_instance WHERE service_id = $1", serviceId)
	if err != nil {
		tx.Rollback()
		log.Printf("[x] Could not delete the instances. Reason: %s", err.Error())
		return
	}
	if err := tx.Commit(); err != nil {
		log.Printf("[x] Could not commit the transaction. Reason: %s", err.Error())
	}
}

func toInstance(rows db.RowMapper) *Instance {
	i := NewInstance()
	err := rows.Scan(
		&i.InstanceId,
		&i.ServiceId,
		&i.Name,
		&i.UrlHealthCheck,
		&i.Status,
	)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("[-] No instance found")
		return nil
	case err != nil:
		log.Printf("[-] Could not scan the instance. Reason: %s", err.Error())
		return nil
	default:
		return i
	}
}
