package main

import (
	"fmt"
	"sync"
	"time"
)

type Food struct {
	id         string
	boiling    int
	fryEggs    int
	toastBread int
}

func (f Food) BoilWater(wg *sync.WaitGroup) error {
	defer wg.Done()
	fmt.Println("Boiling Water", f.id)
	time.Sleep(time.Duration(f.boiling) * time.Millisecond)
	return nil
}
func (f Food) FryEggs(wg *sync.WaitGroup) error {
	defer wg.Done()
	fmt.Println("Frying Eggs", f.id)
	time.Sleep(time.Duration(f.fryEggs) * time.Millisecond)
	return nil
}
func (f Food) ToastBread(wg *sync.WaitGroup) error {
	defer wg.Done()
	fmt.Println("Toasting Bread", f.id)
	time.Sleep(time.Duration(f.toastBread) * time.Millisecond)
	return nil
}

func ProcessCooking(food *Food) error {
	wg := new(sync.WaitGroup)
	wg.Add(3)
	go func() {
		err := food.BoilWater(wg)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		err := food.FryEggs(wg)
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		err := food.ToastBread(wg)
		if err != nil {
			fmt.Println(err)
		}
	}()
	wg.Wait()
	return nil
}

func main() {

	foods := []*Food{
		{id: "Palov", boiling: 3, fryEggs: 4, toastBread: 1},
		{id: "Soup", boiling: 5, fryEggs: 0, toastBread: 1},
		{id: "Kebab", boiling: 0, fryEggs: 1, toastBread: 1},
	}
	start := time.Now()
	for _, food := range foods {
		err := ProcessCooking(food)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	fmt.Println("Time elapsed: ", time.Since(start))
	fmt.Println("All cooking is completed")
}
