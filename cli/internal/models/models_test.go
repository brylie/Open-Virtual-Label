package models

import "testing"

func TestHasMasteringData_NilMastering(t *testing.T) {
	var track Track
	if track.HasMasteringData() {
		t.Error("expected false for nil mastering")
	}
}

func TestHasMasteringData_MissingFields(t *testing.T) {
	track := Track{Mastering: &TrackMastering{}}
	if track.HasMasteringData() {
		t.Error("expected false when LUFS and peak are nil")
	}
}

func TestHasMasteringData_PartialFields(t *testing.T) {
	lufs := -14.5
	track := Track{Mastering: &TrackMastering{IntegratedLUFS: &lufs}}
	if track.HasMasteringData() {
		t.Error("expected false when only LUFS is set (peak missing)")
	}
}

func TestHasMasteringData_AllFields(t *testing.T) {
	lufs := -14.5
	peak := -1.0
	track := Track{
		Mastering: &TrackMastering{
			IntegratedLUFS: &lufs,
			TruePeakDBTP:   &peak,
		},
	}
	if !track.HasMasteringData() {
		t.Error("expected true when both LUFS and peak are set")
	}
}
