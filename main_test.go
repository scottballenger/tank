package main

import (
	"io"
	"strings"
	"testing"
)

func Test_xRange(t *testing.T) {
	type args struct {
		angle float64
		v     float64
	}
	tests := []struct {
		name  string
		args  args
		wantX float64
		wantT float64
	}{
		{
			name:  "Test 1 - 300m/s",
			args:  args{22.5, 300.0},
			wantX: 6489.43424174303,
			wantT: 23.41371002524347,
		},
		{
			name:  "Test 2 - 600 m/s",
			args:  args{22.5, 600.0},
			wantX: 25957.73696697212,
			wantT: 46.82742005048694,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotT := xRange(tt.args.angle, tt.args.v)
			if gotX != tt.wantX {
				t.Errorf("xRange() gotX = %v, want %v", gotX, tt.wantX)
			}
			if gotT != tt.wantT {
				t.Errorf("xRange() gotT = %v, want %v", gotT, tt.wantT)
			}
		})
	}
}

func Test_getRandomValue(t *testing.T) {
	type args struct {
		min float64
		max float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "MinMax",
			args: args{0.0, 1.0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRandomValue(tt.args.min, tt.args.max); got <= tt.args.min || got >= tt.args.max {
				t.Errorf("getRandomValue() = %v, want >= %v && want <= %v", got, tt.args.min, tt.args.max)
			}
		})
	}
}

