// Code generated by "stringer -type MessageIdentifier -trimprefix MessageIdentifier"; DO NOT EDIT.

package xsens

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[MessageIdentifierWakeup-62]
	_ = x[MessageIdentifierWakeupAck-63]
	_ = x[MessageIdentifierReqDID-0]
	_ = x[MessageIdentifierDeviceID-1]
	_ = x[MessageIdentifierInitBus-2]
	_ = x[MessageIdentifierInitBusResults-3]
	_ = x[MessageIdentifierReqPeriod-4]
	_ = x[MessageIdentifierReqPeriodAck-5]
	_ = x[MessageIdentifierSetPeriod-4]
	_ = x[MessageIdentifierSetPeriodAck-5]
	_ = x[MessageIdentifierSetBid-6]
	_ = x[MessageIdentifierSetBidAck-7]
	_ = x[MessageIdentifierAutoStart-6]
	_ = x[MessageIdentifierAutoStartAck-7]
	_ = x[MessageIdentifierBusPower-8]
	_ = x[MessageIdentifierBusPowerAck-9]
	_ = x[MessageIdentifierReqDataLength-10]
	_ = x[MessageIdentifierDataLength-11]
	_ = x[MessageIdentifierReqConfiguration-12]
	_ = x[MessageIdentifierConfiguration-13]
	_ = x[MessageIdentifierRestoreFactoryDef-14]
	_ = x[MessageIdentifierRestoreFactoryDefAck-15]
	_ = x[MessageIdentifierGotoMeasurement-16]
	_ = x[MessageIdentifierGotoMeasurementAck-17]
	_ = x[MessageIdentifierReqFirmwareRevision-18]
	_ = x[MessageIdentifierFirmwareRevision-19]
	_ = x[MessageIdentifierReqBluetoothDisable-20]
	_ = x[MessageIdentifierReqBluetoothDisableAck-21]
	_ = x[MessageIdentifierDisableBluetooth-20]
	_ = x[MessageIdentifierDisableBluetoothAck-21]
	_ = x[MessageIdentifierReqXmOutputMode-22]
	_ = x[MessageIdentifierReqXmOutputModeAck-23]
	_ = x[MessageIdentifierSetXmOutputMode-22]
	_ = x[MessageIdentifierSetXmOutputModeAck-23]
	_ = x[MessageIdentifierReqBaudrate-24]
	_ = x[MessageIdentifierReqBaudrateAck-25]
	_ = x[MessageIdentifierSetBaudrate-24]
	_ = x[MessageIdentifierSetBaudrateAck-25]
	_ = x[MessageIdentifierReqSyncMode-26]
	_ = x[MessageIdentifierReqSyncModeAck-27]
	_ = x[MessageIdentifierSetSyncMode-26]
	_ = x[MessageIdentifierSetSyncModeAck-27]
	_ = x[MessageIdentifierReqProductCode-28]
	_ = x[MessageIdentifierProductCode-29]
	_ = x[MessageIdentifierReqProcessingFlags-32]
	_ = x[MessageIdentifierReqProcessingFlagsAck-33]
	_ = x[MessageIdentifierSetProcessingFlags-32]
	_ = x[MessageIdentifierSetProcessingFlagsAck-33]
	_ = x[MessageIdentifierReqInputTrigger-38]
	_ = x[MessageIdentifierReqInputTriggerAck-39]
	_ = x[MessageIdentifierSetInputTrigger-38]
	_ = x[MessageIdentifierSetInputTriggerAck-39]
	_ = x[MessageIdentifierReqOutputTrigger-40]
	_ = x[MessageIdentifierReqOutputTriggerAck-41]
	_ = x[MessageIdentifierSetOutputTrigger-40]
	_ = x[MessageIdentifierSetOutputTriggerAck-41]
	_ = x[MessageIdentifierSetSyncBoxMode-42]
	_ = x[MessageIdentifierSetSyncBoxModeAck-43]
	_ = x[MessageIdentifierReqSyncBoxMode-42]
	_ = x[MessageIdentifierReqSyncBoxModeAck-43]
	_ = x[MessageIdentifierSetSyncConfiguration-44]
	_ = x[MessageIdentifierSetSyncConfigurationAck-45]
	_ = x[MessageIdentifierReqSyncConfiguration-44]
	_ = x[MessageIdentifierSyncConfiguration-45]
	_ = x[MessageIdentifierDriverDisconnect-46]
	_ = x[MessageIdentifierDriverDisconnectAck-47]
	_ = x[MessageIdentifierXmPowerOff-68]
	_ = x[MessageIdentifierReqOutputConfiguration-192]
	_ = x[MessageIdentifierReqOutputConfigurationAck-193]
	_ = x[MessageIdentifierSetOutputConfiguration-192]
	_ = x[MessageIdentifierSetOutputConfigurationAck-193]
	_ = x[MessageIdentifierReqOutputMode-208]
	_ = x[MessageIdentifierReqOutputModeAck-209]
	_ = x[MessageIdentifierSetOutputMode-208]
	_ = x[MessageIdentifierSetOutputModeAck-209]
	_ = x[MessageIdentifierReqOutputSettings-210]
	_ = x[MessageIdentifierReqOutputSettingsAck-211]
	_ = x[MessageIdentifierSetOutputSettings-210]
	_ = x[MessageIdentifierSetOutputSettingsAck-211]
	_ = x[MessageIdentifierReqOutputSkipFactor-212]
	_ = x[MessageIdentifierReqOutputSkipFactorAck-213]
	_ = x[MessageIdentifierSetOutputSkipFactor-212]
	_ = x[MessageIdentifierSetOutputSkipFactorAck-213]
	_ = x[MessageIdentifierReqSyncInSettings-214]
	_ = x[MessageIdentifierReqSyncInSettingsAck-215]
	_ = x[MessageIdentifierSetSyncInSettings-214]
	_ = x[MessageIdentifierSetSyncInSettingsAck-215]
	_ = x[MessageIdentifierReqSyncOutSettings-216]
	_ = x[MessageIdentifierReqSyncOutSettingsAck-217]
	_ = x[MessageIdentifierSetSyncOutSettings-216]
	_ = x[MessageIdentifierSetSyncOutSettingsAck-217]
	_ = x[MessageIdentifierReqErrorMode-218]
	_ = x[MessageIdentifierReqErrorModeAck-219]
	_ = x[MessageIdentifierSetErrorMode-218]
	_ = x[MessageIdentifierSetErrorModeAck-219]
	_ = x[MessageIdentifierReqTransmitDelay-220]
	_ = x[MessageIdentifierReqTransmitDelayAck-221]
	_ = x[MessageIdentifierSetTransmitDelay-220]
	_ = x[MessageIdentifierSetTransmitDelayAck-221]
	_ = x[MessageIdentifierSetMfmResults-222]
	_ = x[MessageIdentifierSetMfmResultsAck-223]
	_ = x[MessageIdentifierReqObjectAlignment-224]
	_ = x[MessageIdentifierReqObjectAlignmentAck-225]
	_ = x[MessageIdentifierSetObjectAlignment-224]
	_ = x[MessageIdentifierSetObjectAlignmentAck-225]
	_ = x[MessageIdentifierReqXmErrorMode-130]
	_ = x[MessageIdentifierReqXmErrorModeAck-131]
	_ = x[MessageIdentifierSetXmErrorMode-130]
	_ = x[MessageIdentifierSetXmErrorModeAck-131]
	_ = x[MessageIdentifierReqBufferSize-132]
	_ = x[MessageIdentifierReqBufferSizeAck-133]
	_ = x[MessageIdentifierSetBufferSize-132]
	_ = x[MessageIdentifierSetBufferSizeAck-133]
	_ = x[MessageIdentifierReqHeading-130]
	_ = x[MessageIdentifierReqHeadingAck-131]
	_ = x[MessageIdentifierSetHeading-130]
	_ = x[MessageIdentifierSetHeadingAck-131]
	_ = x[MessageIdentifierReqMagneticField-106]
	_ = x[MessageIdentifierReqMagneticFieldAck-107]
	_ = x[MessageIdentifierSetMagneticField-106]
	_ = x[MessageIdentifierSetMagneticFieldAck-107]
	_ = x[MessageIdentifierReqLocationID-132]
	_ = x[MessageIdentifierReqLocationIDAck-133]
	_ = x[MessageIdentifierSetLocationID-132]
	_ = x[MessageIdentifierSetLocationIDAck-133]
	_ = x[MessageIdentifierReqExtOutputMode-134]
	_ = x[MessageIdentifierReqExtOutputModeAck-135]
	_ = x[MessageIdentifierSetExtOutputMode-134]
	_ = x[MessageIdentifierSetExtOutputModeAck-135]
	_ = x[MessageIdentifierReqBatteryLevel-136]
	_ = x[MessageIdentifierBatterylevel-137]
	_ = x[MessageIdentifierReqInitTrackMode-136]
	_ = x[MessageIdentifierReqInitTrackModeAck-137]
	_ = x[MessageIdentifierSetInitTrackMode-136]
	_ = x[MessageIdentifierSetInitTrackModeAck-137]
	_ = x[MessageIdentifierReqMasterSettings-138]
	_ = x[MessageIdentifierMasterSettings-139]
	_ = x[MessageIdentifierStoreFilterState-138]
	_ = x[MessageIdentifierStoreFilterStateAck-139]
	_ = x[MessageIdentifierSetUtcTime-96]
	_ = x[MessageIdentifierReqUtcTime-96]
	_ = x[MessageIdentifierSetUtcTimeAck-97]
	_ = x[MessageIdentifierUtcTime-97]
	_ = x[MessageIdentifierAdjustUtcTime-168]
	_ = x[MessageIdentifierAdjustUtcTimeAck-169]
	_ = x[MessageIdentifierReqActiveClockCorrection-156]
	_ = x[MessageIdentifierActiveClockCorrection-157]
	_ = x[MessageIdentifierStoreActiveClockCorrection-158]
	_ = x[MessageIdentifierStoreActiveClockCorrectionAck-159]
	_ = x[MessageIdentifierReqAvailableFilterProfiles-98]
	_ = x[MessageIdentifierAvailableFilterProfiles-99]
	_ = x[MessageIdentifierReqFilterProfile-100]
	_ = x[MessageIdentifierReqFilterProfileAck-101]
	_ = x[MessageIdentifierSetFilterProfile-100]
	_ = x[MessageIdentifierSetFilterProfileAck-101]
	_ = x[MessageIdentifierReqGravityMagnitude-102]
	_ = x[MessageIdentifierReqGravityMagnitudeAck-103]
	_ = x[MessageIdentifierSetGravityMagnitude-102]
	_ = x[MessageIdentifierSetGravityMagnitudeAck-103]
	_ = x[MessageIdentifierReqGpsLeverArm-104]
	_ = x[MessageIdentifierReqGpsLeverArmAck-105]
	_ = x[MessageIdentifierSetGpsLeverArm-104]
	_ = x[MessageIdentifierSetGpsLeverArmAck-105]
	_ = x[MessageIdentifierReqLatLonAlt-110]
	_ = x[MessageIdentifierReqLatLonAltAck-111]
	_ = x[MessageIdentifierSetLatLonAlt-110]
	_ = x[MessageIdentifierSetLatLonAltAck-111]
	_ = x[MessageIdentifierGotoConfig-48]
	_ = x[MessageIdentifierGotoConfigAck-49]
	_ = x[MessageIdentifierBusData-50]
	_ = x[MessageIdentifierMtData-50]
	_ = x[MessageIdentifierSetNoRotation-34]
	_ = x[MessageIdentifierSetNoRotationAck-35]
	_ = x[MessageIdentifierRunSelfTest-36]
	_ = x[MessageIdentifierSelfTestResults-37]
	_ = x[MessageIdentifierPrepareData-50]
	_ = x[MessageIdentifierReqData-52]
	_ = x[MessageIdentifierReqDataAck-53]
	_ = x[MessageIdentifierMTData2-54]
	_ = x[MessageIdentifierReset-64]
	_ = x[MessageIdentifierResetAck-65]
	_ = x[MessageIdentifierError-66]
	_ = x[MessageIdentifierMasterIndication-70]
	_ = x[MessageIdentifierStopRecordingIndication-18]
	_ = x[MessageIdentifierFlushingIndication-19]
	_ = x[MessageIdentifierReqFilterSettings-160]
	_ = x[MessageIdentifierReqFilterSettingsAck-161]
	_ = x[MessageIdentifierSetFilterSettings-160]
	_ = x[MessageIdentifierSetFilterSettingsAck-161]
	_ = x[MessageIdentifierReqAmd-162]
	_ = x[MessageIdentifierReqAmdAck-163]
	_ = x[MessageIdentifierSetAmd-162]
	_ = x[MessageIdentifierSetAmdAck-163]
	_ = x[MessageIdentifierResetOrientation-164]
	_ = x[MessageIdentifierResetOrientationAck-165]
	_ = x[MessageIdentifierReqGpsStatus-166]
	_ = x[MessageIdentifierGpsStatus-167]
	_ = x[MessageIdentifierWriteDeviceID-176]
	_ = x[MessageIdentifierWriteDeviceIDAck-177]
	_ = x[MessageIdentifierWriteSecurityKey-178]
	_ = x[MessageIdentifierWriteSecurityKeyAck-179]
	_ = x[MessageIdentifierProtectFlash-180]
	_ = x[MessageIdentifierProtectFlashAck-181]
	_ = x[MessageIdentifierReqSecurityCheck-182]
	_ = x[MessageIdentifierSecurityCheck-183]
	_ = x[MessageIdentifierScanChannels-176]
	_ = x[MessageIdentifierScanChannelsAck-177]
	_ = x[MessageIdentifierEnableMaster-178]
	_ = x[MessageIdentifierEnableMasterAck-179]
	_ = x[MessageIdentifierDisableMaster-180]
	_ = x[MessageIdentifierDisableMasterAck-181]
	_ = x[MessageIdentifierReqRadioChannel-182]
	_ = x[MessageIdentifierReqRadioChannelAck-183]
	_ = x[MessageIdentifierSetClientPriority-184]
	_ = x[MessageIdentifierSetClientPriorityAck-185]
	_ = x[MessageIdentifierReqClientPriority-184]
	_ = x[MessageIdentifierReqClientPriorityAck-185]
	_ = x[MessageIdentifierSetWirelessConfig-186]
	_ = x[MessageIdentifierSetWirelessConfigAck-187]
	_ = x[MessageIdentifierReqWirelessConfig-186]
	_ = x[MessageIdentifierReqWirelessConfigAck-187]
	_ = x[MessageIdentifierUpdateBias-188]
	_ = x[MessageIdentifierUpdateBiasAck-189]
	_ = x[MessageIdentifierToggleIoPins-190]
	_ = x[MessageIdentifierToggleIoPinsAck-191]
	_ = x[MessageIdentifierSetTransportMode-194]
	_ = x[MessageIdentifierSetTransportModeAck-195]
	_ = x[MessageIdentifierReqTransportMode-194]
	_ = x[MessageIdentifierReqTransportModeAck-195]
	_ = x[MessageIdentifierAcceptMtw-196]
	_ = x[MessageIdentifierAcceptMtwAck-197]
	_ = x[MessageIdentifierRejectMtw-198]
	_ = x[MessageIdentifierRejectMtwAck-199]
	_ = x[MessageIdentifierInfoRequest-200]
	_ = x[MessageIdentifierInfoRequestAck-201]
	_ = x[MessageIdentifierReqFrameRates-202]
	_ = x[MessageIdentifierReqFrameRatesAck-203]
	_ = x[MessageIdentifierStartRecording-204]
	_ = x[MessageIdentifierStartRecordingAck-205]
	_ = x[MessageIdentifierStopRecording-206]
	_ = x[MessageIdentifierStopRecordingAck-207]
	_ = x[MessageIdentifierInfoBatteryLevel-73]
	_ = x[MessageIdentifierInfoTemperature-74]
	_ = x[MessageIdentifierGotoOperational-192]
	_ = x[MessageIdentifierGotoOperationalAck-193]
	_ = x[MessageIdentifierReqEmts-144]
	_ = x[MessageIdentifierEmtsData-145]
	_ = x[MessageIdentifierRestoreEmts-148]
	_ = x[MessageIdentifierRestoreEmtsAck-149]
	_ = x[MessageIdentifierStoreEmts-150]
	_ = x[MessageIdentifierStoreEmtsAck-151]
	_ = x[MessageIdentifierGotoTransparentMode-80]
	_ = x[MessageIdentifierGotoTransparentModeAck-81]
}

