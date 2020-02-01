package conf

import "syscall"

type Configuration struct {
	Database struct {
		Username string
		Password string
		Host string
		Database string
	}
}

func New() (*Configuration, error) {

}
