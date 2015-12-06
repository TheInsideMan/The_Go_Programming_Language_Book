package main

import (
	"fmt"
	"math"
	"net/http"
)

const (
	width, height = 600, 320
	cells         = 100
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 8
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func main() {
	http.HandleFunc("/test", handler)
	http.ListenAndServe("localhost:8000", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Content-Encoding", "gzip")

	fmt.Println(`<?xml version="1.0" encoding="iso-8859-1"?>
	<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.0//EN"
	"http://www.w3.org/TR/2001/
	REC-SVG-20010904/DTD/svg10.dtd">`)

	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill:white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			errors := 0
			ax, ay, height, ok := corner(i+1, j)
			if !ok {
				errors++
			}
			bx, by, height, ok := corner(i, j)
			if !ok {
				errors++
			}
			cx, cy, height, ok := corner(i, j+1)
			if !ok {
				errors++
			}
			dx, dy, height, ok := corner(i+1, j+1)
			if !ok {
				errors++
			}
			if errors < 1 {
				color := ""
				if height < 0 {
					color = "#0000ff"
				} else {
					color = GetColor(float64(height))
				}
				fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g' style='stroke:%v'/>\n", ax, ay, bx, by, cx, cy, dx, dy, color)
			}
		}
	}
	fmt.Println("</svg>")
}

func GetColor(h float64) string {
	colors := []string{"#0000ff", "#3300ff", "#6600ff", "#9900ff", "#cc00ff", "#ff00ff", "#ff00cc", "#ff0099", "#ff0033", "#ff0000"}
	h *= 10
	h2 := int(h)
	return colors[h2]
}

func corner(i, j int) (float64, float64, float64, bool) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	z, ok := f(x, y)
	if ok {
		sx := width/2 + (x-y)*cos30*xyscale
		sy := height/2 + (x+y)*sin30*xyscale - z*zscale
		return sx, sy, z, true
	}
	return 0.0, 0.0, 0.0, false
}

func f(x, y float64) (float64, bool) {
	r := math.Hypot(x, y) // distance from 0,0
	nan := math.NaN()
	if r == nan {
		return 0.0, false
	}
	return math.Sin(r), true
}

// import (
// 	"flag"
// 	"fmt"
// 	"popcount"
// 	"strings"
// )

// var n = flag.Bool("n", false, "omit trailing newline") // & pointer
// var sep = flag.String("s", " ", "seperator")

// func main() {
// 	x, s := popcount.PopCount(uint64(10))

// 	fmt.Println(x)
// 	fmt.Println(s)
// 	flag.Parse()
// 	fmt.Print(strings.Join(flag.Args(), *sep))
// 	if !*n {
// 		fmt.Println()
// 	}
// }

// import (
// 	"fmt"
// 	"io/ioutil"
// 	// "bytes"
// 	"image"
// 	"image/color"
// 	"image/gif"
// 	"io"
// 	"log"
// 	"math"
// 	"math/rand"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"sync"
// 	"time"
// )

// func fetch(url string, ch chan<- string) {
// 	start := time.Now()
// 	if !strings.HasPrefix(url, "http://") {
// 		fmt.Println("Please prefix URL with 'http://'")
// 		os.Exit(1)
// 	}
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
// 		os.Exit(1)
// 	}
// 	b, err := io.Copy(ioutil.Discard, resp.Body)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
// 		os.Exit(1)
// 	}
// 	secs := time.Since(start).Seconds()
// 	ch <- fmt.Sprintf("%.2fs %7d %s", secs, b, url)
// }

// var mu sync.Mutex
// var count int

// func counter(w http.ResponseWriter, r *http.Request) {
// 	mu.Lock()
// 	fmt.Fprintf(w, "Count %d\n", count)
// 	mu.Unlock()
// }

// const (
// 	whiteIndex = 0
// 	blackIndex = 1
// )

// func lissajous(out io.Writer, cycles float64, red uint8, green uint8, blue uint8, size int) {
// 	var palette = []color.Color{color.White, color.RGBA{red, green, blue, 100}}
// 	const (
// 		// cycles  = cycles
// 		res = 0.001
// 		// size    = 200
// 		nframes = 64
// 		delay   = 8
// 	)
// 	if cycles < 1 {
// 		cycles = 1
// 	}
// 	if size < 1 {
// 		size = 64
// 	}
// 	freq := rand.Float64() * 3.0
// 	anim := gif.GIF{LoopCount: nframes}
// 	phase := 0.0
// 	for i := 0; i < nframes; i++ {
// 		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
// 		img := image.NewPaletted(rect, palette)
// 		for t := 0.0; t < cycles*2*math.Pi; t += res {
// 			x := math.Sin(t)
// 			y := math.Sin(t*freq + phase)
// 			img.SetColorIndex(size+int(x*float64(size)+0.5), size+int(y*float64(size)+0.5), 1)
// 		}
// 		phase += 0.1
// 		anim.Delay = append(anim.Delay, delay)
// 		anim.Image = append(anim.Image, img)
// 	}
// 	gif.EncodeAll(out, &anim)
// }

