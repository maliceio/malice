package main

import (
	"fmt"
	"net/http"
)

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

// Index root route
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

// VBoxVersion route returns VBoxManage version
func VBoxVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	d := vbox.NewDriver("", "")
	outPut, err := d.Version()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(outPut))
}

// VBoxList route lists all VMs
func VBoxList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// machines, err := virtualbox.ListMachines()
	// assert(err)
	// for _, machine := range machines {
	// 	fmt.Println(machine.Name)
	// }

	// if err := json.NewEncoder(w).Encode(machines); err != nil {
	// 	panic(err)
	// }
	d := vbox.NewDriver("", "")
	outPut, err := d.ListVMs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(outPut))
}
