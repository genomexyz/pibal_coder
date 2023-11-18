package main

import (
	"fmt"
	"math"
)

func calcX(az, el, h float64) float64 {
	B_rad := math.Pi/2 - az*math.Pi/180
	C_rad := el * math.Pi / 180

	// Calculate the result
	x := math.Cos(B_rad) * (h / math.Tan(C_rad))
	return x
}

func calcY(az, el, h float64) float64 {
	B_rad := math.Pi/2 - az*math.Pi/180
	C_rad := el * math.Pi / 180

	// Calculate the result
	y := math.Sin(B_rad) * (h / math.Tan(C_rad))
	return y
}

func calcWindspeed(x1, x2, y1, y2 float64) float64 {
	result := math.Sqrt(math.Pow(x2-x1, 2)+math.Pow(y2-y1, 2)) / 120
	return result
}

func calcWinddir(x1, x2, y1, y2 float64) float64 {
	condition := x2-x1 > 0
	dir := 0.0
	if condition {
		dir = (4.712 - math.Atan((y2-y1)/(x2-x1))) / (math.Pi / 180)
	} else {
		dir = (4.712 - math.Atan((y2-y1)/(x2-x1)) - 3.1416) / (math.Pi / 180)
	}
	return dir
}

func calcSpeedDir(elevation1, azimuth1, elevation2, azimuth2, h1, h2 float64) (float64, float64) {
	x1 := calcX(azimuth1, elevation1, h1)
	x2 := calcX(azimuth2, elevation2, h2)
	y1 := calcY(azimuth1, elevation1, h1)
	y2 := calcY(azimuth2, elevation2, h2)

	speed := calcWindspeed(x1, x2, y1, y2)
	dir := calcWinddir(x1, x2, y1, y2)

	return speed, dir
}

func main() {
	// Example usage
	//elevation1 := 30.0
	//azimuth1 := 45.0
	//elevation2 := 25.0
	//azimuth2 := 60.0
	//h1 := 1000.0
	//h2 := 1200.0

	az1 := 230.3
	az2 := 223.
	el1 := 28.2
	el2 := 27.9
	h1 := 153.
	h2 := 433.

	speed, dir := calcSpeedDir(el1, az1, el2, az2, h1, h2)
	fmt.Printf("Wind Speed: %f m/s\nWind Direction: %f degrees\n", speed, dir)
}
