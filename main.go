package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/TwiN/go-color"
)

type Comments struct {
	Comments []Comment `json:"comments"`
}

type Comment struct {
	File string `json:"file"`
	Text string `json:"text"`
}

type Output struct {
	path    string
	comment string
	size    int64
}

var commentAlign int

//var ctreeoutput = []string{}
var ctreeoutput = []Output{}

func printSpaceX(level int) string {
	r := ""
	for i := 0; i < level; i++ {
		r = r + "    "
	}
	c := [5]string{color.White, color.Cyan, color.Yellow, color.Red, color.Blue}
	return r + c[level%5] + "├── " + color.Reset
}

func listDirectory(path string, depth int, level int) {
	if depth < 0 {
		return
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	comments := readComments(path)

	for _, f := range files {
		if f.Name() == ".CTREE_Store" {
			continue
		}

		output := printSpaceX(level)

		var c string
		if val, ok := comments[f.Name()]; ok {
			c = color.Yellow + "\t#" + val + color.Reset
		}

		if f.IsDir() {
			//println(color.Bold + f.Name() + c + color.Reset)
			output = output + color.Bold + f.Name() + color.Reset
			if depth > 0 {
				listDirectory(path+f.Name()+"/", depth-1, level+1)
			}
		} else {
			//println(color.Colorize(color.Green, f.Name()+c))
			output = output + color.Colorize(color.Green, f.Name())
		}
		//println(output)

		if len(output) > commentAlign {
			commentAlign = len(output)
		}

		ctreeoutput = append(ctreeoutput, Output{path: output, comment: c, size: f.Size()})
	}
}

func printCTree() {
	for i := len(ctreeoutput) - 1; i >= 0; i-- {
		space := commentAlign - len(ctreeoutput[i].path)
		r := ""
		for i := 0; i < space; i++ {
			r = r + " "
		}
		//fmt.Sprintf("%d", ctreeoutput[i].size)
		println(ctreeoutput[i].path + r + ctreeoutput[i].comment)
	}
}

func readComments(path string) map[string]string {
	file, _ := ioutil.ReadFile(path + ".CTREE_Store")
	data := Comments{}
	a := make(map[string]string)

	_ = json.Unmarshal([]byte(file), &data)
	for i := 0; i < len(data.Comments); i++ {
		a[data.Comments[i].File] = data.Comments[i].Text
	}
	return a
}

func writeComments(filepath string, comment string) {
	_, err := os.Stat(filepath)
	if err != nil {
		log.Fatal(err)
	}
	comments := readComments(path.Dir(filepath) + "/")

	commentsToSave := []Comment{}
	saveCommets := Comments{commentsToSave}

	for k, v := range comments {
		if k != path.Base(filepath) {
			saveCommets.Comments = append(saveCommets.Comments, Comment{File: k, Text: v})
		}
	}

	msgtext := "REMOVED"

	if comment != "" {
		saveCommets.Comments = append(saveCommets.Comments, Comment{File: path.Base(filepath), Text: comment})
		msgtext = "created"
	}

	file, _ := json.MarshalIndent(saveCommets, "", " ")
	_ = ioutil.WriteFile(path.Dir(filepath)+"/.CTREE_Store", file, 0644)

	println(color.Colorize(color.Yellow, "Comment "+msgtext+" for: "+filepath+" in "+path.Dir(filepath)+"/.CTREE_Store"))
}

func version() {
	println(color.Colorize(color.Cyan, `		    
   ____    _____      ____     U _____ u U _____ u 
U /"___|  |_ " _|  U |  _"\ u  \| ___"|/ \| ___"|/ 
\| | u      | |     \| |_) |/   |  _|"    |  _|"   
 | |/__    /| |\     |  _ <     | |___    | |___   
  \____|  u |_|U     |_| \_\    |_____|   |_____|  
 _// \\   _// \\_    //   \\_   <<   >>   <<   >>  
(__)(__) (__) (__)  (__)  (__) (__) (__) (__) (__)  ctreeVTAG
			`))
}

func main() {

	goinitV := flag.Bool("v", false, "ctree version")
	giDebug := flag.Bool("d", false, "ctree debug (verbose)")
	giDepth := flag.Int("depth", -1, "tree depth")
	giComment := flag.String("c", "#", "file / folder Comment. Ex: ctree -c \"some comment\" path/to/file")

	flag.Parse()

	if *goinitV {
		version()
		return
	}

	path := "."

	if flag.NArg() != 0 {
		_, err := os.Stat(flag.Arg(0))
		if !os.IsNotExist(err) {
			path = flag.Arg(0)
			if path[len(path)-1:] == "/" {
				path = strings.TrimRight(path, "/")
			}
		}
	}

	if *giDebug {
		println(color.Colorize(color.Cyan, "ctree debug mode"))
	}

	if *giComment != "#" {
		if path != "." {
			writeComments(path, *giComment)
		} else {
			flag.Usage()
		}
		return
	}

	if *giDepth == -1 {
		*giDepth = 5
	}

	listDirectory(path+"/", *giDepth-1, 0)
	printCTree()
}
