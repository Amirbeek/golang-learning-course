package main

import "fmt"

func task1() {
	num := 10
	ptr := &num
	fmt.Println("num:", num, "ptr:", ptr, "*ptr:", *ptr)
	*ptr = 100
	fmt.Println("num:", num, "ptr:", ptr, "*ptr:", *ptr)
}

type Car struct {
	brand string
	speed int
}

func changeSpeed(c *Car, newSpeed int) {
	c.speed = newSpeed
}
func task2() {
	car := Car{
		brand: "Porse",
		speed: 100,
	}
	fmt.Println("Before car:", car)
	changeSpeed(&car, 200)
	fmt.Println("After car:", car)

}

func main() {
	//task1()
	task2()

}
