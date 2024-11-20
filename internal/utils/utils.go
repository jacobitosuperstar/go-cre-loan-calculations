// Internal Utilities
// This is used to round the numbers to certain ammount of decimals.
package utils
import "math"


// Round2 returns a float number rounded up to 2 decimals. This is mainly used
// for values that represent money.
func Round2(num float64) float64 {
    return (math.Round(num*100)/100)
}

// Round4 returns a float number rounded up to 4 decimals. This is mainly used
// for decimal values that represent percentages.
func Round4(num float64) float64 {
    return (math.Round(num*10000)/10000)
}
