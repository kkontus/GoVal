package second

import "fmt"

func baz() {
	fmt.Println("Baz Internal Function")
}

func BazExported() {
	fmt.Println("Baz Exported Function")
}
