package main

import (
	"your_package_path/zimgo"
)

func main() {
	// Open a .zim file (replace "your_file.zim" with the actual file path)
	z, err := zimgo.Open("your_file.zim")
	if err != nil {
		panic(err)
	}
	defer z.Close()

	// Initialize Gin router
	router := zimgo.SetupRouter(z)

	// Start the web server on port 8080
	router.Run(":8080")
}
