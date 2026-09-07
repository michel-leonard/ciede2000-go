// This function written in Go is not affiliated with the CIE (International Commission on Illumination),
// and is released into the public domain. It is provided "as is" without any warranty, express or implied.

package main

import "math"

// Convenience function, with parametric factors set to their default values.
func ciede2000(l1 float64, a1 float64, b1 float64, l2 float64, a2 float64, b2 float64) float64 {
	return ciede2000_with_parameters(l1, a1,b1, l2, a2, b2, 1.0, 1.0, 1.0, false)
}

// The classic CIE ΔE2000 implementation, which operates on two L*a*b* colors, and returns their difference.
// "l" ranges from 0 to 100, while "a" and "b" are unbounded and commonly clamped to the range of -128 to 127.
func ciede2000_with_parameters(l1 float64, a1 float64, b1 float64, l2 float64, a2 float64, b2 float64, kl float64, kc float64, kh float64, canonical bool) float64 {
	// Working in Go with the CIEDE2000 color-difference formula.
	// kl, kc, kh are parametric factors to be adjusted according to
	// different viewing parameters such as textures, backgrounds...
	n := (math.Sqrt(a1 * a1 + b1 * b1) + math.Sqrt(a2 * a2 + b2 * b2)) * 0.5
	n = n * n * n * n * n * n * n
	// A factor involving chroma raised to the power of 7 designed to make
	// the influence of chroma on the total color difference more accurate.
	n = 1.0 + 0.5 * (1.0 - math.Sqrt(n / (n + 6103515625.0)))
	// Application of the chroma correction factor.
	c1 := math.Sqrt(a1 * a1 * n * n + b1 * b1)
	c2 := math.Sqrt(a2 * a2 * n * n + b2 * b2)
	// atan2 is preferred over atan because it accurately computes the angle of
	// a point (x, y) in all quadrants, handling the signs of both coordinates.
	h1 := math.Atan2(b1, a1 * n)
	h2 := math.Atan2(b2, a2 * n)
	if h1 < 0.0 { h1 += 2.0 * math.Pi }
	if h2 < 0.0 { h2 += 2.0 * math.Pi }
	// When the hue angles lie in different quadrants, the straightforward
	// average can produce a mean that incorrectly suggests a hue angle in
	// the wrong quadrant, the next lines handle this issue.
	h_mean := (h1 + h2) * 0.5
	h_delta := (h2 - h1) * 0.5
	// The part where most programmers get it wrong.
	if math.Pi + 1E-14 < math.Abs(h2 - h1) {
		h_delta += math.Pi
		if canonical && math.Pi + 1E-14 < h_mean {
			// Sharma’s implementation, OpenJDK, ...
			h_mean -= math.Pi
		} else {
			// Lindbloom’s implementation, Netflix’s VMAF, ...
			h_mean += math.Pi
		}
	}
	p := 36.0 * h_mean - 55.0 * math.Pi
	n = (c1 + c2) * 0.5
	n = n * n * n * n * n * n * n
	// The hue rotation correction term is designed to account for the
	// non-linear behavior of hue differences in the blue region.
	r_t :=	-2.0 * math.Sqrt(n / (n + 6103515625.0)) *
			math.Sin(math.Pi / 3.0 * math.Exp(p * p / (-25.0 * math.Pi * math.Pi)))
	n = (l1 + l2) * 0.5
	n = (n - 50.0) * (n - 50.0)
	// Lightness.
	l := (l2 - l1) / (kl * (1.0 + 0.015 * n / math.Sqrt(20.0 + n)))
	// These coefficients adjust the impact of different harmonic
	// components on the hue difference calculation.
	t := 1.0 -	0.17 * math.Sin(h_mean + math.Pi / 3.0) +
				0.24 * math.Sin(2.0 * h_mean + math.Pi * 0.5) +
				0.32 * math.Sin(3.0 * h_mean + 8.0 * math.Pi / 15.0) -
				0.20 * math.Sin(4.0 * h_mean + 3.0 * math.Pi / 20.0)
	n = c1 + c2
	// Hue.
	h := 2.0 * math.Sqrt(c1 * c2) * math.Sin(h_delta) / (kh * (1.0 + 0.0075 * n * t))
	// Chroma.
	c := (c2 - c1) / (kc * (1.0 + 0.0225 * n))
	// The result reflects the actual geometric distance in the color space, given a tolerance of 3.6e-13.
	return math.Sqrt(l * l + h * h + c * c + c * h * r_t)
}

// If you remove the constant 1E-14, the code will continue to work, but CIEDE2000
// interoperability between all programming languages will no longer be guaranteed.

// Source code tested by Michel LEONARD
// Website: ciede2000.pages-perso.free.fr
