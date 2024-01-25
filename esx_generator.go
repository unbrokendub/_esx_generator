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
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {

	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	arguments := os.Args[1:]
	argumentOne := arguments[1]   //duration
	argumetnTwo := arguments[0]   //number of slices
	argumentThree := arguments[2] //keyword
	argumentFour := arguments[3]  //cut start time (random or not)

	slicesnumber, _ := strconv.Atoi(argumetnTwo)

	durationfloat, _ := strconv.ParseFloat(argumentOne, 64)

	letters := []string{"A", "B", "C", "D", "E", "F", "G", "H", "T", "U", "R", "E", "W", "Q", "X", "V", "N", "M", "K", "L", "O", "P", "J", "S", "T"}
	// Provide the path to your audio file
	//filePath := "cs60_vibra.wav"

	//fmt.Println(audio_duration(filePath))

	//fmt.Println("---------------------------")

	// err := filepath.Walk("/Users/unbrokendub/Desktop/AUDIO/SAMPLES",
	// 	func(path string, info os.FileInfo, err error) error {
	// 		if err != nil {
	// 			return err
	// 		}

	// 		if strings.Contains(strings.ToLower(path), ".wav") || strings.Contains(strings.ToLower(path), ".mp3") || strings.Contains(strings.ToLower(path), ".aiff") {
	// 			if strings.Contains(strings.ToLower(path), ".asd") {
	// 				//do nothing
	// 			} else {
	// 				fmt.Println(path)
	// 			}
	// 		}

	// 		return nil
	// 	})
	// if err != nil {
	// 	fmt.Println("error")
	// }

	readFile, err := os.Open("list.txt")
	if err != nil {
		fmt.Println(err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())

	}
	readFile.Close()
	fmt.Println(fileLines[3])

	for i := 0; i < slicesnumber; {
		randomNumber := rand.Intn(79416)

		if strings.Contains(strings.ToLower(fileLines[randomNumber]), argumentThree) {
			auDuration := audio_duration(fileLines[randomNumber])
			fmt.Println(auDuration)
			if auDuration > durationfloat {
				ffmpegCut(fileLines[randomNumber], strconv.Itoa(randomNumber), auDuration, argumentOne, argumentFour)
				i++
			}
		}

	}

	f, err := os.Create("list128.txt")
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

	newfilename := "/Users/unbrokendub/Desktop/AUDIO/SAMPLES/esx_generator/" + argumetnTwo + "_" + argumentOne + "_" + argumentThree + "_" + letters[rand.Intn(25)] + letters[rand.Intn(25)] + letters[rand.Intn(25)] + letters[rand.Intn(25)] + letters[rand.Intn(25)] + "_" + strconv.Itoa(rand.Intn(10000)) + ".wav"

	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", "list128.txt", "-c", "copy", newfilename)
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
