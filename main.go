package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/0xAX/notificator"
	"github.com/jszwec/csvutil"
)

type stringFlag struct {
	set   bool
	value string
}

type workDay struct {
	Hour1 string `csv:"hour1"`
	Hour2 string `csv:"hour2"`
	Hour3 string `csv:"hour3"`
	Hour4 string `csv:"hour4"`
	Hour5 string `csv:"hour5"`
	Hour6 string `csv:"hour6"`
	Hour7 string `csv:"hour7"`
	Hour8 string `csv:"hour8"`
}

var notify *notificator.Notificator

func main() {
	// init notificator
	notify = notificator.New(notificator.Options{
		DefaultIcon: "icon/default.png",
		AppName:     "Puncher",
	})

	// parse flags
	var inVal stringFlag
	flag.Var(&inVal, "in", "help message for in")
	shouldPrint := flag.Bool("print", false, "help message for print")
	shouldUndo := flag.Bool("undo", false, "help message for undo")
	showNext := flag.Bool("next", false, "help message for next")
	flag.Parse()

	// load hours csv
	days := loadHoursFile()

	if *showNext {
		fmt.Printf("next punch should be at:\n")
		fmt.Printf("00:00 to complete 8 hours\n")
		fmt.Printf("00:00 at maximum\n")

		c := askForConfirmation("Place notification for next punch?")
		if c {
			notify.Push("title", "text", "/home/user/icon.png", notificator.UR_CRITICAL)
			// commandString := `sleep 4; osascript -e 'display notification "Punch in now" with title "PUNCH" sound name "Glass"'`
			// exec.Command("sh", "-c", commandString).Start()
		}
	}

	if *shouldUndo {
		c := askForConfirmation("Do you really want to undo the last punch?")
		if c {
			undoLastPunch(days)
			// write out
			writeHoursFile(days)
		}
	}

	if inVal.set {
		if inVal.value == "now" {
			fmt.Printf("--punching in now\n")
		} else {
			fmt.Printf("--punching in at %q\n", inVal.value)
		}

		setNextPunch(days, inVal.value)

		// write out
		writeHoursFile(days)
	}

	if *shouldPrint {
		printHoursTable(days)
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

func printHoursTable(days []workDay) {
	fmt.Printf("-------------------------------------------------------------------------------\n")
	fmt.Printf("| dia | hora 1 | hora 2 | hora 3 | hora 4 | hora 5 | hora 6 | hora 7 | hora 8 |\n")
	fmt.Printf("|-----|--------|--------|--------|--------|--------|--------|--------|--------|\n")

	for i, day := range days {
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

func getCurrentPunch(days []workDay) *string {
	d := time.Now().Day()
	var punch *string

	switch nilString := ""; nilString {
	case days[d-1].Hour2:
		punch = &days[d-1].Hour1
	case days[d-1].Hour3:
		punch = &days[d-1].Hour2
	case days[d-1].Hour4:
		punch = &days[d-1].Hour3
	case days[d-1].Hour5:
		punch = &days[d-1].Hour4
	case days[d-1].Hour6:
		punch = &days[d-1].Hour5
	case days[d-1].Hour7:
		punch = &days[d-1].Hour6
	case days[d-1].Hour8:
		punch = &days[d-1].Hour7
	default:
		punch = &days[d-1].Hour8
	}

	return punch
}

func getNextPunch(days []workDay) *string {
	d := time.Now().Day()
	var punch *string

	switch nilString := ""; nilString {
	case days[d-1].Hour1:
		punch = &days[d-1].Hour1
	case days[d-1].Hour2:
		punch = &days[d-1].Hour2
	case days[d-1].Hour3:
		punch = &days[d-1].Hour3
	case days[d-1].Hour4:
		punch = &days[d-1].Hour4
	case days[d-1].Hour5:
		punch = &days[d-1].Hour5
	case days[d-1].Hour6:
		punch = &days[d-1].Hour6
	case days[d-1].Hour7:
		punch = &days[d-1].Hour7
	case days[d-1].Hour8:
		punch = &days[d-1].Hour8
	default:
		punch = &days[d-1].Hour8
	}

	return punch
}

func undoLastPunch(days []workDay) {
	punch := getCurrentPunch(days)
	*punch = ""
	fmt.Printf("%+v\n", punch)
}

func setNextPunch(days []workDay, hour string) {
	t := time.Now()
	var newHour = hour

	if hour == "now" {
		newHour = t.Format("15:04")
	}

	punch := getNextPunch(days)
	*punch = newHour
}

func loadHoursFile() []workDay {
	csvInput, err := ioutil.ReadFile("hours.csv")
	if err != nil {
		fmt.Print(err)
	}

	var days []workDay
	if err := csvutil.Unmarshal(csvInput, &days); err != nil {
		fmt.Println("error:", err)
	}

	return days
}

func writeHoursFile(days []workDay) {
	b, err := csvutil.Marshal(days)
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
