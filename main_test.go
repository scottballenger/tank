package main

import "testing"

func TestImpactTimeline(t *testing.T) {
	// Test hit and ensure that '*' is in correct location
	impactTimeline := getImpactTimeline(1.0, 1.0, 100.0, true)
	if impactTimeline[0][:3] != " /\\" {
		t.Errorf("Hit Flight '%s' != ' /\\'", impactTimeline[0][:3])
	}
	if impactTimeline[1][:4] != "/--*" {
		t.Errorf("Hit Impact '%s' != '/--*'", impactTimeline[1][:4])
	}

	// Test undershot and ensure that 'T' is in correct location
	impactTimeline = getImpactTimeline(1.0, 20.0, 100.0, false)
	if impactTimeline[0][:3] != " /\\" {
		t.Errorf("Undershot Flight '%s' != ' /\\'", impactTimeline[0][:3])
	}
	if impactTimeline[1][:11] != "/--\\-----+T" {
		t.Errorf("Undershot Impact '%s' != '/--\\-----+T'", impactTimeline[1])
	}

	// Test short undershot and ensure that 'T' is in correct location
	impactTimeline = getImpactTimeline(1.0, 2.0, 100.0, false)
	if impactTimeline[0][:3] != " /\\" {
		t.Errorf("Short Undershot Flight '%s' != ' /\\'", impactTimeline[0][:3])
	}
	if impactTimeline[1][:5] != "/--\\T" {
		t.Errorf("Short Undershot Impact '%s' != '/--\\T'", impactTimeline[1])
	}

	// Test short overshot and ensure that 'T' is in correct location
	impactTimeline = getImpactTimeline(2.0, 1.0, 100.0, false)
	if impactTimeline[0][:3] != " /\\" {
		t.Errorf("Short Overshot Flight '%s' != ' /\\'", impactTimeline[0][:3])
	}
	if impactTimeline[1][:4] != "/T-\\" {
		t.Errorf("Short Overshot Impact '%s' != '/T--\\'", impactTimeline[1])
	}

	// Test overshot of long target and ensure that 'T' is in correct location
	impactTimeline = getImpactTimeline(30.0, 20.0, 100.0, false)
	if impactTimeline[1][:11] != "/--------+T" {
		t.Errorf("Overshot Long Target Impact '%s' != '/--------+T'", impactTimeline[1])
	}

	// Test overshot of short target and ensure that 'T' is in correct location
	impactTimeline = getImpactTimeline(30.0, 1.0, 100.0, false)
	if impactTimeline[1][:2] != "/T" {
		t.Errorf("Overshot Short Target Impact '%s' != '/T'", impactTimeline[1])
	}
}
