package noencrypt

// function designed to mimic a Decrypt function.
//
// because this is a BlankEncryptor object, the ciphertext
// passed in will be immediately returned, without anything
// being done to it.
func (be *BlankEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	return ciphertext, nil
}

// function designed to mimic an Encrypt function.
//
// because this is a BlankEncryptor object, the plaintext
// passed in will be immediately returned, without anything
// being done to it.
func (be *BlankEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
	return plaintext, nil
}
