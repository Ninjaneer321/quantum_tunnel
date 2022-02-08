package libquantum

import (
	"fmt"
	"log"
	"os/user"
)

// Get the current users home directory
func GetHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot get home directory for current user: %s", err))
	}
	return usr.HomeDir
}