const _MessageIdentifier_name = "ReqDIDDeviceIDInitBusInitBusResultsReqPeriodReqPeriodAckSetBidSetBidAckBusPowerBusPowerAckReqDataLengthDataLengthReqConfigurationConfigurationRestoreFactoryDefRestoreFactoryDefAckGotoMeasurementGotoMeasurementAckReqFirmwareRevisionFirmwareRevisionReqBluetoothDisableReqBluetoothDisableAckReqXmOutputModeReqXmOutputModeAckReqBaudrateReqBaudrateAckReqSyncModeReqSyncModeAckReqProductCodeProductCodeReqProcessingFlagsReqProcessingFlagsAckSetNoRotationSetNoRotationAckRunSelfTestSelfTestResultsReqInputTriggerReqInputTriggerAckReqOutputTriggerReqOutputTriggerAckSetSyncBoxModeSetSyncBoxModeAckSetSyncConfigurationSetSyncConfigurationAckDriverDisconnectDriverDisconnectAckGotoConfigGotoConfigAckBusDataReqDataReqDataAckMTData2WakeupWakeupAckResetResetAckErrorXmPowerOffMasterIndicationInfoBatteryLevelInfoTemperatureGotoTransparentModeGotoTransparentModeAckSetUtcTimeSetUtcTimeAckReqAvailableFilterProfilesAvailableFilterProfilesReqFilterProfileReqFilterProfileAckReqGravityMagnitudeReqGravityMagnitudeAckReqGpsLeverArmReqGpsLeverArmAckReqMagneticFieldReqMagneticFieldAckReqLatLonAltReqLatLonAltAckReqXmErrorModeReqXmErrorModeAckReqBufferSizeReqBufferSizeAckReqExtOutputModeReqExtOutputModeAckReqBatteryLevelBatterylevelReqMasterSettingsMasterSettingsReqEmtsEmtsDataRestoreEmtsRestoreEmtsAckStoreEmtsStoreEmtsAckReqActiveClockCorrectionActiveClockCorrectionStoreActiveClockCorrectionStoreActiveClockCorrectionAckReqFilterSettingsReqFilterSettingsAckReqAmdReqAmdAckResetOrientationResetOrientationAckReqGpsStatusGpsStatusAdjustUtcTimeAdjustUtcTimeAckWriteDeviceIDWriteDeviceIDAckWriteSecurityKeyWriteSecurityKeyAckProtectFlashProtectFlashAckReqSecurityCheckSecurityCheckSetClientPrioritySetClientPriorityAckSetWirelessConfigSetWirelessConfigAckUpdateBiasUpdateBiasAckToggleIoPinsToggleIoPinsAckReqOutputConfigurationReqOutputConfigurationAckSetTransportModeSetTransportModeAckAcceptMtwAcceptMtwAckRejectMtwRejectMtwAckInfoRequestInfoRequestAckReqFrameRatesReqFrameRatesAckStartRecordingStartRecordingAckStopRecordingStopRecordingAckReqOutputModeReqOutputModeAckReqOutputSettingsReqOutputSettingsAckReqOutputSkipFactorReqOutputSkipFactorAckReqSyncInSettingsReqSyncInSettingsAckReqSyncOutSettingsReqSyncOutSettingsAckReqErrorModeReqErrorModeAckReqTransmitDelayReqTransmitDelayAckSetMfmResultsSetMfmResultsAckReqObjectAlignmentReqObjectAlignmentAck"

