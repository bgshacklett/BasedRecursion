package main

import (
	"os"
	"github.com/urfave/cli"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"log"
	"reflect"
)

var (
	counter int
	mapper map[int]interface{}
	array_diff = make(map[string][]Keyvalue)
)

type Keyvalue map[string]interface{}
type Iter chan interface{}


func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
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

func recursion(original Keyvalue, modified Keyvalue, path []string) {

	kListModified := ListStripper(modified)
	kListOriginal := ListStripper(original)

	fmt.Println(":::SHIT AT START OF FUNCTION:::")
	fmt.Println(kListModified)
	fmt.Println(kListOriginal)
	fmt.Println(path)
	fmt.Println(":::::::::::::::::")
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
		var valOrig, valMod interface{}
		var back_orig, back_mod Keyvalue
		counter++

		if reflect.TypeOf(original).Name() == "string" {

			valOrig = original
		} else {
			valOrig = original[k]
		}
		if reflect.TypeOf(modified).Name() == "string" {
			valMod = modified
		} else {
			valMod = modified[k]
		}
		fmt.Println(":::COMPARING INPUTS:::")
		fmt.Println(k)
		fmt.Println(valOrig)
		fmt.Println(valMod)
		fmt.Println("::::::::::::::::::::::::::")

		if !(reflect.DeepEqual(valMod, valOrig)) {
			fmt.Println(":::TYPE OF VALUE THINGY:::")
			fmt.Println(reflect.TypeOf(valOrig).Name())
			fmt.Println("::::::::::::::::::::::::::")

			if reflect.TypeOf(valOrig).Name() != "" {
				fmt.Println("SHOULD BE ADDING SHIT \\\\\\\\\\\\\\\\\\\\")
				array_diff["Changed"] = append(array_diff["Changed"],Keyvalue{"Path": path, "Key": k, "oldValue":valOrig,"newValue":valMod})
				return
			} else {
				npath := append(path, k)

				fmt.Println(valOrig)
				fmt.Println(valMod)

				orig_out,_ := json.Marshal(valOrig)
				_ = json.Unmarshal([]byte(orig_out), &back_orig)
				mod_out,_ := json.Marshal(valOrig)
				_ = json.Unmarshal([]byte(mod_out), &back_mod)
				//fmt.Println(back_orig)
				go recursion(back_orig,back_mod, npath)
				return
			}
		}
		return

	}
	return
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


				if reflect.DeepEqual(json_original, json_modified) {
					fmt.Println("No differences!")
					os.Exit(0)
				} else {
					recursion(json_original, json_modified, path)
				}



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


