package auth

type LoginContent struct {
	Homeserver string `json:"homeserver"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	DeviceName string `json:"device_name" optional:"true"`
}

func Login(_ string, _ LoginContent) (any, error) {
	return nil, nil
}

func Logout(_ string, _ struct{}) (any, error) {
	return nil, nil
}
