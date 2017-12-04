/*
 * Copyright (c) 2017. Braedon Wooding
 * Created under LICENSE, see the file LICENSE for information
 */

package tests

import (
	"bufio"
	"bytes"
	"doml/parser"
	"fmt"
	"math"
	"os"
	"os/user"
	"strings"
	"testing"

	chart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/seq"
)

var singularBlock = `
@ example =System.Cool...
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
`

var block string

const ratioValue = float64(0.01106186)

func TestAccuracyParser(t *testing.T) {
	testBlock := func(startPoint, endPoint, incrementSize int) float64 {
		var avgAccuracy float64
		for blockSize := startPoint; blockSize <= endPoint; blockSize += incrementSize {
			block = strings.Repeat(singularBlock, blockSize)
			expectedTime := float64(len(block)) / ratioValue
			benchmarkSpeed := float64(testing.Benchmark(getSpeedForSingular).NsPerOp())
			avgAccuracy += math.Abs(expectedTime-benchmarkSpeed) / benchmarkSpeed
		}
		return (avgAccuracy / float64(endPoint/incrementSize-(startPoint/incrementSize-1))) * 100
	}

	fmt.Printf("%f%%\n", testBlock(1, 10, 1))
	fmt.Printf("%f%%\n", testBlock(50, 60, 1))
	fmt.Printf("%f%%\n", testBlock(100, 1000, 100))
}

func TestGetPlotForParserSpeed(t *testing.T) {
	maxLength := 10

	testBlock := func(size int) float64 {
		block = strings.Repeat(singularBlock, size)
		benchmarkSpeed := float64(testing.Benchmark(getSpeedForSingular).NsPerOp())
		return float64(len(block)) / (benchmarkSpeed) * 1000
	}

	usr, err := user.Current()
	if err != nil {
		println(err.Error())
		return
	}
	f, err := os.Create(usr.HomeDir + "/Desktop/Example.png")
	if err != nil {
		println(err.Error())
		return
	}
	defer f.Close()

	println(f.Name())

	values := make([]float64, maxLength)

	for i := 0; i < maxLength; i++ {
		values[i] = testBlock(i + 1)
	}

	ticks := make([]chart.Tick, maxLength)

	for i := 1; i < maxLength; i++ {
		ticks[i-1] = chart.Tick{Value: float64(i), Label: fmt.Sprintf("%v", i)}
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:  "Ratio * 1000",
			Ticks: ticks,
			Style: chart.Style{
				Show: true,
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
				XValues: seq.Range(0, float64(maxLength)-1),
				YValues: values,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err = graph.Render(chart.PNG, buffer)
	if err != nil {
		println(err.Error())
		return
	}

	writer := bufio.NewWriter(f)
	writer.Write(buffer.Bytes())
	writer.Flush()
}

func TestGetAverageParserSpeed(t *testing.T) {
	testBlock := func(startPoint, endPoint int) float64 {
		var avg float64
		for blockSize := startPoint; blockSize <= endPoint; blockSize++ {
			block = strings.Repeat(singularBlock, blockSize)
			benchmarkSpeed := float64(testing.Benchmark(getSpeedForSingular).NsPerOp())
			println(float64(len(block)) / benchmarkSpeed)
			avg += float64(len(block)) / benchmarkSpeed
		}
		return avg / float64(endPoint-(startPoint-1))
	}

	fmt.Printf("%f\n", testBlock(1, 5))
	fmt.Printf("%f\n", testBlock(10, 50))
	fmt.Printf("%f\n", testBlock(1, 100))
}

func getSpeedForSingular(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser.ParseDOMLFromString(block)
	}
}
