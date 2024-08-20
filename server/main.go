package main

import "fmt"

func main() {
	a := "id"

	switch a {
	case "ip":
	case "port":
	case "protocol":
	case "response_time":
	case "status_id":
	case "country_id":
	case "id":
		fmt.Println("id")
	default:
		fmt.Println("default")
	}
}