var _MessageIdentifier_map = map[MessageIdentifier]string{
	0:   _MessageIdentifier_name[0:6],
	1:   _MessageIdentifier_name[6:14],
	2:   _MessageIdentifier_name[14:21],
	3:   _MessageIdentifier_name[21:35],
	4:   _MessageIdentifier_name[35:44],
	5:   _MessageIdentifier_name[44:56],
	6:   _MessageIdentifier_name[56:62],
	7:   _MessageIdentifier_name[62:71],
	8:   _MessageIdentifier_name[71:79],
	9:   _MessageIdentifier_name[79:90],
	10:  _MessageIdentifier_name[90:103],
	11:  _MessageIdentifier_name[103:113],
	12:  _MessageIdentifier_name[113:129],
	13:  _MessageIdentifier_name[129:142],
	14:  _MessageIdentifier_name[142:159],
	15:  _MessageIdentifier_name[159:179],
	16:  _MessageIdentifier_name[179:194],
	17:  _MessageIdentifier_name[194:212],
	18:  _MessageIdentifier_name[212:231],
	19:  _MessageIdentifier_name[231:247],
	20:  _MessageIdentifier_name[247:266],
	21:  _MessageIdentifier_name[266:288],
	22:  _MessageIdentifier_name[288:303],
	23:  _MessageIdentifier_name[303:321],
	24:  _MessageIdentifier_name[321:332],
	25:  _MessageIdentifier_name[332:346],
	26:  _MessageIdentifier_name[346:357],
	27:  _MessageIdentifier_name[357:371],
	28:  _MessageIdentifier_name[371:385],
	29:  _MessageIdentifier_name[385:396],
	32:  _MessageIdentifier_name[396:414],
	33:  _MessageIdentifier_name[414:435],
	34:  _MessageIdentifier_name[435:448],
	35:  _MessageIdentifier_name[448:464],
	36:  _MessageIdentifier_name[464:475],
	37:  _MessageIdentifier_name[475:490],
	38:  _MessageIdentifier_name[490:505],
	39:  _MessageIdentifier_name[505:523],
	40:  _MessageIdentifier_name[523:539],
	41:  _MessageIdentifier_name[539:558],
	42:  _MessageIdentifier_name[558:572],
	43:  _MessageIdentifier_name[572:589],
	44:  _MessageIdentifier_name[589:609],
	45:  _MessageIdentifier_name[609:632],
	46:  _MessageIdentifier_name[632:648],
	47:  _MessageIdentifier_name[648:667],
	48:  _MessageIdentifier_name[667:677],
	49:  _MessageIdentifier_name[677:690],
	50:  _MessageIdentifier_name[690:697],
	52:  _MessageIdentifier_name[697:704],
	53:  _MessageIdentifier_name[704:714],
	54:  _MessageIdentifier_name[714:721],
	62:  _MessageIdentifier_name[721:727],
	63:  _MessageIdentifier_name[727:736],
	64:  _MessageIdentifier_name[736:741],
	65:  _MessageIdentifier_name[741:749],
	66:  _MessageIdentifier_name[749:754],
	68:  _MessageIdentifier_name[754:764],
	70:  _MessageIdentifier_name[764:780],
	73:  _MessageIdentifier_name[780:796],
	74:  _MessageIdentifier_name[796:811],
	80:  _MessageIdentifier_name[811:830],
	81:  _MessageIdentifier_name[830:852],
	96:  _MessageIdentifier_name[852:862],
	97:  _MessageIdentifier_name[862:875],
	98:  _MessageIdentifier_name[875:901],
	99:  _MessageIdentifier_name[901:924],
	100: _MessageIdentifier_name[924:940],
	101: _MessageIdentifier_name[940:959],
	102: _MessageIdentifier_name[959:978],
	103: _MessageIdentifier_name[978:1000],
	104: _MessageIdentifier_name[1000:1014],
	105: _MessageIdentifier_name[1014:1031],
	106: _MessageIdentifier_name[1031:1047],
	107: _MessageIdentifier_name[1047:1066],
	110: _MessageIdentifier_name[1066:1078],
	111: _MessageIdentifier_name[1078:1093],
	130: _MessageIdentifier_name[1093:1107],
	131: _MessageIdentifier_name[1107:1124],
	132: _MessageIdentifier_name[1124:1137],
	133: _MessageIdentifier_name[1137:1153],
	134: _MessageIdentifier_name[1153:1169],
	135: _MessageIdentifier_name[1169:1188],
	136: _MessageIdentifier_name[1188:1203],
	137: _MessageIdentifier_name[1203:1215],
	138: _MessageIdentifier_name[1215:1232],
	139: _MessageIdentifier_name[1232:1246],
	144: _MessageIdentifier_name[1246:1253],
	145: _MessageIdentifier_name[1253:1261],
	148: _MessageIdentifier_name[1261:1272],
	149: _MessageIdentifier_name[1272:1286],
	150: _MessageIdentifier_name[1286:1295],
	151: _MessageIdentifier_name[1295:1307],
	156: _MessageIdentifier_name[1307:1331],
	157: _MessageIdentifier_name[1331:1352],
	158: _MessageIdentifier_name[1352:1378],
	159: _MessageIdentifier_name[1378:1407],
	160: _MessageIdentifier_name[1407:1424],
	161: _MessageIdentifier_name[1424:1444],
	162: _MessageIdentifier_name[1444:1450],
	163: _MessageIdentifier_name[1450:1459],
	164: _MessageIdentifier_name[1459:1475],
	165: _MessageIdentifier_name[1475:1494],
	166: _MessageIdentifier_name[1494:1506],
	167: _MessageIdentifier_name[1506:1515],
	168: _MessageIdentifier_name[1515:1528],
	169: _MessageIdentifier_name[1528:1544],
	176: _MessageIdentifier_name[1544:1557],
	177: _MessageIdentifier_name[1557:1573],
	178: _MessageIdentifier_name[1573:1589],
	179: _MessageIdentifier_name[1589:1608],
	180: _MessageIdentifier_name[1608:1620],
	181: _MessageIdentifier_name[1620:1635],
	182: _MessageIdentifier_name[1635:1651],
	183: _MessageIdentifier_name[1651:1664],
	184: _MessageIdentifier_name[1664:1681],
	185: _MessageIdentifier_name[1681:1701],
	186: _MessageIdentifier_name[1701:1718],
	187: _MessageIdentifier_name[1718:1738],
	188: _MessageIdentifier_name[1738:1748],
	189: _MessageIdentifier_name[1748:1761],
	190: _MessageIdentifier_name[1761:1773],
	191: _MessageIdentifier_name[1773:1788],
	192: _MessageIdentifier_name[1788:1810],
	193: _MessageIdentifier_name[1810:1835],
	194: _MessageIdentifier_name[1835:1851],
	195: _MessageIdentifier_name[1851:1870],
	196: _MessageIdentifier_name[1870:1879],
	197: _MessageIdentifier_name[1879:1891],
	198: _MessageIdentifier_name[1891:1900],
	199: _MessageIdentifier_name[1900:1912],
	200: _MessageIdentifier_name[1912:1923],
	201: _MessageIdentifier_name[1923:1937],
	202: _MessageIdentifier_name[1937:1950],
	203: _MessageIdentifier_name[1950:1966],
	204: _MessageIdentifier_name[1966:1980],
	205: _MessageIdentifier_name[1980:1997],
	206: _MessageIdentifier_name[1997:2010],
	207: _MessageIdentifier_name[2010:2026],
	208: _MessageIdentifier_name[2026:2039],
	209: _MessageIdentifier_name[2039:2055],
	210: _MessageIdentifier_name[2055:2072],
	211: _MessageIdentifier_name[2072:2092],
	212: _MessageIdentifier_name[2092:2111],
	213: _MessageIdentifier_name[2111:2133],
	214: _MessageIdentifier_name[2133:2150],
	215: _MessageIdentifier_name[2150:2170],
	216: _MessageIdentifier_name[2170:2188],
	217: _MessageIdentifier_name[2188:2209],
	218: _MessageIdentifier_name[2209:2221],
	219: _MessageIdentifier_name[2221:2236],
	220: _MessageIdentifier_name[2236:2252],
	221: _MessageIdentifier_name[2252:2271],
	222: _MessageIdentifier_name[2271:2284],
	223: _MessageIdentifier_name[2284:2300],
	224: _MessageIdentifier_name[2300:2318],
	225: _MessageIdentifier_name[2318:2339],
}

func (i MessageIdentifier) String() string {
	if str, ok := _MessageIdentifier_map[i]; ok {
		return str
	}
	return "MessageIdentifier(" + strconv.FormatInt(int64(i), 10) + ")"
}
