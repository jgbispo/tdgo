package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jgbispo/todo-app-go"
)

const (
	todoFile = ".todos.json"
)

func main() {
	add := flag.Bool("add", false, "add a new todo")
	done := flag.Int("done", 0, "mark a todo as completed")
	del := flag.Int("del", 0, "delete a todo")
	list := flag.Bool("list", false, "list all")
	work := flag.Int("work", 0, "Start current task")

	flag.Parse()

	todos := &todo.Todos{}

	if err := todos.Load(todoFile); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	switch {
	case *add:
		task, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}

		todos.Add(task)

		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		todos.Print()
	case *done > 0:
		err := todos.Complete(*done)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		todos.Print()
	case *del > 0:
		err := todos.Delete(*del)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		todos.Print()
	case *list:
		todos.Print()
	case *work > 0:
		err := todos.Working(*work)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		todos.Print()
	default:
		fmt.Fprintln(os.Stdout, "Invalid command")
		os.Exit(0)
	}
}

func getInput(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Scan()

	if err := scanner.Err(); err != nil {
		return "", err
	}

	text := scanner.Text()

	if len(text) == 0 {
		return "", errors.New("empty todo is not allowed")
	}

	return scanner.Text(), nil
}
