package harvesine

import "math"

//EarthRadius is 6372.8
func ReferenceHarvesine(x0, y0, x1, y1, EarthRadius float64) float64 {
    lat1 := y0
    lat2 := y1
    lon1 := x0
    lon2 := x1

    dLat := radiansFromDegrees(lat2 - lat1) 
    dLon := radiansFromDegrees(lon2 - lon1) 
    lat1 = radiansFromDegrees(lat1)
    lat2 = radiansFromDegrees(lat2)

    a := Square(math.Sin(dLat/2.0) + math.Cos(lat1) * math.Cos(lat2) * Square(math.Sin(dLon/2))) 
    c := 2.0*math.Asin(math.Sqrt(a))

    return EarthRadius * c
}

func Square(A float64) float64 {
    return A * A
}

func radiansFromDegrees(degrees float64) float64 {
    return (math.Pi / 180) * degrees
}
