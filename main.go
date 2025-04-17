package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)


type Task struct {
	Id     int    `json:"id"`
	Item   string `json:"item"`
	Status string `json:"status"` // "pending" or "done"
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

func saveTask(tasks *[]Task)error{
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
		fmt.Println("\nAvailable commands: add | list | delete | done | help | exit")
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
	
		switch input {
			case "add":{
				handleAdd(&tasks)
			}
			case "list": {
				handleList(&tasks)
			}
			case "delete": {
				handleDelete(&tasks)
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

	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		fmt.Println("⚠️ Task cannot be empty")
	}

	newTask := Task{
		Id: len(*tasks) + 1,
		Item: input,
		Status: "pending",
	}

	*tasks = append(*tasks, newTask)

	fmt.Printf("Saving Task: %s", newTask.Item)

	saveTask(tasks)
	fmt.Printf("✅ Task added: \"%s\"\n", newTask.Item)

}

func handleList(tasks *[]Task) {
	if len(*tasks) == 0 {
		fmt.Println("No item in the list.")
	}
	for _, task := range *tasks {
		fmt.Printf("item: %d - %s - status: %s\n", task.Id,task.Item, task.Status)
	}
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
	saveTask(tasks)
	fmt.Println("✅ Task deleted.")
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
	saveTask(tasks)
	fmt.Printf("✅ Task %d marked as done.\n", id)
}

func handleExit(){
	fmt.Println("👋 Goodbye!")
	os.Exit(0)
}

	


