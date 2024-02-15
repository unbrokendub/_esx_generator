package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/vansante/go-ffprobe"
)

func randFloats(min, max float64, n int) []float64 {
	source := rand.NewSource(time.Now().UnixNano())
	rndm := rand.New(source)

	res := make([]float64, n)
	for i := range res {
		res[i] = min + rndm.Float64()*(max-min)
	}
	return res
}

func audio_duration(filename string) float64 {
	info, err := ffprobe.GetProbeData(filename, 5*time.Second)
	if err != nil {
		fmt.Printf("Error getting probe data: %s\n", err)

	}

	// Print the duration of the audio file
	duration := info.Format.DurationSeconds
	return duration
}

func ffmpegCut(input string, ouptut string, duration float64, cuttime string, starttime string) {
	ouptut = "/Users/unbrokendub/projects/esx_generator/tmp/" + ouptut + ".wav"

	seek := randFloats(0, duration, 1)
	seek_string := "0"

	if starttime == "r" {
		seek_string = fmt.Sprintf("%f", seek[0])
	}

	fmt.Println("ffmpeg", "-ss", seek_string, "-i", input, "-vn", "-acodec", "pcm_s16le", "-ar", "44100", "-ac", "1", "-t", cuttime, ouptut)

	cmd := exec.Command("ffmpeg", "-ss", seek_string, "-i", input, "-vn", "-acodec", "pcm_s16le", "-ar", "44100", "-ac", "1", "-t", cuttime, ouptut)
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {

	}
}

func main() {
	lines_number := 0

	source := rand.NewSource(time.Now().UnixNano())
	rndm := rand.New(source)
	arguments := os.Args[1:]
	argDuration := arguments[1]       //duration
	argNumberOfSlices := arguments[0] //number of slices
	argKeyword := arguments[2]        //keyword
	aargCutStartTime := arguments[3]  //cut start time (random or not)

	slicesnumber, _ := strconv.Atoi(argNumberOfSlices)

	durationfloat, _ := strconv.ParseFloat(argDuration, 64)

	letters := []string{"A", "B", "C", "D", "E", "F", "G", "H", "T", "U", "R", "E", "W", "Q", "X", "V", "N", "M", "K", "L", "O", "P", "J", "S", "T"}

	sam, err := os.Create("sampleslist.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file

	err = filepath.Walk("/Users/unbrokendub/Desktop/AUDIO/SAMPLES",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if strings.Contains(strings.ToLower(path), ".wav") || strings.Contains(strings.ToLower(path), ".mp3") || strings.Contains(strings.ToLower(path), ".aiff") {
				if strings.Contains(strings.ToLower(path), ".asd") {
					//do nothing
				} else {
					_, err2 := sam.WriteString(path + "\n")
					if err != nil {
						log.Fatal(err2)
					}
				}
			}

			return nil
		})
	if err != nil {
		fmt.Println("error")
	}

	sam.Close()

	readFile, err := os.Open("sampleslist.txt")
	if err != nil {
		fmt.Println(err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		if strings.Contains(strings.ToLower(fileScanner.Text()), argKeyword) {
			fileLines = append(fileLines, fileScanner.Text())
			lines_number++
		}
	}

	readFile.Close()

	rndm.Shuffle(len(fileLines), func(i, j int) {
		fileLines[i],
			fileLines[j] = fileLines[j],
			fileLines[i]
	})

	fmt.Println("Number of lines:", lines_number)
	for i := 0; i < slicesnumber; {
		randomNumber := rndm.Intn(lines_number)

		if strings.Contains(strings.ToLower(fileLines[randomNumber]), argKeyword) {
			fmt.Println("Random number is ", randomNumber)
			fmt.Println("File is ", fileLines[randomNumber])
			auDuration := audio_duration(fileLines[randomNumber])
			fmt.Println(auDuration)
			if auDuration > durationfloat {
				ffmpegCut(fileLines[randomNumber], strconv.Itoa(i), auDuration, argDuration, aargCutStartTime)
				i++
			}
			fmt.Println("===============================================================")
		}

	}

	f, err := os.Create("tmp_list.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file
	defer f.Close()

	err1 := filepath.Walk("/Users/unbrokendub/projects/esx_generator/tmp",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			//fmt.Println(path)
			if strings.Contains(path, ".wav") {
				_, err2 := f.WriteString("file " + `'` + path + `'` + "\n")
				if err != nil {
					log.Fatal(err2)
				}
			}
			return nil
		})
	if err1 != nil {
		fmt.Println("error")
	}

	newfilename := "/Users/unbrokendub/Desktop/AUDIO/SAMPLES/esx_generator/" + argNumberOfSlices + "_" + argDuration + "_" + argKeyword + "_" + letters[rndm.Intn(25)] + letters[rndm.Intn(25)] + letters[rndm.Intn(25)] + letters[rndm.Intn(25)] + letters[rndm.Intn(25)] + "_" + strconv.Itoa(rndm.Intn(10000)) + ".wav"

	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", "tmp_list.txt", "-c", "copy", newfilename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err3 := cmd.Run()
	if err3 != nil {

	}

	cmd2 := exec.Command("rm", "-rf", "tmp")
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	err4 := cmd2.Run()
	if err4 != nil {

	}

	cmd3 := exec.Command("mkdir", "tmp")
	cmd3.Stdout = os.Stdout
	cmd3.Stderr = os.Stderr
	err5 := cmd3.Run()
	if err5 != nil {

	}

}
