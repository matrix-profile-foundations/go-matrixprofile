[![Build Status](https://travis-ci.com/matrix-profile-foundation/go-matrixprofile.svg?branch=master)](https://travis-ci.com/matrix-profile-foundation/go-matrixprofile)
[![codecov](https://codecov.io/gh/matrix-profile-foundation/go-matrixprofile/branch/master/graph/badge.svg)](https://codecov.io/gh/matrix-profile-foundation/go-matrixprofile)
[![Go Report Card](https://goreportcard.com/badge/github.com/matrix-profile-foundation/go-matrixprofile)](https://goreportcard.com/report/github.com/matrix-profile-foundation/go-matrixprofile)
[![GoDoc](https://godoc.org/github.com/matrix-profile-foundation/go-matrixprofile?status.svg)](https://godoc.org/github.com/matrix-profile-foundation/go-matrixprofile)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

# go-matrixprofile

Golang library for computing a matrix profiles and matrix profile indexes. Features also include time series discords, time series segmentation, and motif discovery after computing the matrix profile. Visit [The UCR Matrix Profile Page](https://www.cs.ucr.edu/~eamonn/MatrixProfile.html) for more details into matrix profiles.

Features:
* STMP
* STAMP (parallelized)
* STAMPI
* STOMP (parallelized)
* mSTOMP
* MPX (parallelized)
* TopKMotifs - finds the top K motifs from a computed matrix profile
* TopKDiscords - finds the top K discords from a computed matrix profile
* Segement - computes the corrected arc curve for time series segmentation
* Annotation Vectors
  * Complexity
  * Mean Standard Deviation
  * Clipping

## Contents
- [Installation](#installation)
- [Quick start](#quick-start)
- [Case Studies](#case-studies)
  * [Matrix Profile](#matrix-profile)
  * [Multi-Dimensional Matrix Profile](#multi-dimensional-matrix-profile)
- [Benchmarks](#benchmarks)
- [Contributing](#contributing)
- [Testing](#testing)
- [Other Libraries](#other-libraries)
- [Contact](#contact)
- [License](#license)
- [Citations](#citations)

## Installation
```sh
$ go get -u github.com/matrix-profile-foundation/go-matrixprofile
$ cd $GOPATH/src/github.com/matrix-profile-foundation/go-matrixprofile
$ make setup
```

## Quick start
```sh
$ cat example_test.go
```
```go
package main

import (
	"fmt"

	"github.com/matrix-profile-foundation/go-matrixprofile"
)

func main() {
	sig := []float64{0, 0.99, 1, 0, 0, 0.98, 1, 0, 0, 0.96, 1, 0}

	mp, err := matrixprofile.New(sig, nil, 4)
	if err != nil {
		panic(err)
	}

	if err = mp.Compute(matrixprofile.NewOptions()); err != nil {
		panic(err)
	}

	fmt.Printf("Signal:         %.3f\n", sig)
	fmt.Printf("Matrix Profile: %.3f\n", mp.MP)
	fmt.Printf("Profile Index:  %5d\n", mp.Idx)
}
```
```sh
$ go run example_mp.go
Signal:         [0.000 0.990 1.000 0.000 0.000 0.980 1.000 0.000 0.000 0.960 1.000 0.000]
Matrix Profile: [0.014 0.014 0.029 0.029 0.014 0.014 0.029 0.029 0.029]
Profile Index:  [    4     5     6     7     0     1     2     3     4]
```

## Case studies
### Matrix Profile
Going through a completely synthetic scenario, we'll cover what features to look for in a matrix profile, and what the additional Discords, TopKMotifs, and Segment tell us. We'll first be generating a fake signal that is composed of sine waves, noise, and sawtooth waves. We then run STOMP on the signal to calculte the matrix profile and matrix profile indexes.

![mpsin](https://github.com/matrix-profile-foundation/go-matrixprofile/blob/master/mp_sine.png)
subsequence length: 32

* signal: This shows our raw data. Theres several oddities and patterns that can be seen here. 
* matrix profile: generated by running STOMP on this signal which generates both the matrix profile and the matrix profile index. In the matrix profile we see several spikes which indicate that these may be time series discords or anomalies in the time series.
* corrected arc curve: This shows the segmentation of the time series. The two lowest dips around index 420 and 760 indicate potential state changes in the time series. At 420 we see the sinusoidal wave move into a more pulsed pattern. At 760 we see the pulsed pattern move into a sawtooth pattern.
* discords: The discords graph shows the top 3 potential discords of the defined subsequence length, m, based on the 3 highest peaks in the matrix profile. This is mostly composed of noise.
* motifs: These represent the top 6 motifs found from the time series. The first being the initial sine wave pattern. The second is during the pulsed sequence on a fall of the pulse to the noise. The third is during the pulsed sequence on the rise from the noise to the pulse. The fourth and fifth are the sawtooth patterns.

The code to generate the graph can be found in [this example](https://github.com/matrix-profile-foundation/go-matrixprofile/blob/master/matrixprofile/example_caseStudy_test.go#L121).

### Multi-Dimensional Matrix Profile
Based on [4] we can extend the matrix profile algorithm to multi-dimensional scenario.

![mpkdim](https://github.com/matrix-profile-foundation/go-matrixprofile/blob/master/mp_kdim.png)
subsequence length: 25

* signal 0-2: the 3 time series dimensions
* matrix profile 0-2: the k-dimensional matrix profile representing choose k from d time series. matrix profile 1 minima represent motifs that span at that time across 2 time series of the 3 available. matrix profile 2 minima represents the motifs that span at that time across 3 time series.

The plots can be generated by running
```sh
$ make example
go test ./... -run=Example
ok  	github.com/matrix-profile-foundation/go-matrixprofile	0.256s
ok  	github.com/matrix-profile-foundation/go-matrixprofile/av	(cached) [no tests to run]
?   	github.com/matrix-profile-foundation/go-matrixprofile/method	[no test files]
ok  	github.com/matrix-profile-foundation/go-matrixprofile/siggen	(cached) [no tests to run]
ok  	github.com/matrix-profile-foundation/go-matrixprofile/util	(cached) [no tests to run]
```
A png file will be saved in the top level directory of the repository as `mp_sine.png` and `mp_kdim.png`

## Benchmarks
Benchmark name                      | NumReps |    Time/Rep    |  Memory/Rep  |     Alloc/Rep   |
-----------------------------------:|--------:|---------------:|-------------:|----------------:|
BenchmarkMStomp-4                   |       50|  28559842 ns/op|  7335193 B/op| 227071 allocs/op|
BenchmarkZNormalize-4               | 10000000|       159 ns/op|      256 B/op|      1 allocs/op|
BenchmarkMovmeanstd-4               |    50000|     26689 ns/op|    65537 B/op|      4 allocs/op|
BenchmarkCrossCorrelate-4           |    10000|    138444 ns/op|    49180 B/op|      3 allocs/op|
BenchmarkMass-4                     |    10000|    144664 ns/op|    49444 B/op|      4 allocs/op|
BenchmarkDistanceProfile-4          |    10000|    147884 ns/op|    49444 B/op|      4 allocs/op|
BenchmarkCalculateDistanceProfile-4 |   200000|      8959 ns/op|        2 B/op|      0 allocs/op|
BenchmarkStmp/m32_pts1k-4           |        5| 300883003 ns/op| 97396006 B/op|   7883 allocs/op|
BenchmarkStmp/m128_pts1k-4          |        5| 297214441 ns/op| 94091148 B/op|   7498 allocs/op|
BenchmarkStamp/m32_p2_pts1k-4       |       10| 193601139 ns/op| 97498281 B/op|   7898 allocs/op|
BenchmarkStomp/m_32_p1_pts1024-4    |       50|  42061763 ns/op|   156119 B/op|     16 allocs/op|
BenchmarkStomp/m128_p1_pts1024-4    |       50|  37975960 ns/op|   156124 B/op|     16 allocs/op|
BenchmarkStomp/m128_p2_pts1024-4    |      100|  24062106 ns/op|   302562 B/op|     25 allocs/op|
BenchmarkStomp/m128_p2_pts2048-4    |       20|  99007571 ns/op|   638441 B/op|     26 allocs/op|
BenchmarkStomp/m128_p2_pts4096-4    |       10| 403500289 ns/op|  1318713 B/op|     27 allocs/op|
BenchmarkStomp/m128_p2_pts8192-4    |       10|1775433560 ns/op|  2616211 B/op|     27 allocs/op|
BenchmarkStomp/m128_p4_pts8192-4    |       10|1742241625 ns/op|  4992480 B/op|     48 allocs/op|
BenchmarkMpx/m_32_p1_pts1024-4      |       50|  14109571 ns/op|   137484 B/op|     15 allocs/op|
BenchmarkMpx/m128_p1_pts1024-4      |       50|  22611401 ns/op|   137484 B/op|     15 allocs/op|
BenchmarkMpx/m128_p2_pts1024-4      |      100|  21189096 ns/op|   167374 B/op|     19 allocs/op|
BenchmarkMpx/m128_p2_pts2048-4      |       20|  63714124 ns/op|   359726 B/op|     20 allocs/op|
BenchmarkMpx/m128_p2_pts4096-4      |       10| 203365439 ns/op|   777936 B/op|     21 allocs/op|
BenchmarkMpx/m128_p2_pts8192-4      |       10| 703186642 ns/op|  1551251 B/op|     21 allocs/op|
BenchmarkMpx/m128_p4_pts8192-4      |       10| 640190693 ns/op|  2075916 B/op|     29 allocs/op|
BenchmarkStampUpdate-4              |       10| 173804923 ns/op|  2031161 B/op|     24 allocs/op|

Ran on a 2018 MacBookAir on Dec 18, 2019
```sh
    Processor: 1.6 GHz Intel Core i5
       Memory: 8GB 2133 MHz LPDDR3
           OS: macOS Mojave v10.14.2
 Logical CPUs: 4
Physical CPUs: 2
```
```sh
$ make bench
```

## Contributing
* Fork the repository
* Create a new branch (feature_\* or bug_\*)for the new feature or bug fix
* Run tests
* Commit your changes
* Push code and open a new pull request

## Testing
Run all tests including benchmarks
```sh
$ make all
```
Just run benchmarks
```sh
$ make bench
```
Just run tests
```sh
$ make test
```

## Other libraries
* R: [github.com/franzbischoff/tsmp](https://github.com/franzbischoff/tsmp)
* Python: [github.com/target/matrixprofile-ts](https://github.com/target/matrixprofile-ts)

## Contact
* Austin Ouyang (aouyang1@gmail.com)

## License
The MIT License (MIT). See [LICENSE](https://github.com/matrix-profile-foundation/go-matrixprofile/blob/master/LICENSE) for more details.

Copyright (c) 2018 Austin Ouyang

## Citations
[1] Chin-Chia Michael Yeh, Yan Zhu, Liudmila Ulanova, Nurjahan Begum, Yifei Ding, Hoang Anh Dau, Diego Furtado Silva, Abdullah Mueen, Eamonn Keogh (2016). [Matrix Profile I: All Pairs Similarity Joins for Time Series: A Unifying View that Includes Motifs, Discords and Shapelets](https://www.cs.ucr.edu/~eamonn/PID4481997_extend_Matrix%20Profile_I.pdf). IEEE ICDM 2016.

[2] Yan Zhu, Zachary Zimmerman, Nader Shakibay Senobari, Chin-Chia Michael Yeh, Gareth Funning, Abdullah Mueen, Philip Berisk and Eamonn Keogh (2016). [Matrix Profile II: Exploiting a Novel Algorithm and GPUs to break the one Hundred Million Barrier for Time Series Motifs and Joins](https://www.cs.ucr.edu/~eamonn/STOMP_GPU_final_submission_camera_ready.pdf). IEEE ICDM 2016.

[3] Hoang Anh Dau and Eamonn Keogh (2017). [Matrix Profile V: A Generic Technique to Incorporate Domain Knowledge into Motif Discovery](https://www.cs.ucr.edu/~eamonn/guided-motif-KDD17-new-format-10-pages-v005.pdf). KDD 2017.

[4] Chin-Chia Michael Yeh, Nickolas Kavantzas, Eamonn Keogh (2017).[Matrix Profile VI: Meaningful Multidimensional Motif Discovery](https://www.cs.ucr.edu/%7Eeamonn/Motif_Discovery_ICDM.pdf). ICDM 2017.

[5] Shaghayegh Gharghabi, Yifei Ding, Chin-Chia Michael Yeh, Kaveh Kamgar, Liudmila Ulanova, Eamonn Keogh (2017). [Matrix Profile VIII: Domain Agnostic Online Semantic Segmentation at Superhuman Performance Levels](https://www.cs.ucr.edu/%7Eeamonn/Segmentation_ICDM.pdf). ICDM 2017.
