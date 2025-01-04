package domain

import (
	"os"
	"seva/utils"
)

func GetDomains() ([]string, *utils.Error) {
	dir := "Var/Domains"

	// Read the directory
	files, be := os.ReadDir(dir)
	if be != nil {
		return nil, utils.NewDefaultErrorFromBase(be)
	}

	r := []string{}
	// Iterate over the files and directories
	for _, file := range files {
		// Check if the file is a directory
		if file.IsDir() {
			r = append(r, file.Name())
		}
	}
	return r, nil

}

func IsDomainRegistered(domain string) (bool, *utils.Error) {
	domains, e := GetDomains()
	if e != nil {
		return false, e
	}
	for _, d := range domains {
		if domain == d {
			return true, nil
		}
	}
	return false, nil
}

func RegisterDomain(domain string) *utils.Error {
	registered, e := IsDomainRegistered(domain)
	if e != nil {
		return e
	}
	if !registered {
		return utils.NewDefaultError("Domain already registered.")
	}

	be := os.Mkdir("Var/Domains/"+domain, 0755)
	if be != nil {
		return utils.NewDefaultErrorFromBase(be)
	}

	return nil
}