// func handler(w http.ResponseWriter, r *http.Request) {
// 	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
// 	red, _ := strconv.Atoi(r.URL.Query().Get("r"))
// 	green, _ := strconv.Atoi(r.URL.Query().Get("g"))
// 	blue, _ := strconv.Atoi(r.URL.Query().Get("b"))
// 	cycles := r.URL.Query().Get("cycles")
// 	c, _ := strconv.Atoi(cycles)
// 	lissajous(w, float64(c), uint8(red), uint8(green), uint8(blue), int(size))
// 	mu.Lock()
// 	count++
// 	mu.Unlock()
// 	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
// }

// func main() {
// 	http.HandleFunc("/", handler)
// 	http.HandleFunc("/count", counter)
// 	// http.HandleFunc("/lissajous", lissajous)
// 	log.Fatal(http.ListenAndServe("localhost:8000", nil))
// 	// start := time.Now()
// 	// ch := make(chan string)
// 	// for _, url := range os.Args[1:] {
// 	// 	go fetch(url, ch)
// 	// }
// 	// for range os.Args[1:] {
// 	// 	fmt.Println(<-ch)
// 	// }
// 	// fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())

// }

// import (
// 	"image"
// 	"image/color"
// 	"image/gif"
// 	"io"
// 	"math"
// 	"math/rand"
// 	"os"
// 	// "io/ioutil"
// 	// "bufio"
// 	// "fmt"
// 	// "strconv"
// 	// "strings"
// 	// "time"
// )

// var palette = []color.Color{color.Black, color.RGBA{63, 250, 6, 100}, color.RGBA{24, 6, 250, 100}}

// func main() {
// 	lissajous(os.Stdout)
// }

// func main() {
// 	counts := make(map[string]int)
// 	for _, filename := range os.Args[1:] {
// 		data, err := ioutil.ReadFile(filename)
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
// 			continue
// 		}
// 		for _, line := range strings.Split(string(data), "\n") {
// 			counts[line]++
// 		}
// 	}
// 	fmt.Println(len(counts))
// 	for line, n := range counts {
// 		if n > 1 {
// 			fmt.Printf("%d\t%s\n", n, line)
// 		}
// 	}
// }

// func main() {
// 	counts := make(map[string]int)
// 	files := os.Args[1:]
// 	// fmt.Printf("file %v\n", files)
// 	if len(files) == 0 {
// 		countLines(os.Stdin, counts, "")
// 	} else {
// 		// fmt.Println("here I am")
// 		for _, arg := range files {
// 			f, err := os.Open(arg)
// 			if err != nil {
// 				// fmt.Println("some error")
// 				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
// 				continue //go to next iteration of loop
// 			}

// 			countLines(f, counts, arg)
// 			f.Close()
// 		}
// 	}
// 	for line, n := range counts {
// 		// fmt.Printf("got lines - %v\n", n)
// 		// fmt.Printf("%v - %v\n", n, line)
// 		if n > 1 {
// 			fmt.Printf("%d\t%s\n", n, line)
// 		}

// 	}
// }
// func countLines(f *os.File, counts map[string]int, fname string) {
// 	// fmt.Println("inside countLines()")
// 	intput := bufio.NewScanner(f)
// 	for intput.Scan() {
// 		counts[intput.Text()]++
// 		if fname != "" && counts[intput.Text()] > 1 {
// 			fmt.Printf("file - %v\n", fname)
// 		}
// 	}
// }

// func main() {
// 	counts := make(map[string]int)
// 	input := bufio.NewScanner(os.Stdin)
// 	for input.Scan() {
// 		counts[input.Text()]++
// 	}
// 	for line, n := range counts {
// 		// if n > 1 {
// 		fmt.Printf("%d\t%s\n", n, line)
// 		// }
// 	}
// }

// func main() {
// 	start := time.Now()
// 	// s, sep := "", ""
// 	// strings.Join(os.Args[1:], " ")
// 	// for i, arg := range os.Args[1:] {
// 	// 	s += sep + "\n" + strconv.Itoa(i) + " - " + arg + " - " + os.Args[0]

// 	// 	sep = " "
// 	// }
// 	// fmt.Println(s)
// 	fmt.Println(strings.Join(os.Args[1:], " "))
// 	fmt.Printf("%v elapsed time\n", time.Since(start))
// }
