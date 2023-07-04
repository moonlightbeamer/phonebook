/*
Copyright © 2023 Lin Zhu <moonlightbeamer@ymail.com>
*/
package cmd

import (
	  "encoding/json"
	  "fmt"
	  "io"
	  "math/rand"
	  "os"
	_ "path/filepath"
	  "regexp"
	_ "sort"
	  "strconv"
	  "strings"
	  "time"

	  "github.com/spf13/cobra"
)

// JSONFILE resides in the current directory
var JSONFILE = "./phonebook.json"
var INDFILE = "./phonebook_index.json"

type Entry struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Tel        string `json:"tel"`
	LastAccess string `json:"lastaccess"`
}

type PhoneBook []Entry

var data = PhoneBook{}

type Index struct {
	Tel        string `json:"tel"`
	Index_num  int    `json:"index_num"`
}

type IndexBook []Index

var index = IndexBook{}

var data_match = []int{}

// Implement sort.Interface
func (a PhoneBook) Len() int {
	return len(a)
}

// First based on surname. If they have the same
// surname take into account the name.
func (a PhoneBook) Less(i, j int) bool {
	if a[i].Surname == a[j].Surname {
		if a[i].Name == a[j].Name {
      return a[i].Tel < a[j].Tel
    }
    return a[i].Name < a[j].Name
	}
	return a[i].Surname < a[j].Surname
}

func (a PhoneBook) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// DeSerialize decodes a serialized slice with JSON records
func DeSerialize(slice_of_struct interface{}, json_file io.Reader) error {
	decoder := json.NewDecoder(json_file)
	return decoder.Decode(slice_of_struct)
}

// Serialize serializes a slice with JSON records
func Serialize(slice_of_struct interface{}, json_file io.Writer) error {
	encoder := json.NewEncoder(json_file)
  encoder.SetIndent("", "\t")
	return encoder.Encode(slice_of_struct)
}

func readJsonDataFile(filepath string, d *PhoneBook) error {
  // check filepath exists or not
  fileInfo, info_err := os.Stat(filepath)
  if info_err != nil {
    fmt.Println("Didn't find pre-existing \"", filepath, "\" data file, creating a new one.")
    init_data_file, create_err := os.Create(filepath)
		if create_err != nil {
			return create_err
		}
    l := 100
    init_data := PhoneBook{}
    for i := 0; i < l; i++ {
      name := getString(4)
      surname := getString(5)
      tel := strconv.Itoa(random(1999999999, 9999999999))
      time := strconv.FormatInt(time.Now().Unix(), 10)
      temp := Entry{
        Name:       name,
        Surname:    surname,
        Tel:        tel,
        LastAccess: time,
      }
      init_data = append(init_data, temp)
    }
    init_data_err := Serialize(init_data, init_data_file)
    if init_data_err != nil {
      return fmt.Errorf("initializing phonebook data file error: %s", init_data_err)
    }
    init_data_file.Close()
  } else if !(fileInfo.Mode().IsRegular()) {
		return fmt.Errorf("%s is not a regular file!", filepath)
	}
  
  // check filepath opens ok or not
  data_file, open_err := os.Open(filepath)
  if open_err != nil {
    return open_err
  }
  defer data_file.Close()

  json_deserialize_err := DeSerialize(&d, data_file)
  if json_deserialize_err != nil {
    return json_deserialize_err
  } else if len(*d) == 0 {
    return fmt.Errorf("%s is an empty file, no data found, exiting.", filepath)
  } else {
    return nil
  }
}

func saveJsonDataFile(filepath string, d *PhoneBook) error {
  // check filepath creates ok or not
	data_file, create_err := os.Create(filepath)
	if create_err != nil {
		return create_err
	}
	defer data_file.Close()

	json_serialize_err := Serialize(&d, data_file)
	if json_serialize_err != nil {
		return json_serialize_err
	}

	return nil
}

func readJsonIndexFile(filepath string, d *PhoneBook, ind *IndexBook) error {
  // check filepath exists or not
  fileInfo, info_err := os.Stat(filepath)
  if info_err != nil {
    fmt.Println("Didn't find pre-existing \"", filepath, "\" index file, creating a new one.")
    init_index_file, create_err := os.Create(filepath)
		if create_err != nil {
			return create_err
		}
    init_index := IndexBook{}
    for i, v := range *d {
      temp := Index{
        Tel: v.Tel, 
        Index_num: i,
      }
      init_index = append(init_index, temp)
    }
    init_index_err := Serialize(init_index, init_index_file)
    if init_index_err != nil {
      return fmt.Errorf("initializing phonebook index file error: %s", init_index_err)
    }
    init_index_file.Close()
  } else if !(fileInfo.Mode().IsRegular()) {
		return fmt.Errorf("%s is not a regular file!", filepath)
	}
  
  // check filepath opens ok or not
  index_file, open_err := os.Open(filepath)
  if open_err != nil {
    return open_err
  }
  defer index_file.Close()

  json_deserialize_err := DeSerialize(&ind, index_file)
  if json_deserialize_err != nil {
    return json_deserialize_err
  } else if len(*ind) == 0 {
    return fmt.Errorf("%s is an empty file, no index data found, exiting.", filepath)
  } else {
    return nil
  }
}

func saveJsonIndexFile(filepath string, d *PhoneBook, ind *IndexBook) error {
  // check filepath creates ok or not
	index_file, create_err := os.Create(filepath)
	if create_err != nil {
		return create_err
	}
	defer index_file.Close()
	json_serialize_err := Serialize(&ind, index_file)
	if json_serialize_err != nil {
		return json_serialize_err
	}
	return nil
}

// Initialized by the user – returns a pointer
// If it returns nil, there was an error
func initS(N, S, T string) *Entry {
	// Both of them should have a value
	if T == "" || S == "" {
		return nil
	}
	// Give LastAccess a value
	LastAccess := strconv.FormatInt(time.Now().Unix(), 10)
	return &Entry{Name: N, Surname: S, Tel: T, LastAccess: LastAccess}
}

func setJsonDataFILE() error {
	filepath := os.Getenv("PHONEBOOK")
	if filepath != "" {
		JSONFILE = filepath
	}
  return nil
}

func matchNameSur(s string) bool {
	return regexp.MustCompile("^[A-Z][a-z]*$").Match([]byte(s))
}

func matchTel(s string) bool {
	return regexp.MustCompile(`\d+$`).Match([]byte(s))
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func getString(l int64) string {
	startChar := "A"
	temp := ""
	var i int64 = 1
	for {
		myRand := random(0, 26)
		newChar := string(startChar[0] + byte(myRand))
		temp = temp + newChar
		if i == l {
			break
		}
		i++
	}
	return strings.Title(strings.ToLower(temp))
}
// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "phonebook",
	Short: "A phonebook app with commands",
	Long: `A phonebook app with commands built on cobra`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  setJsonDataFILE()

  read_data_err := readJsonDataFile(JSONFILE, &data)
  // io.EOF is fine because it means the file is empty
	if read_data_err != nil && read_data_err != io.EOF {
		return
	}

  read_index_err := readJsonIndexFile(INDFILE, &data, &index)
  // io.EOF is fine because it means the file is empty
	if read_index_err != nil && read_index_err != io.EOF {
		return
	}
  // can also be simplified as cobra.CheckErr(rootCmd.Execute())
	root_exec_err := rootCmd.Execute()
	if root_exec_err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.phonebook.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


