package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jszwec/csvutil"
)

type stringFlag struct {
	set   bool
	value string
}

type hour struct {
	Hour1 string `csv:"hour1"`
	Hour2 string `csv:"hour2"`
	Hour3 string `csv:"hour3"`
	Hour4 string `csv:"hour4"`
	Hour5 string `csv:"hour5"`
	Hour6 string `csv:"hour6"`
	Hour7 string `csv:"hour7"`
	Hour8 string `csv:"hour8"`
}

func main() {
	var inVal stringFlag
	flag.Var(&inVal, "in", "help message for in")
	shouldPrint := flag.Bool("print", false, "help message for print")
	shouldUndo := flag.Bool("undo", false, "help message for undo")
	showNext := flag.Bool("next", false, "help message for next")
	flag.Parse()

	// load hours csv
	hours := loadHoursFile()

	if *showNext {
		fmt.Printf("next punch should be at:\n")
		fmt.Printf("00:00 to complete 8 hours\n")
		fmt.Printf("00:00 at maximum\n")

		c := askForConfirmation("Place notification for next punch?")
		if c {
			commandString := `sleep 4; osascript -e 'display notification "Punch in now" with title "PUNCH" sound name "Glass"'`
			exec.Command("sh", "-c", commandString).Start()
		}
	}

	if *shouldUndo {
		c := askForConfirmation("Do you really want to undo the last punch?")
		if c {
			undoLastPunch(hours)
			// write out
			writeHoursFile(hours)
		}
	}

	if inVal.set {
		if inVal.value == "now" {
			fmt.Printf("--punching in now\n")
		} else {
			fmt.Printf("--punching in at %q\n", inVal.value)
		}

		setNextPunch(hours, inVal.value)

		// write out
		writeHoursFile(hours)
	}

	if *shouldPrint {
		printHoursTable(hours)
	}
}

func (sf *stringFlag) Set(x string) error {
	sf.value = x
	sf.set = true
	return nil
}

func (sf *stringFlag) String() string {
	return sf.value
}

func printHoursTable(hours []hour) {
	fmt.Printf("-------------------------------------------------------------------------------\n")
	fmt.Printf("| dia | hora 1 | hora 2 | hora 3 | hora 4 | hora 5 | hora 6 | hora 7 | hora 8 |\n")
	fmt.Printf("|-----|--------|--------|--------|--------|--------|--------|--------|--------|\n")

	for i, day := range hours {
		fmt.Printf("| %02d  |", i+1)
		fmt.Printf(" %+v  | %+v  | %+v  | %+v  | %+v  | %+v  | %+v  | %+v  |",
			formatHour(day.Hour1), formatHour(day.Hour2), formatHour(day.Hour3),
			formatHour(day.Hour4), formatHour(day.Hour5), formatHour(day.Hour6),
			formatHour(day.Hour7), formatHour(day.Hour8))

		fmt.Printf("\n")
	}
	fmt.Printf("-------------------------------------------------------------------------------\n")
}

func formatHour(hour string) string {
	if hour != "" {
		return hour
	} else {
		return "     "
	}
}

func undoLastPunch(hours []hour) {
	t := time.Now()
	d := t.Day()

	switch nilString := ""; nilString {
	case hours[d-1].Hour2:
		hours[d-1].Hour1 = ""
	case hours[d-1].Hour3:
		hours[d-1].Hour2 = ""
	case hours[d-1].Hour4:
		hours[d-1].Hour3 = ""
	case hours[d-1].Hour5:
		hours[d-1].Hour4 = ""
	case hours[d-1].Hour6:
		hours[d-1].Hour5 = ""
	case hours[d-1].Hour7:
		hours[d-1].Hour6 = ""
	case hours[d-1].Hour8:
		hours[d-1].Hour7 = ""
	default:
		// hours[d+1].Hour8 = newHour
	}
}

func setNextPunch(hours []hour, hour string) {
	t := time.Now()
	d := t.Day()
	var newHour = hour

	if hour == "now" {
		newHour = t.Format("15:04")
	}

	fmt.Printf("%v\n", newHour)

	switch nilString := ""; nilString {
	case hours[d-1].Hour1:
		hours[d-1].Hour1 = newHour
	case hours[d-1].Hour2:
		hours[d-1].Hour2 = newHour
	case hours[d-1].Hour3:
		hours[d-1].Hour3 = newHour
	case hours[d-1].Hour4:
		hours[d-1].Hour4 = newHour
	case hours[d-1].Hour5:
		hours[d-1].Hour5 = newHour
	case hours[d-1].Hour6:
		hours[d-1].Hour6 = newHour
	case hours[d-1].Hour7:
		hours[d-1].Hour7 = newHour
	case hours[d-1].Hour8:
		hours[d-1].Hour8 = newHour
	default:
		// hours[d+1].Hour8 = newHour
	}
}

func loadHoursFile() []hour {
	csvInput, err := ioutil.ReadFile("hours.csv")
	if err != nil {
		fmt.Print(err)
	}

	var hours []hour
	if err := csvutil.Unmarshal(csvInput, &hours); err != nil {
		fmt.Println("error:", err)
	}

	return hours
}

func writeHoursFile(hours []hour) {
	b, err := csvutil.Marshal(hours)
	if err != nil {
		fmt.Println("error:", err)
	}

	err = ioutil.WriteFile("hours.csv", b, 0644)
	if err != nil {
		panic(err)
	}
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else {
			return false
		}
	}
}