func Test_getMilesOrKilometers(t *testing.T) {
	type args struct {
		value        float64
		englishUnits bool
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Metric",
			args: args{1000.0, false},
			want: 1000.0 / metersPerKilometer,
		},
		{
			name: "English",
			args: args{1000.0, true},
			want: 1000.0 * feetPerMeter / feetPerMile,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMilesOrKilometers(tt.args.value, tt.args.englishUnits); got != tt.want {
				t.Errorf("getMilesOrKilometers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFeetOrMeters(t *testing.T) {
	type args struct {
		value        float64
		englishUnits bool
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Metric",
			args: args{1000.0, false},
			want: 1000.0,
		},
		{
			name: "English",
			args: args{1000.0, true},
			want: 1000.0 * feetPerMeter,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFeetOrMeters(tt.args.value, tt.args.englishUnits); got != tt.want {
				t.Errorf("getFeetOrMeters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getGameModes(t *testing.T) {
	type args struct {
		shootModeAuto  bool
		targetModeAuto bool
	}
	tests := []struct {
		name               string
		args               args
		wantTargetModeAuto bool
		wantGameSpeed      int
	}{
		{
			name:               "Default",
			args:               args{false, false},
			wantTargetModeAuto: false,
			wantGameSpeed:      1,
		},
		{
			name:               "Auto Shoot",
			args:               args{true, false},
			wantTargetModeAuto: true,
			wantGameSpeed:      10,
		},
		{
			name:               "Manual Shoot - Pause Target",
			args:               args{false, false},
			wantTargetModeAuto: false,
			wantGameSpeed:      1,
		},
		{
			name:               "Manual Shoot - Realtime Target",
			args:               args{false, true},
			wantTargetModeAuto: true,
			wantGameSpeed:      1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTargetModeAuto, gotGameSpeed := getTargetMode(tt.args.shootModeAuto, tt.args.targetModeAuto)
			if gotTargetModeAuto != tt.wantTargetModeAuto {
				t.Errorf("getGameModes() gotTargetModeAuto = %v, want %v", gotTargetModeAuto, tt.wantTargetModeAuto)
			}
			if gotGameSpeed != tt.wantGameSpeed {
				t.Errorf("getGameModes() gotGameSpeed = %v, want %v", gotGameSpeed, tt.wantGameSpeed)
			}
		})
	}
}

func Test_getImpactTimelineIndices(t *testing.T) {
	type args struct {
		shotDistance   float64
		targetDistance float64
		maxDistance    float64
	}
	tests := []struct {
		name            string
		args            args
		wantShotIndex   int
		wantTargetIndex int
	}{
		{
			name:            "Test undershot and ensure that 'T' is in correct location",
			args:            args{1.0, 20.0, 100.0},
			wantShotIndex:   4,
			wantTargetIndex: 10,
		},
		{
			name:            "Test short undershot and ensure that 'T' is in correct location",
			args:            args{1.0, 2.0, 1000.0},
			wantShotIndex:   4,
			wantTargetIndex: 5,
		},
		{
			name:            "Test short overshot and ensure that 'T' is in correct location",
			args:            args{2.0, 1.0, 1000.0},
			wantShotIndex:   4,
			wantTargetIndex: 2,
		},
		{
			name:            "Test overshot of short target and ensure that 'T' is in correct location",
			args:            args{30.0, 1.0, 100.0},
			wantShotIndex:   15,
			wantTargetIndex: 2,
		},
		{
			name:            "Test overshot at same location",
			args:            args{10.0, 10.1, 100.0},
			wantShotIndex:   4,
			wantTargetIndex: 5,
		},
		{
			name:            "Test undershot at same location",
			args:            args{10.1, 10.0, 100.0},
			wantShotIndex:   6,
			wantTargetIndex: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotShotIndex, gotTargetIndex := getImpactTimelineIndices(tt.args.shotDistance, tt.args.targetDistance, tt.args.maxDistance)
			if gotShotIndex != tt.wantShotIndex {
				t.Errorf("getImpactTimelineIndices() gotShotIndex = %v, want %v", gotShotIndex, tt.wantShotIndex)
			}
			if gotTargetIndex != tt.wantTargetIndex {
				t.Errorf("getImpactTimelineIndices() gotTargetIndex = %v, want %v", gotTargetIndex, tt.wantTargetIndex)
			}
		})
	}
}

func Test_isGameOverMan(t *testing.T) {
	type args struct {
		targetRange float64
		deathRadius float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Not Yet",
			args: args{100.0, 20.0},
			want: false,
		},
		{
			name: "Game Over Man",
			args: args{10.0, 20.0},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isGameOverMan(tt.args.targetRange, tt.args.deathRadius); got != tt.want {
				t.Errorf("isGameOverMan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_printImpactResults(t *testing.T) {
	type args struct {
		shotRange   float64
		targetRange float64
		shotDelta   float64
		deathRadius float64
		shotCount   int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Direct Hit",
			args: args{100.0, 100.0, 10.0, 20.0, 30},
			want: true,
		},
		{
			name: "Game Over Man",
			args: args{100.0, 10.0, 30.0, 20.0, 30},
			want: true,
		},
		{
			name: "Undershot",
			args: args{100.0, 100.0, 30.0, 20.0, 30},
			want: false,
		},
		{
			name: "Overshot",
			args: args{100.0, 100.0, -30.0, 20.0, 30},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := printImpactResults(tt.args.shotRange, tt.args.targetRange, tt.args.shotDelta, tt.args.deathRadius, tt.args.shotCount); got != tt.want {
				t.Errorf("printImpactResults() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_takeShot(t *testing.T) {
	type args struct {
		targetRange    float64
		targetVmps     float64
		shotCount      int
		shotAngle      float64
		projectileVmps float64
	}
	tests := []struct {
		name          string
		args          args
		wantShotRange float64
		wantShotTime  float64
		wantShotDelta float64
	}{
		{
			name:          "Test 1",
			args:          args{10000.0, 0.0, 1, 22.5, 300.0},
			wantShotRange: 6489.43424174303,
			wantShotTime:  23.41371002524347,
			wantShotDelta: 3510.56575825697,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetRange = tt.args.targetRange
			targetVmps = tt.args.targetVmps
			gotShotRange, gotShotTime, gotShotDelta := takeShot(tt.args.shotCount, tt.args.shotAngle, tt.args.projectileVmps)
			if gotShotRange != tt.wantShotRange {
				t.Errorf("takeShot() gotShotRange = %v, want %v", gotShotRange, tt.wantShotRange)
			}
			if gotShotTime != tt.wantShotTime {
				t.Errorf("takeShot() gotShotTime = %v, want %v", gotShotTime, tt.wantShotTime)
			}
			if gotShotDelta != tt.wantShotDelta {
				t.Errorf("takeShot() gotShotDelta = %v, want %v", gotShotDelta, tt.wantShotDelta)
			}
		})
	}
}

func Test_predictNextShotAngle(t *testing.T) {
	type args struct {
		targetVmps     float64
		projectileVmps float64
		shotRange      float64
		shotTime       float64
		shotDelta      float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Test 1",
			args: args{0.0, 600.0, 2689.835682980686, 23.429597899903456, 7310.164317019314},
			want: 7.903772202260505,
		},
		{
			name: "Out of Range",
			args: args{0.0, 300.0, 2689.835682980686, 23.429597899903456, 7310.164317019314},
			want: maxShotAngle,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetVmps = tt.args.targetVmps
			projectileVmps = tt.args.projectileVmps
			if got := predictNextShotAngle(tt.args.shotRange, tt.args.shotTime, tt.args.shotDelta); got != tt.want {
				t.Errorf("predictNextShotAngle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNextShotAngle(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Return 0",
			args: args{strings.NewReader("0\n")},
			want: 0.0,
		},
		{
			name: "Return 1.359",
			args: args{strings.NewReader("1.359\n")},
			want: 1.359,
		},
		{
			name: "Return 10.0 after errors",
			args: args{strings.NewReader("10.a\n10\n")},
			want: 10.0,
		},
		{
			name: "Return 45.0 after 46 out of range",
			args: args{strings.NewReader("46\n45\n")},
			want: 45.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNextShotAngle(tt.args.reader); got != tt.want {
				t.Errorf("getNextShotAngle() = %v, want %v", got, tt.want)
			}
		})
	}
}
