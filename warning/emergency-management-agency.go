package warning

import (
	"errors"
	"math/rand"
	"time"

	"github.com/HydrologicEngineeringCenter/go-statistics/paireddata"
)

type EmergencyManagementAgency struct {
	WarningIssuance            paireddata.UncertaintyPairedData
	FirstAlertDay              paireddata.UncertaintyPairedData
	FirstAlertNight            paireddata.UncertaintyPairedData
	ProtectiveActionInitiation paireddata.UncertaintyPairedData
}

type EmergencyManagementAgencySample struct {
	WarningIssuance            paireddata.PairedData
	FirstAlertDay              paireddata.PairedData
	FirstAlertNight            paireddata.PairedData
	ProtectiveActionInitiation paireddata.PairedData
}
type WarningScenario struct {
	WarningIssuanceTime        time.Time
	FirstAlert                 paireddata.PairedData
	ProtectiveActionInitiation paireddata.PairedData
}

func (ema EmergencyManagementAgency) Sample(seed int64) (EmergencyManagementAgencySample, error) {
	r := rand.New(rand.NewSource(seed))
	wi := ema.WarningIssuance.SampleValueSampler(r.Float64())
	wipd, wiok := wi.(paireddata.PairedData)
	if !wiok {
		return EmergencyManagementAgencySample{}, errors.New("could not sample warning issuance")
	}
	fad := ema.FirstAlertDay.SampleValueSampler(r.Float64())
	fadpd, fadok := fad.(paireddata.PairedData)
	if !fadok {
		return EmergencyManagementAgencySample{}, errors.New("could not sample first alert day")
	}
	fan := ema.FirstAlertNight.SampleValueSampler(r.Float64())
	fanpd, fanok := fan.(paireddata.PairedData)
	if !fanok {
		return EmergencyManagementAgencySample{}, errors.New("could not sample first alert night")
	}
	pai := ema.ProtectiveActionInitiation.SampleValueSampler(r.Float64())
	paipd, paiok := pai.(paireddata.PairedData)
	if !paiok {
		return EmergencyManagementAgencySample{}, errors.New("could not sample protective action initiation")
	}
	return EmergencyManagementAgencySample{WarningIssuance: wipd, FirstAlertDay: fadpd, FirstAlertNight: fanpd, ProtectiveActionInitiation: paipd}, nil
}
func (emas EmergencyManagementAgencySample) IssueWarning(seed int64, startTime time.Time) (WarningScenario, error) {
	r := rand.New(rand.NewSource(seed))
	minsFromStart := emas.WarningIssuance.SampleValue(r.Float64())
	warnTime := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), startTime.Hour(), startTime.Minute()+int(minsFromStart), 0, 0, time.UTC)
	fac := interpolateCurves(emas.FirstAlertDay, emas.FirstAlertNight, warnTime)
	return WarningScenario{WarningIssuanceTime: warnTime, FirstAlert: fac, ProtectiveActionInitiation: emas.ProtectiveActionInitiation}, nil
}
