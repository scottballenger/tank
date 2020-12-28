package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	maxProjectileV = 600  // meters/sec
	minProjectileV = 300  // meters/sec
	maxTargetV     = 60   // kilometers/hour
	minTargetV     = 0    // kilometers/hour
	impactRadius   = 20.0 // meters
	maxShotAngle   = 45.0 // degrees
	minShotAngle   = 1.0  // degrees

	gameOverMan = "GAME OVER MAN, you just got crushed by the tank!"
)

var (
	projectileV         float64
	targetV             float64
	targetVmps          float64
	maxRange            float64
	targetRange         float64
	deathRadius         float64
	wg                  sync.WaitGroup
	playModeAuto        bool // true = Auto Shoot Mode, false = Single Player
	playModePauseTarget bool // true = target pauses when deciding shot, false = target moves when deciding shot
	gameSpeed           int  // times faster than real-time
	englishUnits        bool // true = english units, false = metric units
)

func init() {
	flag.BoolVar(&playModeAuto, "a", false, "Auto Shoot Mode (default - Single Player)")
	flag.BoolVar(&playModePauseTarget, "p", false, "Pause Target During Shot Decision (default - Realtime Target Movement)")
	flag.BoolVar(&englishUnits, "e", false, "English Units (default - Metric)")
	flag.Float64Var(&deathRadius, "d", impactRadius, "Detonation Radius (meters)")
	flag.Parse()

	if playModeAuto {
		fmt.Println("Play Mode: Auto")
		gameSpeed = 10
		if playModePauseTarget {
			fmt.Println("Play Mode: Error: Pause Target During Shot Decision has no efffect during Auto Shot Mode")
		}
	} else {
		fmt.Println("Play Mode: Single Player")
		gameSpeed = 1
		if playModePauseTarget {
			fmt.Println("Play Mode: Pause Target During Shot Decision")
		} else {
			fmt.Println("Play Mode: Realtime Target Movement")
		}
	}

	if englishUnits {
		fmt.Println("Units: English")
	} else {
		fmt.Println("Units: Metric")
	}

	fmt.Printf("Detonation Radius = %s\n", getDisplayText(deathRadius))

	rand.Seed(time.Now().UnixNano())
	projectileV = minProjectileV + float64(rand.Intn(10000))*float64(maxProjectileV-minProjectileV)/10000.0
	targetV = minTargetV + float64(rand.Intn(10000))*float64(maxTargetV-minTargetV)/10000.0
	targetVmps = targetV * (1000.0 / 3600.0)
	maxRange, _ = xRange(45.0, projectileV)
	targetRange = float64(maxRange*0.2) + float64(rand.Intn(10000))*(maxRange*0.8)/10000.0
}

func xRange(angle, v float64) (x, t float64) {
	radians := (2 * math.Pi) * (angle / 360.0)
	t = (math.Sin(radians) * v) / 4.9
	x = math.Sin(radians) * t * v
	return
}

func getDisplayText(value float64) string {
	if englishUnits {
		return fmt.Sprintf("%3.1f feet", value*3.28084)
	}
	return fmt.Sprintf("%3.1f meters", value)
}

func printHeader() {
	fmt.Println("==================================")
	fmt.Printf("Projectile Velocity  = %s/sec\n", getDisplayText(projectileV))
	fmt.Printf("Max Projectile Range = %s\n", getDisplayText(maxRange))
	if englishUnits {
		fmt.Printf("Target Velocity      = %3.1f miles/hour\n", targetV*0.62137121212121)
	} else {
		fmt.Printf("Target Velocity      = %3.1f kilometers/hour\n", targetV)
	}
	fmt.Printf("Target Velocity      = %s/sec\n", getDisplayText(targetVmps))
	fmt.Printf("Current Target Range = %s\n", getDisplayText(targetRange))
	fmt.Println("----------------------------------")
}

