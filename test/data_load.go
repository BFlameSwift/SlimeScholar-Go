package main

import (
	"bufio"
	"fmt"
	"os"
)

func printMagPaper(file_path string) {
	open, err := os.Open(file_path)
	if err != nil {
		fmt.Println("打开失败")
	}
	scanner := bufio.NewScanner(open)
	i := 0
	for scanner.Scan() {
		if i < 10 {
			fmt.Println(scanner.Text())
		}
		i++
	}
	fmt.Println("line sum", i)

}
func printAminerPaper(file_path string) {
	open, err := os.Open(file_path)
	if err != nil {
		fmt.Println("打开失败")
	}
	scanner := bufio.NewScanner(open)
	i := 0
	for scanner.Scan() {
		if i < 10 {
			fmt.Println(scanner.Text())
		}
		i++
	}
	fmt.Println("line sum", i)

}

func main() {
	printMagPaper("E:\\Paper\\mag_papers_0\\mag_papers_1.txt")
	printAminerPaper("E:\\Paper\\aminer_papers_0\\aminer_papers_1.txt")
}
