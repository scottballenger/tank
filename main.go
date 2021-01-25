package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// -----------------------------------------------------------
// NOTE:
//    ALL distance units are METERS except where specified!
//    ALL angle units are DEGREES except where specified!
// -----------------------------------------------------------

const (
	minProjectileVmps  = 300  // meters/sec
	maxProjectileVmps  = 600  // meters/sec
	minTargetVkph      = 0    // kilometers/hour
	maxTargetVkph      = 60   // kilometers/hour
	impactRadius       = 20.0 // meters
	minShotAngle       = 1.0  // degrees
	maxShotAngle       = 45.0 // degrees
	metersPerKilometer = 1000.0
	feetPerMile        = 5280.0
	feetPerMeter       = 3.28084
	secondsPerHour     = 60.0 * 60.0
	gravityAmps2       = 9.80665

	flightPath = " /~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
	impactPath = "/--------+---------+---------+---------+---------|"
	// Ruler      01234567890123456789012345678901234567890123456789
	gameOverMan = "GAME OVER MAN, you just got crushed by the other tank!"
)

var (
	projectileVmps        float64
	targetVkph            float64
	targetVmps            float64
	maxRange              float64
	targetRange           float64
	deathRadius           float64
	wg                    sync.WaitGroup
	shootModeAuto         bool = false // true = Auto Shoot Mode, false = Manual Shot
	targetModeAuto        bool = false // true = target moves when deciding shot, false = target pauses when deciding shot
	targetSpeedMultiplier int          // times faster than real-time
	englishUnits          bool = false // true = english units, false = metric units
	printShotProfile      bool = false // true = print the shot profile on startup, false = don't
	rulerText             string

	englishOrMetric   = map[bool]string{true: "English", false: "Metric"}
	feetOrMeters      = map[bool]string{true: "feet", false: "meters"}
	milesOrKilometers = map[bool]string{true: "miles", false: "kilometers"}
)

func parseFlags() {
	flag.BoolVar(&shootModeAuto, "a", shootModeAuto, "Auto Shot Mode (default - Manual Shot)")
	flag.BoolVar(&targetModeAuto, "m", targetModeAuto, "Real-time Target Movement (default - Pause Target During Shot Decision)")
	flag.BoolVar(&englishUnits, "e", englishUnits, "English Units (default - Metric)")
	flag.BoolVar(&printShotProfile, "p", printShotProfile, "Print Shot Profile")
	flag.Float64Var(&deathRadius, "d", impactRadius, "Detonation Radius (meters)")
	flag.Parse()
}

func getTargetMode(shootModeAuto, targetModeAutoDefault bool) (targetModeAuto bool, targetSpeedMultiplier int) {
	targetModeAuto = targetModeAutoDefault
	if shootModeAuto {
		fmt.Println("Shot Mode: Auto")
		targetSpeedMultiplier = 10
		targetModeAuto = true
	} else {
		fmt.Println("Shot Mode: Manual")
		targetSpeedMultiplier = 1
		if targetModeAuto {
			fmt.Println("Target Mode: Realtime Target Movement")
		} else {
			fmt.Println("Target Mode: Pause Target During Shot Decision")
		}
	}
	return
}

func getRandomValue(min, max float64) float64 {
	return min + float64(rand.Intn(10000))*float64(max-min)/10000.0
}

func displayShotProfile() {
	fmt.Println("")
	fmt.Println("Shot Profile:")
	fmt.Printf("+-------+------------+-------+\n")
	fmt.Printf("| Angle | Shot Range | Time  |\n")
	fmt.Printf("| (deg) | %10s | (sec) |\n", "("+feetOrMeters[englishUnits]+")")
	fmt.Printf("+-------+------------+-------+\n")
	for angle := minShotAngle; angle <= maxShotAngle; angle += 1.0 {
		shotRange, shotTime := xRange(angle, projectileVmps)
		fmt.Printf("| %5.1f | %10.1f | %5.1f |\n", angle, getFeetOrMeters(shotRange, englishUnits), shotTime)
	}
	fmt.Printf("+-------+------------+-------+\n")
	fmt.Println("")
}

