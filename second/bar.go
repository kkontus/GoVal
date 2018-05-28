package second

import "fmt"

func bar() {
	fmt.Println("Bar Internal Function")
}

func BarExported() {
	fmt.Println("Bar Exported Function")
}
