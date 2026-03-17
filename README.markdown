# CIEDE2000 color difference formula in Go

This page presents the CIEDE2000 color difference, implemented in the Go programming language.

![Logo for CIEDE2000 in Golang](https://raw.githubusercontent.com/michel-leonard/ciede2000-color-matching/refs/heads/main/docs/assets/images/logo.jpg)

## Our CIEDE2000 offer

This production-ready file, released in 2026, contain the CIEDE2000 algorithm.

Source File|Type|Bits|Purpose|Advantage|
|:--:|:--:|:--:|:--:|:--:|
[ciede2000.go](./ciede2000.go)|`float64`|64|General|Interoperability|

### Software Versions

- Go 1.X

### Example Usage

We calculate the CIEDE2000 distance between two colors, first without and then with parametric factors.

```go
func main() {
	// Example of two L*a*b* colors
	l1, a1, b1 := 78.7, 65.2, -2.9
	l2, a2, b2 := 77.5, 60.7, 2.8

	delta_e := ciede2000(l1, a1, b1, l2, a2, b2);
	println("CIEDE2000 =", delta_e); // ΔE2000 = 2.919895

	// Example of parametric factors used in the textile industry
	kl, kc, kh := 2.0, 1.0, 1.0

	// Perform a CIEDE2000 calculation compliant with that of Gaurav Sharma
	canonical := true;

	delta_e = ciede2000_with_parameters(l1, a1, b1, l2, a2, b2, kl, kc, kh, canonical);
	println("CIEDE2000 =", delta_e); // ΔE2000 = 2.826179
}
```

## Public Domain Licence

You are free to use these files, even for commercial purposes.
