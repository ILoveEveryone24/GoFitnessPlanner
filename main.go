package main

import(
	"fmt"
	"strings"
	"bufio"
	"os"
	"os/exec"
	"runtime"
	"log"
	"encoding/json"
)

const(
	Reset = "\033[0m"
	Red = "\033[31m"
	Green = "\033[32m"
	Yellow = "\033[33m"
)

type Set struct{
	Reps int `json:"reps"`
	Weight float32 `json:"weight"`
}

type Exercise struct{
	Name string `json:"name"`
	Sets []Set `json:"sets"`
}

type Program struct{
	Name string `json:"name"`
	Exercises []Exercise `json:"exercises"`
}

func clear() {
        var cmd *exec.Cmd
        if runtime.GOOS == "windows" {
            cmd = exec.Command("cmd", "/c", "cls")
        } else {
            cmd = exec.Command("clear")
        }

        cmd.Stdout = os.Stdout
        cmd.Run()
}

func addProgramToFile(program Program){
	file, err := os.OpenFile("program.json", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil{
		log.Fatal(err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(program)
	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("Program has been added to file")
}

func addNewProgram(){
	var name string
	
	fmt.Println("Name of the program: ")

	fmt.Scan(&name)
	
	program := Program{
		Name: name,
		Exercises: make([]Exercise, 0),
	}
	addProgramToFile(program)

	fmt.Printf("Successfully created a fitness program with the name: %v %v %v\n\n", Green, program.Name, Reset)
}

func listPrograms(input int) int{
	file, err := os.OpenFile("program.json", os.O_RDWR, 0644)
	if err != nil{
		if os.IsNotExist(err){
			fmt.Println(Red +"No programs were created! Please, first create a program." + Reset)
		}else{
			log.Fatal(err)
		}
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	cnt := 1

	for scanner.Scan(){
		line := scanner.Text()

		if strings.TrimSpace(line) == ""{
			continue
		}

		var program Program

		err := json.Unmarshal([]byte(line), &program)
		if err != nil{
			log.Fatal(err)
		}
	
		if input > 0 && cnt == input{
			fmt.Println("Program Name: ", program.Name)
			if len(program.Exercises) == 0 {
				fmt.Println("No exercises found.\n")
			}else{
				fmt.Println("Exercises: ")
				for _, exercise := range program.Exercises{
					fmt.Printf("	Exercise: %v\n", exercise.Name)
				}
			}
		} else if input == 0{
			fmt.Print(cnt)
			fmt.Println(". Program Name: ", program.Name)
		}
		cnt++

	}

	if err := scanner.Err(); err != nil{
		log.Fatal(err)
	}

	return cnt
}

func editProgram(){
	cnt := listPrograms(0)
	
	if cnt <= 0{
		return
	}

	var input int
	
	fmt.Println("Which program would you like to edit?\n")

	for{
		fmt.Scan(&input)

		if input < 1 || input > cnt - 1{
			fmt.Println("Invalid program")
		} else{
			clear()
			listPrograms(input)		
		}
	}
}

func greet(){
	clear()
	var input int

	fmt.Println("Welcome to GoFitnessPlanner!")
	for{
		fmt.Println("What would you like to do?")
		fmt.Println("1: New program \n2: Edit program \n3: Exit")

		fmt.Scan(&input)
		if input == 1{
			clear()
			addNewProgram()
		} else if input == 2{
			clear()
			editProgram()
		} else if input == 3{
			return
		} else{
			clear()
			fmt.Println("Invalid input\n")
		}
	}
	fmt.Println("You've entered: ", input)
	
}

func main(){
	greet()
}
