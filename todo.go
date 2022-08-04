package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

type item struct {
	Task        string
	Done        bool
	Working     bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type Todos []item

func (t *Todos) Add(task string) {
	todo := item{
		Task:        task,
		Done:        false,
		Working:     false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Now(),
	}

	*t = append(*t, todo)
}

func (t *Todos) Complete(index int) error {
	ls := *t
	fmt.Println(index)
	if index <= 0 || index > len(ls) {
		return errors.New("Invalid index\n")
	}

	ls[index-1].CompletedAt = time.Now()
	ls[index-1].Done = true

	return nil
}

func (t *Todos) Working(index int) error {
	ls := *t

	if index <= 0 || index > len(ls) {
		return errors.New("Invalid index\n")
	}

	ls[index-1].Working = !ls[index-1].Working

	return nil
}

func (t *Todos) Delete(index int) error {
	ls := *t

	if index <= 0 || index > len(ls) {
		return errors.New("Invalid index")
	}

	*t = append(ls[:index-1], ls[index:]...)

	return nil
}

func (t *Todos) Load(filename string) error {
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return err
	}

	err = json.Unmarshal(file, t)

	if err != nil {
		return err
	}

	return nil
}

func (t *Todos) Store(filename string) error {

	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 8644)
}

func (t *Todos) Print() {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Done?"},
			{Align: simpletable.AlignRight, Text: "CreatedAt"},
			{Align: simpletable.AlignRight, Text: "CompletedAt"},
		},
	}

	var cells [][]*simpletable.Cell

	for idx, item := range *t {
		idx++

		task := purple(item.Task)
		id := purple(fmt.Sprintf("%d", idx))
		done := purple(fmt.Sprintf("%t", item.Done))
		createdAt := purple(item.CreatedAt.Format(time.RFC822))
		completedAt := purple(item.CompletedAt.Format(time.RFC822))

		if item.Done {
			task = green(fmt.Sprintf("[Done] %s", item.Task))
			id = green(fmt.Sprintf("%d", idx))
			done = green(fmt.Sprintf("%t", item.Done))
			createdAt = green(item.CreatedAt.Format(time.RFC822))
			completedAt = green(item.CreatedAt.Format(time.RFC822))
		}

		if item.Working && !item.Done {
			task = blue(fmt.Sprintf("[Working] %s", item.Task))
			id = blue(fmt.Sprintf("%d", idx))
			done = blue(fmt.Sprintf("%t", item.Done))
			createdAt = blue(item.CreatedAt.Format(time.RFC822))
			completedAt = blue(item.CreatedAt.Format(time.RFC822))
		}

		cells = append(cells, *&[]*simpletable.Cell{
			{Text: id},
			{Text: task},
			{Text: done},
			{Text: createdAt},
			{Text: completedAt},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}

	pending := red(fmt.Sprintf("You have %d pending todos", t.CountPending()))
	if t.CountPending() == 0 {
		pending = green(fmt.Sprintf("You have %d pending todos", t.CountPending()))

	}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter,
			Span: 5,
			Text: pending,
		},
	}}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}

func (t *Todos) CountPending() int {
	total := 0

	for _, item := range *t {
		if !item.Done {
			total++
		}
	}

	return total
}
