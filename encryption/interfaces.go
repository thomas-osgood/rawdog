package encryption

// interface describing an encryptiion object
// that can be used with the Rawdog system.
type RawdogEncryptor interface {
	Decrypt([]byte) ([]byte, error)
	Encrypt([]byte) ([]byte, error)
}
