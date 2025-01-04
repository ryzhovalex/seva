package domains

import (
	"os"
	"seva/lib/utils"
)

func GetDomains() ([]string, *utils.Error) {
	dir := "Var/Domains"

	// Read the directory
	files, be := os.ReadDir(dir)
	if be != nil {
		return nil, utils.CreateDefaultErrorFromBase(be)
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

func CheckDomainNotCreated(domain string) *utils.Error {
	registered, e := IsDomainCreated(domain)
	if e != nil {
		return e
	}
	if !registered {
		return utils.CreateDefaultError("Domain already registered: " + domain)
	}
	return nil
}

func CheckDomainCreated(domain string) *utils.Error {
	registered, e := IsDomainCreated(domain)
	if e != nil {
		return e
	}
	if !registered {
		return utils.CreateDefaultError("Domain not registered: " + domain)
	}
	return nil
}

func IsDomainCreated(domain string) (bool, *utils.Error) {
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

func GetDomainDir(domain string) string {
	return "Var/Domains/" + domain
}

func CreateDomain(domain string) *utils.Error {
	e := CheckDomainCreated(domain)
	if e != nil {
		return e
	}

	be := os.Mkdir("Var/Domains/"+domain, 0755)
	if be != nil {
		return utils.CreateDefaultErrorFromBase(be)
	}

	return nil
}
