// Testing the utils functionality

package utils
import "testing"


func TestRound2(t *testing.T){
    var testCases = []struct {
        name string
        input float64
        want float64
    }{
        {"Zero Value", 0.00, 0.00},
        {"Float rounding", 1.56789, 1.57},
        {"Float round to zero", 0.00375, 0.00},
    }

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            if got  := Round2(test.input); got != test.want {
                t.Errorf("got: %g, wanted: %g", got, test.want)
            }
        })
    }
}

func TestRound4(t *testing.T){
    var testCases = []struct {
        name string
        input float64
        want float64
    }{
        {"Zero Value", 0.0000, 0.0000},
        {"Float rounding", 1.56789, 1.5679},
        {"Float round to zero", 0.00375, 0.0038},
    }

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            if got  := Round4(test.input); got != test.want {
                t.Errorf("got: %g, wanted: %g", got, test.want)
            }
        })
    }
}