func getImpactTimeline(shotDistance, targetDistance, maxDistance float64, hit bool) (impactString [2]string) {
	flightPath := 0
	impactPath := 1
	impactString[flightPath] = " /~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
	impactString[impactPath] = "/--------+---------+---------+---------+----------|"
	maxString := len(impactString[impactPath]) - 1
	targetIndex := int(targetDistance/maxDistance*float64(maxString)) + 1
	shotIndex := int(shotDistance/maxDistance*float64(maxString)) + 1
	if hit {
		if targetIndex < 4 {
			targetIndex = 4
		}
		impactString[flightPath] = impactString[flightPath][:targetIndex-2] + "\\"
		impactString[impactPath] = impactString[impactPath][:targetIndex-1] + "*" + impactString[impactPath][targetIndex:]
	} else {
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
		impactString[flightPath] = impactString[flightPath][:shotIndex-2] + "\\"
		impactString[impactPath] = impactString[impactPath][:shotIndex-1] + "\\" + impactString[impactPath][shotIndex:]
		impactString[impactPath] = impactString[impactPath][:targetIndex-1] + "T" + impactString[impactPath][targetIndex:]
	}
	return
}

func printImpactTimeline(shotDistance float64, hit bool) {

	impactStrings := getImpactTimeline(shotDistance, targetRange, maxRange, hit)

	fmt.Println("")
	for _, impactString := range impactStrings {
		fmt.Println(impactString)
	}
	fmt.Println("")
}

func printImpactResults(shotRange, shotDelta float64, shotCount int) bool {
	fmt.Printf("Target Range = %s at time of impact.\n", getDisplayText(targetRange))
	if math.Abs(shotDelta) <= deathRadius {
		printImpactTimeline(shotRange, true)
		fmt.Println("")
		fmt.Println("Direct hit after", shotCount, "shots!!")
		fmt.Println("")
		return true
	} else if targetRange <= deathRadius {
		fmt.Println("")
		fmt.Println(gameOverMan)
		fmt.Println("")
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

func singlePlayer() {
	shotCount := 0
	var shotAngle float64
	for {
		printHeader()

		goodValue := false
		for !goodValue {
			fmt.Printf("Enter a shot angle from %3.1f to %3.1f degrees:\n", minShotAngle, maxShotAngle)
			_, err := fmt.Scanf("%f\n", &shotAngle)
			if err != nil {
				fmt.Println(err)
			} else {
				if shotAngle >= 0.0 && shotAngle <= 45.0 {
					if shotAngle == 0.0 {
						return
					}
					goodValue = true
				}
			}
		}

		shotCount++
		shotRange, shotTime := xRange(shotAngle, projectileV)
		fmt.Printf("Taking shot #%d at %4.2f degrees.\n", shotCount, shotAngle)
		fmt.Printf("Shot #%d took %3.1f seconds, and went %s.\n", shotCount, shotTime, getDisplayText(shotRange))
		targetRange -= targetVmps * shotTime
		shotDelta := targetRange - shotRange
		if printImpactResults(shotRange, shotDelta, shotCount) {
			return
		}
	}
}

func targetMovement() {
	defer wg.Done()

	movementCount := 0
	for {
		time.Sleep(time.Second / time.Duration(gameSpeed))
		targetRange -= targetVmps
		movementCount++
		if (movementCount % 10) == 0 {
			fmt.Printf("Target Range = %s after %d seconds.\n", getDisplayText(targetRange), 10)
		}
		if targetRange < deathRadius {
			fmt.Println("")
			fmt.Println(gameOverMan)
			fmt.Println("")
			return
		}
	}
}

func battleManager() {
	defer wg.Done()

	shotAngle := maxShotAngle / 2.0
	shotCount := 0
	for {
		printHeader()
		shotCount++
		shotRange, shotTime := xRange(shotAngle, projectileV)
		fmt.Printf("Taking shot #%d at %4.2f degrees.\n", shotCount, shotAngle)
		time.Sleep(time.Second * time.Duration(shotTime) / time.Duration(gameSpeed))
		fmt.Printf("Shot #%d took %3.1f seconds, and went %s.\n", shotCount, shotTime, getDisplayText(shotRange))
		shotDelta := targetRange - shotRange
		if printImpactResults(shotRange, shotDelta, shotCount) {
			return
		}
		predictedShotDelta := shotDelta - (targetVmps * shotTime)
		shotAngleDelta := maxShotAngle * (predictedShotDelta / maxRange)
		shotAngle += shotAngleDelta
		shotAngle = math.Max(math.Min(shotAngle, maxShotAngle), minShotAngle)
	}
}

func main() {
	if playModeAuto {
		wg.Add(1)
		go targetMovement()
		go battleManager()
		wg.Wait()
	} else {
		if !playModePauseTarget {
			go targetMovement()
		}
		singlePlayer()
	}
}
