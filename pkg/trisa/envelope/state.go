package envelope

type State uint16

const (
	Unknown       State = iota
	Clear               // The envelope has been decrypted and the payload is available on it
	Unsealed            // The envelope is unsealed and can be decrypted without any other information
	Sealed              // The envelope is sealed and must be unsealed with a private key or it is ready to send
	Error               // The envelope does not contain a payload but does contain an error field
	ClearError          // The envelope contains both a decrypted payload and an error
	UnsealedError       // The envelope contains both an error and a payload and is unsealed
	SealedError         // The envelope contains both an error and a payload and is sealed
	Corrupted           // The envelope is in an invalid state and cannot be moved into a correct state
)

var stateNames = []string{"unknown", "clear", "unsealed", "sealed", "error", "clear-error", "unsealed-error", "sealed-error", "corrupted"}

func (s State) String() string {
	idx := int(s)
	if idx >= len(stateNames) {
		idx = 0
	}
	return stateNames[idx]
}