func initialize() {
	parseFlags()
	targetModeAuto, targetSpeedMultiplier = getTargetMode(shootModeAuto, targetModeAuto)

	fmt.Printf("Units: %s\n", englishOrMetric[englishUnits])
	fmt.Printf("Detonation Radius = %s\n", getDisplayText(deathRadius))

	// Initialize random values.
	rand.Seed(time.Now().UnixNano())
	projectileVmps = getRandomValue(minProjectileVmps, maxProjectileVmps)
	targetVkph = getRandomValue(minTargetVkph, maxTargetVkph)
	targetVmps = targetVkph * (metersPerKilometer / secondsPerHour)
	maxRange, _ = xRange(maxShotAngle, projectileVmps)
	targetRange = getRandomValue(maxRange*0.2, maxRange)

	// Calculate distance markers for the legend undel the timeline.
	rulerText = fmt.Sprintf("       %5s     %5s     %5s     %5s     %5s",
		getRulerText(1.0*maxRange/5.0),
		getRulerText(2.0*maxRange/5.0),
		getRulerText(3.0*maxRange/5.0),
		getRulerText(4.0*maxRange/5.0),
		getRulerText(5.0*maxRange/5.0),
	)
	rulerText = rulerText[:len(rulerText)-1] + strings.Title(milesOrKilometers[englishUnits])

	if printShotProfile {
		displayShotProfile()
	}

}

func xRange(angle, v float64) (x, t float64) {
	radians := (2.0 * math.Pi) * (angle / 360.0)
	t = (math.Sin(radians) * v * 2.0) / gravityAmps2
	x = math.Cos(radians) * t * v
	return
}

func xAngle(x, v float64) float64 {
	radians := math.Asin((x*gravityAmps2)/(v*v)) / 2.0
	angle := (radians * 360.0) / (2.0 * math.Pi)
	return angle
}

func getMilesOrKilometers(value float64, englishUnits bool) float64 {
	if englishUnits {
		return ((value * feetPerMeter) / feetPerMile)
	}
	return value / metersPerKilometer
}

func getRulerText(value float64) string {
	return fmt.Sprintf("%4.1f%1s", getMilesOrKilometers(value, englishUnits), strings.Title(milesOrKilometers[englishUnits][:1]))
}

func getFeetOrMeters(value float64, englishUnits bool) float64 {
	if englishUnits {
		return value * feetPerMeter
	}
	return value
}

func getDisplayText(value float64) string {
	return fmt.Sprintf("%3.1f %s", getFeetOrMeters(value, englishUnits), feetOrMeters[englishUnits])
}

func printHeader() {
	fmt.Println("==================================")
	fmt.Printf("Projectile Velocity  = %s/sec\n", getDisplayText(projectileVmps))
	fmt.Printf("Max Projectile Range = %s\n", getDisplayText(maxRange))
	fmt.Printf("Target Velocity      = %3.1f %s/hour\n", getMilesOrKilometers(targetVkph*metersPerKilometer, englishUnits), milesOrKilometers[englishUnits])
	fmt.Printf("Target Velocity      = %s/sec\n", getDisplayText(targetVmps))
	fmt.Printf("Current Target Range = %3.1f %s\n", getMilesOrKilometers(targetRange, englishUnits), milesOrKilometers[englishUnits])
	fmt.Printf("Current Target Range = %s\n", getDisplayText(targetRange))
	fmt.Println("----------------------------------")
}

func getImpactTimelineIndices(shotDistance, targetDistance, maxDistance float64) (shotIndex, targetIndex int) {
	maxString := len(impactPath) - 1
	targetIndex = int(targetDistance/maxDistance*float64(maxString)) + 1
	shotIndex = int(shotDistance/maxDistance*float64(maxString)) + 1
	if shotIndex == targetIndex {
		if shotDistance < targetDistance {
			shotIndex = targetIndex - 1
		} else {
			shotIndex = targetIndex + 1
		}
	}
	if shotIndex < 4 || targetIndex < 5 {
		if shotIndex < 4 && targetIndex < 5 {
			if targetIndex < shotIndex {
				shotIndex += 4 - shotIndex
				targetIndex = int(math.Max(2, float64(targetIndex)))
			} else {
				inc := int(math.Max(4-float64(shotIndex), 5-float64(targetIndex)))
				shotIndex += inc
				targetIndex += inc
			}
		} else if shotIndex < 4 {
			shotIndex += 4 - shotIndex
		} else { // targetIndex < 5
			targetIndex = int(math.Max(2, float64(targetIndex)))
		}
	}
	return
}

func printImpactTimeline(shotDistance float64, hit bool) {

	shotIndex, targetIndex := getImpactTimelineIndices(shotDistance, targetRange, maxRange)
	curFlightPath := flightPath
	curImpactPath := impactPath
	if hit {
		curFlightPath = curFlightPath[:targetIndex-2] + "\\"
		curImpactPath = curImpactPath[:targetIndex-1] + "*" + curImpactPath[targetIndex:]
	} else {
		curFlightPath = curFlightPath[:shotIndex-2] + "\\"
		curImpactPath = curImpactPath[:shotIndex-1] + "\\" + curImpactPath[shotIndex:]
		curImpactPath = curImpactPath[:targetIndex-1] + "T" + curImpactPath[targetIndex:]
	}
	fmt.Println("")
	fmt.Println(curFlightPath)
	fmt.Println(curImpactPath)
	fmt.Println(rulerText)
	fmt.Println("")
}

