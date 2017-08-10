package main

import (
	"os"
	"github.com/urfave/cli"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"log"
	"reflect"
	"go/types"
)

var (
	counter int
	mapper Keyvalue
	array_diff = make(map[string][]Keyvalue)
)

type Keyvalue map[string]interface{}
type Iter chan interface{}


func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func recursion(original Keyvalue, modified Keyvalue, path []string) {


	kListModified := ListStripper(modified)
	kListOriginal := ListStripper(original)

	if len(kListModified) > 1 || len(kListOriginal) > 1 {
		proc := true
		for k, v := range original {
			if IndexOf(kListModified, k) == -1 {
				array_diff["Removed"] = append(array_diff["Removed"],Keyvalue{"Path": path, "Key": k, "Value":v})
				proc = false
			}
		}

		for k, v := range modified {
			if IndexOf(kListOriginal, k) == -1 {
				array_diff["Added"] = append(array_diff["Added"],Keyvalue{"Path": path, "Key": k, "Value":v})
				proc = false
			}
		}
		if proc {
			for k := range original {
				recursion(Keyvalue{k:original[k]},Keyvalue{k:modified[k]},path)
			}
		}
		return
	}
	for k := range original {
		counter++

		if reflect.TypeOf(original) == types.String {

		}


		/*
		switch v := interface{}(original).(type) {
		case string:
			valOrig = v
		default:
			valOrig = v[k]
		}
		switch v := interface{}(modified).(type) {
 	  	case string:
			valMod = v
		default:
			valMod = v[k]
		}
		if valOrig != valMod {
			fmt.Println("help")
		}
		*/


	}

}


func main() {
	var patch, object, original_obj, modified_obj string

	app := cli.NewApp()
	app.Name = "JsonDiffer"
	app.Version = "0.1"
	app.Usage = "Used to get an object-based difference between two json objects."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "test, t",
			Usage: "just taking up space",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "diff",
			Aliases: []string{"d"},
			Usage:   "Diff json objects",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "origin, o",
					Usage: "Original `OBJECT` to compare against",
					Value: "",
					Destination: &original_obj,
				},
				cli.StringFlag{
					Name: "modified, m",
					Usage: "Modified `OBJECT` to compare against",
					Value: "",
					Destination: &modified_obj,
				},
				cli.StringFlag{
					Name: "output, O",
					Usage: "File output location",
					Value: "",
					Destination: &modified_obj,
				},
			},
			Action:  func(c *cli.Context) error {
				var json_original, json_modified Keyvalue
				var path []string
				if original_obj == "" {
					fmt.Print("ORIGIN is required!\n\n")
					cli.ShowCommandHelp(c, "diff")
					os.Exit(1)
				}
				if modified_obj == "" {
					fmt.Print("MODIFIED is required!\n\n")
					cli.ShowCommandHelp(c, "diff")
					os.Exit(1)
				}

				/* TODO WE WANT TO DO ALL OUR INIT STUFF IN THIS AREA */

				/*
				array_diff["Changed"] = []Keyvalue{}
				array_diff["Added"] = []Keyvalue{}
				array_diff["Removed"] = []Keyvalue{}
				*/

				read,err := ioutil.ReadFile(original_obj)
				check(err)
				_ = json.Unmarshal([]byte(read), &json_original)

				read,err = ioutil.ReadFile(modified_obj)
				check(err)
				_ = json.Unmarshal([]byte(read), &json_modified)

				recursion(json_original, json_modified, path)



				output,_ := json.Marshal(array_diff)
				os.Stdout.Write(output)

				return nil
			},
		},
		{
			Name: "patch",
			Aliases: []string{"p"},
			Usage:	"Apply patch file to json object",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "patch, p",
					Usage: "`PATCH` the OBJECT",
					Value: "",
					Destination: &patch,
				},
				cli.StringFlag{
					Name: "object, o",
					Usage: "`OBJECT` to PATCH",
					Value: "",
					Destination: &object,
				},
			},
			Action: func(c *cli.Context) error {

				return nil
			},
		},
	}

	app.Run(os.Args)

}


func ListStripper(input Keyvalue ) []string {
	var r []string
	for key := range input {
		r = append(r, key)
	}
	return r
}

func IndexOf(inputList []string, inputKey string) int {
	for i, v := range inputList {
		if v == inputKey {
			return i
		}
	}
	return -1
}