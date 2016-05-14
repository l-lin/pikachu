package web

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/l-lin/pikachu/entity"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Handler to fetch the services
func GetServices(w http.ResponseWriter, r *http.Request) {
	log.Printf("[-] Fetching services")
	write(w, http.StatusOK, entity.GetServices())
}

// Handler to fetch a service
func GetService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceId, _ := strconv.Atoi(vars["serviceId"])
	log.Printf("[-] Fetching service %d", serviceId)

	s := entity.GetService(serviceId)
	if s != nil {
		write(w, http.StatusOK, s)
		return
	}
	// If we didn't find it, 404
	log.Printf("[-] Could not find the service %d", serviceId)
	write(w, http.StatusNotFound, JsonErr{Code: http.StatusNotFound, Text: fmt.Sprintf("Service not found for serviceId %d", serviceId)})
}

// Handler to save a service
func SaveService(w http.ResponseWriter, r *http.Request) {
	var s entity.Service
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Printf("[x] Could not read the body. Reason: %s", err.Error())
		write(w, http.StatusInternalServerError, JsonErr{Code: http.StatusInternalServerError, Text: "Could not read the body."})
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Printf("[x] Could not close ready the body. Reason: %s", err.Error())
		write(w, http.StatusInternalServerError, JsonErr{Code: http.StatusInternalServerError, Text: "Could not close the body."})
		return
	}
	if err := json.Unmarshal(body, &s); err != nil {
		// 422: unprocessable entity
		write(w, 422, JsonErr{Code: 422, Text: "Could not parse the given parameter"})
		return
	}

	log.Printf("[-] Creating new service")
	s.Save()
	write(w, http.StatusCreated, s)
}

// Handler to update a service
func UpdateService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceId, _ := strconv.Atoi(vars["serviceId"])

	oldService := entity.GetService(serviceId)
	if oldService == nil {
		log.Printf("[-] Could not find the service %d", serviceId)
		write(w, http.StatusNotFound, JsonErr{Code: http.StatusNotFound, Text: fmt.Sprintf("Service not found for serviceId %d", serviceId)})
		return
	}

	var s entity.Service
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatalf("[x] Could not read the body. Reason: %s", err.Error())
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalf("[x] Could not close ready the body. Reason: %s", err.Error())
	}
	if err := json.Unmarshal(body, &s); err != nil {
		// 422: unprocessable entity
		write(w, 422, JsonErr{Code: 422, Text: "Could not parse the given parameter"})
		return
	}
	s.ServiceId = serviceId

	log.Printf("[-] Updating service %d", serviceId)
	s.Update()
	write(w, http.StatusOK, s)
}

// Handler to delete a service
func DeleteService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceId, _ := strconv.Atoi(vars["serviceId"])
	s := entity.NewService()
	s.ServiceId = serviceId

	log.Printf("[-] Deleting service id %d", serviceId)
	s.Delete()
	write(w, http.StatusNoContent, nil)
}

// Handler to fetch the instances of a service
func GetInstances(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceId, _ := strconv.Atoi(vars["serviceId"])

	instances := entity.GetInstances(serviceId)
	if instances != nil && len(instances) > 0 {
		log.Printf("[-] Found %d instances of serviceId %d", len(instances), serviceId)
		write(w, http.StatusOK, instances)
		return
	}

	// If we didn't find it, 404
	log.Printf("[-] Could not find the instances for serviceId %d", serviceId)
	write(w, http.StatusNotFound, JsonErr{Code: http.StatusNotFound, Text: fmt.Sprintf("Service not Found for serviceId %d", serviceId)})
}

// Handler to save an instance
func SaveInstance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceId, _ := strconv.Atoi(vars["serviceId"])

	var i entity.Instance
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Printf("[x] Could not read the body. Reason: %s", err.Error())
		write(w, http.StatusInternalServerError, JsonErr{Code: http.StatusInternalServerError, Text: "Could not read the body."})
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Printf("[x] Could not close ready the body. Reason: %s", err.Error())
		write(w, http.StatusInternalServerError, JsonErr{Code: http.StatusInternalServerError, Text: "Could not close the body."})
		return
	}
	if err := json.Unmarshal(body, &i); err != nil {
		// 422: unprocessable entity
		write(w, 422, JsonErr{Code: 422, Text: "Could not parse the given parameter"})
		return
	}
	i.ServiceId = serviceId

	log.Printf("[-] Creating new instance")
	i.Save()
	write(w, http.StatusCreated, i)
}

// Handler to fetch an instance
func GetInstance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceId, _ := strconv.Atoi(vars["instanceId"])

	instance := entity.GetInstance(instanceId)
	if instance != nil {
		log.Printf("[-] Found instances %d", instanceId)
		write(w, http.StatusOK, instance)
		return
	}

	// If we didn't find it, 404
	log.Printf("[-] Could not find the instances %d", instanceId)
	write(w, http.StatusNotFound, JsonErr{Code: http.StatusNotFound, Text: fmt.Sprintf("Instance not found for instanceId %d", instanceId)})
}

// Handler to update an instance
func UpdateInstance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceId, _ := strconv.Atoi(vars["instanceId"])

	oldInstance := entity.GetInstance(instanceId)
	if oldInstance == nil {
		log.Printf("[-] Could not find the instances for instanceId %d", instanceId)
		write(w, http.StatusNotFound, JsonErr{Code: http.StatusNotFound, Text: fmt.Sprintf("Instance not found for instanceId %d", instanceId)})
		return
	}

	var i entity.Instance
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatalf("[x] Could not read the body. Reason: %s", err.Error())
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalf("[x] Could not close ready the body. Reason: %s", err.Error())
	}
	if err := json.Unmarshal(body, &i); err != nil {
		// 422: unprocessable entity
		write(w, 422, JsonErr{Code: 422, Text: "Could not parse the given parameter"})
		return
	}
	i.InstanceId = instanceId

	log.Printf("[-] Updating instance %d", instanceId)
	i.Update()
	write(w, http.StatusOK, i)
}

// Handler to delete an instance
func DeleteInstance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceId, _ := strconv.Atoi(vars["instanceId"])
	i := entity.NewInstance()
	i.InstanceId = instanceId

	log.Printf("[-] Deleting instance id %d", instanceId)
	i.Delete()
	write(w, http.StatusNoContent, nil)
}

// Write the response in JSON Content-type
func write(w http.ResponseWriter, status int, n interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	if n != nil {
		if err := json.NewEncoder(w).Encode(n); err != nil {
			panic(err)
		}
	}
}
