package resiliency

// Define constants using the uint32 type
const (
	OK                  uint32 = 0
	CANCELED            uint32 = 1
	UNKNOWN             uint32 = 2
	INVALID_ARGUMENT    uint32 = 3
	DEADLINE_EXCEEDED   uint32 = 4
	NOT_FOUND           uint32 = 5
	ALREADY_EXISTS      uint32 = 6
	PERMISSION_DENIED   uint32 = 7
	RESOURCE_EXHAUSTED  uint32 = 8
	FAILED_PRECONDITION uint32 = 9
	ABORTED             uint32 = 10
	OUT_OF_RANGE        uint32 = 11
	UNIMPLEMENTED       uint32 = 12
	INTERNAL            uint32 = 13
	UNAVAILABLE         uint32 = 14
	DATA_LOSS           uint32 = 15
)
