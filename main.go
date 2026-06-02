package main

import "fmt"

type User struct {
	Name string
	Age  int
}

func main() {
	user := User{}

	createUser(&user.Name, &user.Age)

	fmt.Println("user:", user)
}

func createUser(name *string, age *int) {
	*name = "mohamad"
	*age = 2
}
