package docker

// Sample Virtualbox create independent of Machine CLI.
// import (
// 	"encoding/json"
// 	"fmt"

// 	log "github.com/Sirupsen/logrus"
// 	"github.com/docker/machine/commands/mcndirs"
// 	"github.com/docker/machine/drivers/virtualbox"
// 	"github.com/docker/machine/libmachine"
// 	er "github.com/maliceio/malice/malice/errors"
// )

// // MakeDockerMachine creates a new docker host via docker-machine
// func MakeDockerMachine(host string) {
// 	// log.SetDebug(true)

// 	client := libmachine.NewClient(mcndirs.GetBaseDir(), mcndirs.GetMachineCertDir())

// 	hostName := host

// 	// Set some options on the provider...
// 	driver := virtualbox.NewDriver(hostName, mcndirs.GetBaseDir())
// 	driver.CPU = 2
// 	driver.Memory = 2048

// 	data, err := json.Marshal(driver)
// 	er.CheckError(err)

// 	// pluginDriver, err := client.NewPluginDriver("virtualbox", data)
// 	// er.CheckError(err)

// 	h, err := client.NewHost("virtualbox", data)
// 	// h, err := client.NewHost(pluginDriver)
// 	er.CheckError(err)

// 	h.HostOptions.EngineOptions.StorageDriver = "overlay"

// 	if err := client.Create(h); err != nil {
// 		log.Fatal(err)
// 	}

// 	out, err := h.RunSSHCommand("df -h")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Printf("Results of your disk space query:\n%s\n", out)

// 	fmt.Println("Powering down machine now...")
// 	if err := h.Stop(); err != nil {
// 		log.Fatal(err)
// 	}
// }

// // MachineURL returns the IP of the docker-machine
// func MachineURL(name string) (url string, err error) {

// 	api := libmachine.NewClient(mcndirs.GetBaseDir(), mcndirs.GetMachineCertDir())

// 	host, err := api.Load(name)
// 	er.CheckError(err)
// 	url, err = host.URL()
// 	er.CheckError(err)

// 	return
// }

// // MachineIP returns the IP of the docker-machine
// func MachineIP(name string) (ip string, err error) {

// 	api := libmachine.NewClient(mcndirs.GetBaseDir(), mcndirs.GetMachineCertDir())

// 	host, err := api.Load(name)
// 	er.CheckError(err)
// 	ip, err = host.Driver.GetIP()
// 	er.CheckError(err)

// 	return
// }

// // MachineStop stops the docker-machine
// func MachineStop(name string) error {

// 	api := libmachine.NewClient(mcndirs.GetBaseDir(), mcndirs.GetMachineCertDir())

// 	host, err := api.Load(name)
// 	er.CheckError(err)
// 	err = host.Driver.Stop()

// 	return err
// }
