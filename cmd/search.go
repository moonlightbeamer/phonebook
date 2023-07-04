/*
Copyright © 2021 Mihalis Tsoukalos <mihalistsoukalos@gmail.com>
*/

package cmd

import (
	  "fmt"
	_ "regexp"
	  "strings"
      "strconv"
	  "time"
	  "github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search for the number",
	Long: `search whether a partial telephone number exists in the
	phone book application or not.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get key
		searchKey, _ := cmd.Flags().GetString("tel")
		if searchKey == "" {
			fmt.Println("Not a valid partial telephone number:", searchKey)
			return
		}
		t := strings.ReplaceAll(searchKey, "-", "")

		if !matchTel(t) {
			fmt.Println("Not a valid telephone number:", t)
			return
		}

		// Search for it
		search_err := search(t, &data, &index, &data_match)
		if search_err != nil {
			fmt.Println("Number not found, or update data file error:", searchKey)
			return
		}
		for _, v := range data_match {
			fmt.Println(data[v])
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringP("tel", "t", "", "Telephone number to search")
}

func search(key string, d *PhoneBook, ind *IndexBook, re *[]int) error {
	*re = []int{}
    for _, v := range *ind {
		if strings.Count(v.Tel, key) > 0 {
			*re = append(*re, v.Index_num)
		}
	}
    if len(*re) == 0 {
        return fmt.Errorf("%s not found.", key)
    } else {
		for _, re_v := range *re {
			(*d)[re_v].LastAccess = strconv.FormatInt(time.Now().Unix(), 10)
		}
    	save_err := saveJsonDataFile(JSONFILE, d)
		if save_err != nil {
			return save_err
		} else {
    		return nil
		}
  	}
}

/*
func matchTel(s string) bool {
	t := []byte(s)
	re := regexp.MustCompile(`\d+$`)
	return re.Match(t)
}
*/
