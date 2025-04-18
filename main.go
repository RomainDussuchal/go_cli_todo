package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)


type Task struct {
	Id     int    `json:"id"`
	Item   string `json:"item"`
	Status string `json:"status"` // "pending" or "done"
	CreatedAt string `json:"createdAt"`
	UpdatedAt string  `json:"UpdatedAt"`

}

const filename = "todo.json"

func loadTasks()([]Task, error){
	var tasks []Task
	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		return tasks, nil
	}
	
	defer file.Close()

	err = json.NewDecoder(file).Decode(&tasks)
	if err != nil {
		if err == io.EOF {
			return tasks, nil // empty file, return empty task list
		}
		return nil, err
	}
	return tasks, nil
}

func saveTasks(tasks *[]Task)error{
	file, err := os.Create(filename) // truncates if exists
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(tasks)

}


func main(){

	fmt.Println("-------------------------------------------")
	fmt.Println("Welcome to Todo CLI")
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Error Loading intials task file")
	}

	for {
		fmt.Println("\nAvailable commands: add | edit | list | timer | delete | delete -all | done | help | exit")
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		args := strings.Fields(strings.ToLower(input))
	
		switch args[0] {
			case "add":{
				handleAdd(&tasks)
			}
			case "edit": {
				handleEdit(&tasks)
			}
			case "list": {
				filter := ""
				if len(args) == 1 {
					fmt.Println("All items in the list:")
					handleList(&tasks, filter)
				} else if len(args) == 2 && (args[1] == "pending" || args[1] == "done") {
					filter = args[1]
					handleList(&tasks, filter)
				} else {
					fmt.Println("❌ Invalid usage. Try: list, list done, or list pending")
				}	
			}
			case "timer": {
				handleTimer()
			}
			
			case "delete": {
				if len(args) == 1 {
				handleDelete(&tasks)
				} else if 
					len(args) == 2 && (args[1] == "-all") {
						handleDeleteAll(&tasks)
					}
				}
		
			case "done": {
				handleDone(&tasks)
			}
		case "exit": {
			handleExit()
		}
		}
	}
}

func handleAdd(tasks *[]Task) {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	input = strings.TrimSpace(input)
	if input == "" {
		fmt.Println("⚠️ Task cannot be empty")
	}
	newTask := Task{
		Id: len(*tasks) + 1,
		Item: input,
		Status: "pending",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	*tasks = append(*tasks, newTask)

	fmt.Printf("Saving Task: %s", newTask.Item)

	saveTasks(tasks)
	fmt.Printf("✅ Task added: \"%s\"\n", newTask.Item)
}

func handleEdit(tasks *[]Task){
	fmt.Println("Select ID of the task to edit")
	fmt.Println("> ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	input = strings.TrimSpace(input)

	id, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Please enter a number")
		return
	}
	index := id-1

	fmt.Printf("Current task: %s\n", (*tasks)[index].Item)
	fmt.Println("Enter new task content:")
	fmt.Printf("> ")

	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	(*tasks)[index].Item = input
	(*tasks)[index].UpdatedAt = time.Now().Format(time.RFC3339)
	

	saveTasks(tasks)

	fmt.Printf("Task Id: %d properly edited: %s", id, input)
}

func handleList(tasks *[]Task, filter string) {
	
	if len(*tasks) == 0 {
		fmt.Println("No item in the list.")
	}

	if len(filter) > 0 {
	fmt.Printf("📋 Listing tasks")
	if filter != "" {
		fmt.Printf(" (filter: %s)", filter)
	}
	fmt.Println(":")
		
	for _, task := range *tasks {
		if filter == "done" && task.Status != "done" {
			continue
		}
		if filter == "pending" && task.Status != "pending" {
			continue
		}
		status := "[ ]"
		if task.Status == "done" {
			status = "[x]"
		}
		fmt.Printf("%d. %s %s\n", task.Id, status, task.Item)
	}
		return 
	}
	for _, task := range *tasks {
		
		status := "[ ]"
		if task.Status == "done" {
			status = "[x]"
		}
		fmt.Printf("%d. %s %s\n", task.Id, status, task.Item)
	}
}
func handleTimer(){
	fmt.Println("Start Timer")
	duration := 3 * time.Minute
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	start := time.Now()
	end := start.Add(duration)
	
	for now := range ticker.C {
		remaining := end.Sub(now)

		if remaining <= 0 {
			break
		}

		fmt.Printf("\r⏳ %s remaining", remaining.Truncate(time.Second))
	}

	fmt.Println("\n✅ Time's up!")
}

func handleDelete(tasks *[]Task){
	fmt.Print("Indicate task ID to delete :\n")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	id, err := strconv.Atoi(input)
	fmt.Printf("ID to delete : %d",id)

	if err != nil || id <= 0 {
		fmt.Println("⚠️ Invalid ID")
		return
	}
	index := -1
	for i, task := range *tasks {
		if task.Id == id {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("❌ Task not found")
		return
	}

	// Remove task
	*tasks = append((*tasks)[:index], (*tasks)[index+1:]...)

	// Re-assign IDs
	for i := range *tasks {
		(*tasks)[i].Id = i + 1 
	}
	saveTasks(tasks)
	fmt.Println("✅ Task deleted.")
}

func handleDeleteAll(tasks *[]Task) {
	fmt.Println("⚠️ Are you sure you want to DELETE ALL tasks? [y/n]:")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input != "y" && input != "n" {
		fmt.Println("❌ Invalid input. Only 'y' or 'n' are accepted.")
		return
	}

	if input == "y" {
		*tasks = []Task{}
		err := saveTasks(tasks)
		if err != nil {
			fmt.Println("❌ Failed to delete all tasks:", err)
			return
		}
		fmt.Println("✅ All tasks deleted.")
	} else {
		fmt.Println("❌ Delete All Canceled.")
	}
}


func handleDone(tasks *[]Task){
	fmt.Println("Enter the ID of the task to mark as done:")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	id, err := strconv.Atoi(input)
	if err != nil || id <= 0 {
		fmt.Println("⚠️ Invalid ID")
		return
	}
	index :=  id-1

	if (*tasks)[index].Status == "done" {
		fmt.Println("⚠️ Task already marked as done.")
		return
	}
	(*tasks)[index].Status = "done"
	saveTasks(tasks)
	fmt.Printf("✅ Task %d marked as done.\n", id)
}

func handleExit(){
	fmt.Println("👋 Goodbye!")
	os.Exit(0)
}


	


