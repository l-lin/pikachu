package service

import (
	"database/sql"
	"github.com/l-lin/pikachu/db"
	"log"
)

// A service contains a list of instances that do a specific business purpose
type Service struct {
	ServiceId int         `json:"serviceId"`
	Name      string      `json:"name"`
	Instances []*Instance `json:"instances"`
}

func NewService() *Service {
	s := &Service{}
	s.Instances = make([]*Instance, 0)
	return s
}

func (s *Service) AddInstance(i *Instance) {
	s.Instances = append(s.Instances, i)
}

// Get all services
func GetServices() []*Service {
	services := make([]*Service, 0)
	database := db.Connect()
	defer database.Close()

	rows, err := database.Query(`
		SELECT service_id, name
		FROM pika_service
	`)
	if err != nil {
		log.Printf("[x] Error when getting the list of services. Reason: %s", err.Error())
		return services
	}
	for rows.Next() {
		s := toService(rows)
		services = append(services, s)
	}
	if err := rows.Err(); err != nil {
		log.Printf("[x] Error when getting the list of services. Reason: %s", err.Error())
	}

	if len(services) > 0 {
		for _, s := range services {
			s.Instances = GetInstances(s.ServiceId)
		}
	}

	return services
}

// Get the service from a given serviceID
func GetService(serviceId int) *Service {
	database := db.Connect()
	defer database.Close()

	row := database.QueryRow("SELECT service_id, name FROM pika_service WHERE service_id = $1", serviceId)
	s := toService(row)
	if s != nil {
		s.Instances = GetInstances(serviceId)
	}
	return s
}

// Save the service in the db
func (s *Service) Save() {
	database := db.Connect()
	defer database.Close()
	tx, err := database.Begin()
	if err != nil {
		log.Printf("[x] Could not start the transaction. Reason: %s", err.Error())
	}
	row := tx.QueryRow("INSERT INTO pika_service (name) VALUES ($1) RETURNING service_id", s.Name)
	var lastId int
	if err := row.Scan(&lastId); err != nil {
		tx.Rollback()
		log.Printf("[x] Could not fetch the service_id of the newly created service. Reason: %s", err.Error())
	}
	s.ServiceId = lastId

	if s.Instances != nil && len(s.Instances) > 0 {
		for _, i := range s.Instances {
			i.Save()
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[x] Could not commit the transaction. Reason: %s", err.Error())
	}
}

// Update the service
func (s *Service) Update() {
	database := db.Connect()
	defer database.Close()
	tx, err := database.Begin()
	if err != nil {
		log.Printf("[x] Could not start the transaction. Reason: %s", err.Error())
		return
	}
	_, err = tx.Exec("UPDATE pika_service SET name = $1 WHERE service_id = $2", s.Name, s.ServiceId)
	if err != nil {
		tx.Rollback()
		log.Printf("[x] Could not update the service. Reason: %s", err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[x] Could not commit the transaction. Reason: %s", err.Error())
	}
}

// Delete a service
func (s *Service) Delete() {
	database := db.Connect()
	defer database.Close()
	tx, err := database.Begin()
	if err != nil {
		log.Fatalf("[x] Could not start the transaction. Reason: %s", err.Error())
	}

	DeleteInstances(s.ServiceId)

	_, err = tx.Exec("DELETE FROM pika_service WHERE service_id = $1", s.ServiceId)
	if err != nil {
		tx.Rollback()
		log.Printf("[x] Could not delete the service. Reason: %s", err.Error())
		return
	}
	if err := tx.Commit(); err != nil {
		log.Printf("[x] Could not commit the transaction. Reason: %s", err.Error())
	}
}

func toService(rows db.RowMapper) *Service {
	s := NewService()
	err := rows.Scan(
		&s.ServiceId,
		&s.Name,
	)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("[-] No service found")
		return nil
	case err != nil:
		log.Printf("[-] Could not scan the service. Reason: %s", err.Error())
		return nil
	default:
		return s
	}
}
