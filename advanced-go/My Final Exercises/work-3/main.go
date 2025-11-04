package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Pomodoro struct {
	timerName string
	workTime  time.Duration
	breakTime time.Duration
}

var cancelGoroutines = make(map[string]func())
var mu sync.Mutex

func NewPomodoro(name string, w time.Duration, b time.Duration) *Pomodoro {
	return &Pomodoro{
		timerName: name,
		workTime:  w * time.Second,
		breakTime: b * time.Second,
	}
}

func (p *Pomodoro) Start(ctx context.Context) {
	start := time.Now()
	endWork := start.Add(p.workTime)
	breakEnd := endWork.Add(p.breakTime)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	workDone := false

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("The %s Pomodoro is stopped.\n", p.timerName)
			mu.Lock()
			delete(cancelGoroutines, p.timerName)
			mu.Unlock()
			return
		default:
			now := time.Now()

			// Study time finished
			if !workDone && now.After(endWork) {
				fmt.Printf("The %s Pomodoro's study time finished.\n", p.timerName)
				fmt.Printf("Break time has started!\n")
				workDone = true
			}

			// Pomodoro finished
			if now.After(breakEnd) {
				fmt.Printf("The %s Pomodoro is Finished.\n", p.timerName)
				mu.Lock()
				delete(cancelGoroutines, p.timerName)
				mu.Unlock()
				return
			}
		}
	}
}

func ask(command string) {
	cm := strings.Fields(command)
	if len(cm) == 0 {
		return
	}

	cmd := strings.ToLower(cm[0])

	if cmd == "exit" {
		os.Exit(0)
	} else if cmd == "time" {
		if len(cm) < 4 {
			fmt.Println("Usage: time <name> <work_duration> <break_duration>")
			return
		}

		name := cm[1]
		pomodoroNum, err1 := strconv.Atoi(cm[2])
		breakNum, err2 := strconv.Atoi(cm[3])
		if err1 != nil || err2 != nil {
			fmt.Println("Invalid durations")
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		p := NewPomodoro(name, time.Duration(pomodoroNum), time.Duration(breakNum))
		go p.Start(ctx)
		mu.Lock()
		cancelGoroutines[p.timerName] = cancel
		mu.Unlock()
		fmt.Printf("Timer '%s' started!\n", p.timerName)

	} else if cmd == "stop" {
		if len(cm) < 2 {
			fmt.Println("Usage: stop <name>")
			return
		}
		value, ok := cancelGoroutines[cm[1]]
		if !ok {
			fmt.Println("Value not found")
			return
		}
		fmt.Printf("Timer stopped with value: %v\n", cm[1])
		value()
		mu.Lock()
		delete(cancelGoroutines, cm[1])
		mu.Unlock()

	} else if cmd == "list" {
		keys := make([]string, 0, len(cancelGoroutines))
		for k := range cancelGoroutines {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for i, k := range keys {
			fmt.Printf("%d: %s\n", i+1, k)
		}
	} else if cmd == "clear" {
		ClearTerminal()
		return
	} else {
		PrintS()
	}
}
func ClearTerminal() {
	fmt.Print("\n\n\n\n\n\n\n\n\n\n\n\n\n\\")
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	PrintS()
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		input = strings.TrimSpace(input)
		ask(input)
	}
}

func PrintS() {
	fmt.Println("Please choose an option:")
	fmt.Println("  - Stop system          : [Exit]")
	fmt.Println("  - Set timer            : [time]")
	fmt.Println("  - List timers          : [list]")
	fmt.Println("  - Clear the Window     : [clear]")
	fmt.Println("  - Stop a specific timer: [stop name]")
}
