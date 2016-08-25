package main

import "fmt"
import "io/ioutil"
import "os"
import "strings"

const PATH = "/sys/bus/usb/devices"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing command line argument")
		return
	}
	args := os.Args[1:]

	if args[0] == "list" {
		var authorized []string
		var unauthorized []string

		files, err := ioutil.ReadDir(PATH)
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			path := PATH + "/" + file.Name()

			product_file := path + "/product"
			if _, err := os.Stat(product_file); os.IsNotExist(err) {
				// Not actually a device.
				continue
			}
			product_contents, err := ioutil.ReadFile(product_file)
			if err != nil {
				panic(err)
			}
			product := strings.TrimSpace(string(product_contents))

			authorized_file := path + "/authorized"
			authorized_contents, err := ioutil.ReadFile(authorized_file)
			if err != nil {
				panic(err)
			}
			if authorized_contents[0] == 49 {
				authorized = append(authorized, product+" (ID: "+file.Name()+")")
				//fmt.Println(file.Name() + " - authorized")
			} else {
				unauthorized = append(unauthorized, product+" (ID: "+file.Name()+")")
				//fmt.Println(file.Name() + " - unauthorized")
			}
		}

		fmt.Println("Authorized devices")
		fmt.Println()
		for _, device := range authorized {
			fmt.Println("\t" + device)
		}
		fmt.Println()
		fmt.Println("Unauthorized devices")
		fmt.Println()
		for _, device := range unauthorized {
			fmt.Println("\t" + device)
		}
	}
	if args[0] == "enable" {
		if len(args) < 2 {
			fmt.Println("enable requires a device to enable")
			return
		}
		device := args[1]

		path := PATH + "/" + device
		authorized_file := path + "/authorized"
		authorized_contents, err := ioutil.ReadFile(authorized_file)
		if err != nil {
			panic(err)
		}
		if authorized_contents[0] == 49 {
			fmt.Println("Device already enabled")
			return
		}

		err = ioutil.WriteFile(authorized_file, []byte{49}, 0)
		if err != nil {
			panic(err)
		}
	}
	if args[0] == "disable" {
		if len(args) < 2 {
			fmt.Println("enable requires a device to enable")
			return
		}
		device := args[1]

		path := PATH + "/" + device
		authorized_file := path + "/authorized"
		authorized_contents, err := ioutil.ReadFile(authorized_file)
		if err != nil {
			panic(err)
		}
		if authorized_contents[0] == 48 {
			fmt.Println("Device already disabled")
			return
		}

		err = ioutil.WriteFile(authorized_file, []byte{48}, 0)
		if err != nil {
			panic(err)
		}

	}
	if args[0] == "default" {
		fmt.Println("default configurations are not yet supported")
		if len(args) < 2 {
			return
		}
	}
}
