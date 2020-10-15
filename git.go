package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/ddliu/go-httpclient"
	"github.com/iancoleman/strcase"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

type Task struct {
	Name        string
	Id          float64
	Description string
	Assignees   []User
	Status      float64
}

type User struct {
	Name  string
	Email string
	Id    string
	Oid   string
}

func (t Task) ListAssignees() string {

	var names []string
	for i := range t.Assignees {
		names = append(names, t.Assignees[i].Name)
	}

	return fmt.Sprintf(strings.Join(names, ","))
}

var GitCommand = &cli.Command{
	Name:  "checkout",
	Usage: "checkout a new branch",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "projectId",
			Aliases:  []string{"p"},
			Usage:    "Project id to get the tasks for",
			Required: true,
		},
		&cli.StringFlag{
			Name: "from",
			Aliases: []string{
				"branch",
				"b",
				"f",
			},
			Usage: "branch name to checkout from",
			Value: "master",
		},
	},
	Action: func(c *cli.Context) error {
		conf, err := LoadConfig()
		if err != nil {
			return err
		}

		httpclient.Defaults(httpclient.Map{
			"Authorization": fmt.Sprintf("%s %s", conf.TokenType, conf.AccessToken),
		})

		projRes, err := httpclient.Get(fmt.Sprintf("https://quire.io/api/task/list/id/%s", c.String("projectId")))
		defer projRes.Body.Close()
		tasksArr, err := ReadArrBody(projRes)
		if err != nil {
			return err
		}

		projUserRes, err := httpclient.Get(fmt.Sprintf("https://quire.io/api/user/list/project/id/%s", c.String("projectId")))
		defer projUserRes.Body.Close()
		usersArr, err := ReadArrBody(projUserRes)
		if err != nil {
			return err
		}

		users := make([]User, len(usersArr))
		for i := range usersArr {
			users[i] = User{
				Name:  usersArr[i]["name"].(string),
				Email: usersArr[i]["email"].(string),
				Id:    usersArr[i]["id"].(string),
				Oid:   usersArr[i]["oid"].(string),
			}
		}

		filterUsers := func(assignees []interface{}) []User {
			var res []User

			for _, u := range users {
				for _, a := range assignees {
					if u.Oid == a.(string) {
						res = append(res, u)
					}
				}
			}

			return res
		}

		var tasks []Task
		for i := range tasksArr {
			status := tasksArr[i]["status"].(float64)

			if status < 100 && status > 0 {
				desc := tasksArr[i]["description"].(string)

				if desc == "" {
					desc = "NO_DESCRIPTION_PROVIDED"
				}

				tasks = append(tasks, Task{
					Name:        tasksArr[i]["name"].(string),
					Id:          tasksArr[i]["id"].(float64),
					Description: desc,
					Assignees:   filterUsers(tasksArr[i]["assignees"].([]interface{})),
					Status:      status,
				})
			}
		}

		taskPrompt := &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "\U0001f538 ({{ .Id | red }}) {{ .Name | cyan }} ({{ .Status | cyan }}%)",
			Inactive: "  ({{ .Id | red }}) {{ .Name | white }} ({{ .Status | white }}%)",
			Selected: "\U0001f538 {{ .Name | red | cyan }}",
			Details: `
------------------- Task --------------------
{{ "Name:" }}	{{ .Name }}
{{ "Id:" }}	{{ .Id }}
{{ "Status:" }}	{{ .Status }}%
{{ "Assigned to:" }}	{{ .ListAssignees }}`,
		}

		prompt := promptui.Select{
			Label:   "Select a task",
			Items:   tasks,
			Pointer: promptui.DefaultCursor,
			Searcher: func(input string, index int) bool {
				task := tasks[index]
				name := strings.Replace(strings.ToLower(task.Name), " ", "", -1)
				id := fmt.Sprintf("%.0f", task.Id)
				ass := strings.Replace(strings.ToLower(tasks[index].ListAssignees()), ",", " ", -1)
				input = strings.Replace(strings.ToLower(input), " ", "", -1)

				return strings.Contains(name, input) || strings.Contains(id, input) || strings.Contains(ass, input)
			},
			Templates:         taskPrompt,
			StartInSearchMode: true,
		}

		idx, _, err := prompt.Run()
		task := tasks[idx]

		branchName := fmt.Sprintf("%s/%.0f-%s", strcase.ToKebab(task.ListAssignees()), task.Id, strcase.ToKebab(task.Name))

		log.Printf("creating branch %q based on %q", branchName, c.String("from"))

		err = exec.Command("git", "checkout", "-b", branchName, c.String("from")).Run()
		if err != nil {
			log.Fatalln(err)
			return err
		}

		log.Println("updating task status")
		_, err = httpclient.PutJson(fmt.Sprintf("https://quire.io/api/task/id/%s/%.0f", c.String("projectId"), task.Id), `{"status": 15}`)
		if err != nil {
			log.Fatalln(err)
			return err
		}

		log.Println("task successfully updated")
		return nil
	},
}
