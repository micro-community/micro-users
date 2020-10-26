package repository

//passwd for saving into store
type passwd struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}