func printImpactResults(shotRange, targetRange, shotDelta, deathRadius float64, shotCount int) bool {
	fmt.Printf("Target Range = %s at time of impact.\n", getDisplayText(targetRange))
	if math.Abs(shotDelta) <= deathRadius {
		printImpactTimeline(shotRange, true)
		fmt.Println("")
		fmt.Printf("Direct hit (within %s) after %d shots!!\n", getDisplayText(math.Abs(shotDelta)), shotCount)
		fmt.Println("")
		return true
	} else if isGameOverMan(targetRange, deathRadius) {
		return true
	} else {
		if shotDelta > 0.0 {
			fmt.Printf("<< Undershot target by %s.\n", getDisplayText(-shotDelta))
		} else {
			fmt.Printf(">> Overshot target by %s.\n", getDisplayText(-shotDelta))
		}
		printImpactTimeline(shotRange, false)
	}
	return false
}

func isGameOverMan(targetRange, deathRadius float64) bool {
	if targetRange <= deathRadius {
		fmt.Println("")
		fmt.Println(gameOverMan)
		fmt.Println("")
		return true
	}
	return false
}

func takeShot(shotCount int, shotAngle, projectileVmps float64) (shotRange, shotTime, shotDelta float64) {
	shotRange, shotTime = xRange(shotAngle, projectileVmps)
	fmt.Printf("Taking shot #%d at %4.2f degrees. Flight time is %3.1f seconds.\n", shotCount, shotAngle, shotTime)
	if targetModeAuto {
		// Wait here so that the target has time to move in targetMovement() during the shot.
		time.Sleep(time.Second * time.Duration(shotTime) / time.Duration(targetSpeedMultiplier))
	} else {
		// Fast forward the target to the correct location.
		targetRange -= targetVmps * shotTime
	}
	fmt.Printf("Shot #%d took %3.1f seconds, and went %s (%3.1f %s).\n", shotCount, shotTime, getDisplayText(shotRange), getMilesOrKilometers(shotRange, englishUnits), milesOrKilometers[englishUnits])
	shotDelta = targetRange - shotRange
	return
}

func getNextShotAngle(reader io.Reader) float64 {
	scanner := bufio.NewScanner(reader)
	for {
		fmt.Printf("Enter a shot angle from %3.1f to %3.1f degrees (0 to quit): ", minShotAngle, maxShotAngle)
		scanner.Scan()
		input := scanner.Text()
		if shotAngle, err := strconv.ParseFloat(input, 64); err == nil {
			if shotAngle >= 0.0 && shotAngle <= maxShotAngle {
				return shotAngle
			}
		}
		fmt.Printf("  Invalid Value: `%s`\n", input)
	}
}

func predictNextShotAngle(shotRange, shotTime, shotDelta float64) float64 {
	predictedLocation := shotRange + shotDelta - ((targetVmps * 0.95) * (shotTime * 0.95))
	predictedAngle := xAngle(predictedLocation, projectileVmps)
	if math.IsNaN(predictedAngle) {
		predictedAngle = maxShotAngle
	}
	return math.Max(math.Min(predictedAngle, maxShotAngle), minShotAngle)
}

func battleManager() {
	defer wg.Done()

	shotAngle := 0.0
	predictedShotAngle := maxShotAngle / 2.0
	shotCount := 0
	for {
		printHeader()
		if shootModeAuto {
			shotAngle = predictedShotAngle
		} else {
			shotAngle = getNextShotAngle(os.Stdin)
			if shotAngle == 0.0 {
				return
			}
		}
		shotCount++
		shotRange, shotTime, shotDelta := takeShot(shotCount, shotAngle, projectileVmps)
		if printImpactResults(shotRange, targetRange, shotDelta, deathRadius, shotCount) {
			return
		}
		predictedShotAngle = predictNextShotAngle(shotRange, shotTime, shotDelta)
	}
}

func targetMovement() {
	defer wg.Done()

	movementCount := 0
	for {
		time.Sleep(time.Second / time.Duration(targetSpeedMultiplier))
		targetRange -= targetVmps
		movementCount++
		if (movementCount % 10) == 0 {
			note := ""
			if targetSpeedMultiplier > 1 {
				note = fmt.Sprintf(" (at %dx real-time)", targetSpeedMultiplier)
			}
			fmt.Printf("Target Range = %s after %d seconds%s.\n", getDisplayText(targetRange), 10, note)
		}
		if isGameOverMan(targetRange, deathRadius) {
			return
		}
	}
}

func main() {
	initialize()
	wg.Add(1)
	if targetModeAuto {
		go targetMovement()
	}
	go battleManager()
	wg.Wait()
}
