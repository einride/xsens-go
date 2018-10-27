package xsensgo

import (
	"sync"
)

var imuStatus = IMUData{
	AltitudeMSL:     -1,
	AltitudeMEllips: -1,
	Quat:            XDIQuaternion{Q0: -1, Q1: -1, Q2: -1, Q3: -1},
	Euler:           XDIEulerAngles{Roll: -1, Pitch: -1, Yaw: -1},
	Vel:             XDIVelocityXYZ{VelX: -1, VelY: -1, VelZ: -1},
	Latlng:          XDILatLng{Lat: -1, Lng: -1},
	Acc:             XDIAccelerationXYZ{AccX: -1, AccY: -1, AccZ: -1},
	Mag:             XDIMagneticXYZ{MagX: -1, MagY: -1, MagZ: -1},
	AngularVel:      XDIRateOfTurnXYZ{GyrX: -1, GyrY: -1, GyrZ: -1},
	Timestamp:       XDIUTCTime{NS: 0, Year: 0, Month: 0, Day: 0, Hour: 0, Minute: 0, Second: 0, Conf: 0},
	FreeAcc:         XDIFreeAccelerationXYZ{FreeAccX: -1, FreeAccY: -1, FreeAccZ: -1},
}

type IMUData struct {
	AltitudeMSL     float64
	AltitudeMEllips float64
	Quat            XDIQuaternion
	Euler           XDIEulerAngles
	Vel             XDIVelocityXYZ
	Latlng          XDILatLng
	Acc             XDIAccelerationXYZ
	Mag             XDIMagneticXYZ
	AngularVel      XDIRateOfTurnXYZ
	Timestamp       XDIUTCTime
	FreeAcc         XDIFreeAccelerationXYZ
}

var mutex = &sync.Mutex{}

func SetStatusAtomic(callback func(data IMUData) IMUData) {
	mutex.Lock()
	defer mutex.Unlock()
	imuStatus = callback(imuStatus)
}

func GetStatus() (status IMUData) {
	mutex.Lock()
	defer mutex.Unlock()
	return imuStatus
}
