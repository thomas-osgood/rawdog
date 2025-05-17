package noencrypt

// function designed to return a new instance
// of a BlankEncryptor object.
func New() (*BlankEncryptor, error) {
	return &BlankEncryptor{}, nil
}
