package api

// TRISA transfer state constants. See protocol buffers documentation for more details.
const (
	TransferStateUnspecified = TransferState_UNSPECIFIED
	TransferStarted          = TransferState_STARTED
	TransferPending          = TransferState_PENDING
	TransferRepair           = TransferState_REPAIR
	TransferAccepted         = TransferState_ACCEPTED
	TransferCompleted        = TransferState_COMPLETED
	TransferRejected         = TransferState_REJECTED
)
