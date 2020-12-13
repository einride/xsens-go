package xsens

// ErrorCode represents an Xsens error code.
type ErrorCode uint8

//go:generate stringer -type ErrorCode -trimprefix ErrorCode

const (
	// ErrorCodeOK: No error.
	ErrorCodeOK ErrorCode = 0

	// ErrorCodeNoBus: No bus communication possible.
	ErrorCodeNoBus ErrorCode = 1

	// ErrorCodeBusNotReady: InitBus and/or SetBID are not issued.
	ErrorCodeBusNotReady ErrorCode = 2

	// ErrorCodeInvalidPeriod: Period sent is invalid.
	ErrorCodeInvalidPeriod ErrorCode = 3

	// ErrorCodeInvalidMessage: Message is invalid or not implemented.
	ErrorCodeInvalidMessage ErrorCode = 4

	// ErrorCodeInitBusFail1: A slave did not respond to WaitForSetBID.
	ErrorCodeInitBusFail1 ErrorCode = 16

	// ErrorCodeInitBusFail2: An incorrect answer received after WaitForSetBID.
	ErrorCodeInitBusFail2 ErrorCode = 17

	// ErrorCodeInitBusFail3: After four bus-scans still undetected Motion Trackers.
	ErrorCodeInitBusFail3 ErrorCode = 18

	// ErrorCodeSetBIDFail1: No reply to SetBID message during SetBID procedure.
	ErrorCodeSetBIDFail1 ErrorCode = 20

	// ErrorCodeSetBIDFail2: Other than SetBIDAck received.
	ErrorCodeSetBIDFail2 ErrorCode = 21

	// ErrorCodeMeasurementFail1: Period too short to collect all data from Motion Trackers.
	ErrorCodeMeasurementFail1 ErrorCode = 24

	// ErrorCodeMeasurementFail2: Motion Tracker responds with other than SlaveData message.
	ErrorCodeMeasurementFail2 ErrorCode = 25

	// ErrorCodeMeasurementFail3: Total bytes of data of Motion Trackers exceeds 255 bytes.
	ErrorCodeMeasurementFail3 ErrorCode = 26

	// ErrorCodeMeasurementFail4: Timer overflows during measurement.
	ErrorCodeMeasurementFail4 ErrorCode = 27

	// ErrorCodeMeasurementFail5: Timer overflows during measurement.
	ErrorCodeMeasurementFail5 ErrorCode = 28

	// ErrorCodeMeasurementFail6: No correct response from Motion Tracker during measurement.
	ErrorCodeMeasurementFail6 ErrorCode = 29

	// ErrorCodeTimerOverflow: Timer overflow during measurement.
	ErrorCodeTimerOverflow ErrorCode = 30

	// ErrorCodeBaudrateInvalid: Baud rate does not comply with valid range.
	ErrorCodeBaudrateInvalid ErrorCode = 32

	// ErrorCodeInvalidParam: An invalid parameter is supplied.
	ErrorCodeInvalidParam ErrorCode = 33

	// ErrorCodeMeasurementFail7: TX PC Buffer is full.
	ErrorCodeMeasurementFail7 ErrorCode = 35

	// ErrorCodeMeasurementFail8: TX PC Buffer overflow, cannot fit full message.
	ErrorCodeMeasurementFail8 ErrorCode = 36

	// ErrorCodeDeviceError: Device generated an error, try updating the firmware.
	ErrorCodeDeviceError ErrorCode = 40

	// ErrorCodeDataOverflow: Device generates more data than bus can handle (baud rate too low).
	ErrorCodeDataOverflow ErrorCode = 41

	// ErrorCodeBufferOverflow: Sample buffer of the device was full during a communication outage.
	ErrorCodeBufferOverflow ErrorCode = 42
)
