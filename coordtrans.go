package coordtrans

import (
	"math"
)

const pi  = 3.14159265358979324
const a = 6378245.0
const ee = 0.00669342162296594323
const x_pi = pi * 3000.0 / 180.0

type Location struct {
	Lat float64
	Lng float64
}

func out_of_china(wgs Location) bool {
	if wgs.Lng < 72.004 || wgs.Lng > 137.8347 {
		return true
	}
	if wgs.Lat < 0.8293 || wgs.Lat > 55.8271 {
		return true
	}
	return false
}

func transform_lat(x, y float64) float64 {
	ret := -100.0 + 2.0 * x + 3.0 * y + 0.2 * y * y + 0.1 * x * y + 0.2 * math.Sqrt(math.Abs(x))
	ret += (20.0 * math.Sin(6.0 * x * pi) + 20.0 * math.Sin(2.0 * x * pi)) * 2.0 / 3.0
	ret += (20.0 * math.Sin(y * pi) + 40.0 * math.Sin(y / 3.0 * pi)) * 2.0 / 3.0
	ret += (160.0 * math.Sin(y / 12.0 * pi) + 320 * math.Sin(y * pi / 30.0)) * 2.0 / 3.0
	return ret
}

func transform_lng(x, y float64) float64 {
	ret := 300.0 + x + 2.0 * y + 0.1 * x * x + 0.1 * x * y + 0.1 * math.Sqrt(math.Abs(x))
	ret += (20.0 * math.Sin(6.0 * x * pi) + 20.0 * math.Sin(2.0 * x * pi)) * 2.0 / 3.0
	ret += (20.0 * math.Sin(x * pi) + 40.0 * math.Sin(x / 3.0 * pi)) * 2.0 / 3.0
	ret += (150.0 * math.Sin(x / 12.0 * pi) + 300.0 * math.Sin(x / 30.0 * pi)) * 2.0 / 3.0
	return ret
}

func Wgs2gcj(wgs Location) Location {
	if out_of_china(wgs) {
		return wgs
	}
	d_lat := transform_lat(wgs.Lng - 105.0, wgs.Lat - 35.0)
	d_lng := transform_lng(wgs.Lng - 105.0, wgs.Lat - 35.0)
	rad_lat := wgs.Lat / 180.0 * pi
	magic := math.Sin(rad_lat)
	magic = 1 - ee * magic * magic
	sqrt_magic := math.Sqrt(magic)
	d_lat = (d_lat * 180.0) / ((a * (1 - ee)) / (magic * sqrt_magic) * pi)
	d_lng = (d_lng * 180.0) / (a / sqrt_magic * math.Cos(rad_lat) * pi)
	gcj := Location{
		Lat: wgs.Lat + d_lat,
		Lng: wgs.Lng + d_lng,
	}
	return gcj
}

func Gcj2wgs(gcj Location) Location {
	g_pt := Wgs2gcj(gcj)
	d_lng := g_pt.Lng - gcj.Lng
	d_lat := g_pt.Lat - gcj.Lat
	wgs := Location{
		Lat: gcj.Lat - d_lat,
		Lng: gcj.Lng - d_lng,
	}
	return wgs
}

func gcj2bd(gcj Location) Location {
	x := gcj.Lng
	y := gcj.Lat
	z := math.Sqrt(x * x + y * y) + 0.00002 * math.Sin(y * x_pi)
	theta := math.Atan2(y, x) + 0.000003 * math.Cos(x * x_pi)
	bd_lat := z * math.Cos(theta) + 0.0065
	bd_lng := z * math.Sin(theta) + 0.006
	bd := Location{
		Lat: bd_lat,
		Lng: bd_lng,
	}
	return bd
}

func bd2gcj(bd Location) Location {
	bd_lat := bd.Lat
	bd_lng := bd.Lng
	x := bd_lng - 0.0065
	y := bd_lat - 0.006
	z := math.Sqrt(x * x + y * y) - 0.00002 * math.Sin(y * x_pi)
	theta := math.Atan2(y, x) - 0.000003 * math.Cos(x * x_pi)
	gcj_lng := z * math.Cos(theta)
	gcj_lat := z * math.Sin(theta)
	gcj := Location{
		Lat: gcj_lat,
		Lng: gcj_lng,
	}
	return gcj
}

func Wgs2bd(wgs Location) Location {
	return gcj2bd(Wgs2gcj(wgs))
}

func Bd2wgs(bd Location) Location {
	return Gcj2wgs(bd2gcj(bd))
}
