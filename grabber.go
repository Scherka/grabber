package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var dst = ""    //переменная каталога, в которую будут записаны файлы
var counter = 1 //счётчик успешно записанных файлов

func main() {
	start := time.Now()                             //время начал программы
	srcfl := flag.String("src", "url.txt", "third") //флаг файла
	dirfl := flag.String("dst", "`/dst/", "third")  //флаг папки
	flag.Parse()
	dst = *dirfl
	CheckDir(dst)                       //проверяем наличие каталога, создаём, если его нет
	if !(strings.HasSuffix(dst, "/")) { //добавление "/" к концу введённого каталога, если его нет
		dst += "/"
	}

	fmt.Println("Открываем", *srcfl)
	src := *srcfl
	ReadLines(src)              //считывание содержимого файла построчно
	finish := time.Since(start) //время завершения программы
	fmt.Println("Время выполнения программы:", finish)
}
func CheckDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) { //проверка существования каталога
		os.Mkdir(path, os.ModeDir|0755) //создание каталога
	}
}

func ReadLines(src string) {
	file, err := os.Open(src) //открываем файл
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { //проходим все строки документа
		line := scanner.Text()
		Parse(line) //
	}
}
func Parse(urlog string) {
	var url string
	if !(strings.HasPrefix(urlog, "http://")) || !(strings.HasPrefix(urlog, "https://")) { //проверка наличия "http:// в начале строки"
		url = "http://" + urlog //приведение url  к нужному формату
	} else {
		url = urlog
	}
	resp, err := http.Get(url) //отправка get запроса
	if err != nil {
		fmt.Println("Не удалось открыть '" + urlog + "'")
	} else {
		//t := time.Now()
		file, err := os.Create(dst + "file" + strconv.Itoa(counter) + ".html") //создание html-файла
		if err != nil {
			fmt.Println("Ошибка при создании файла")
		}
		resp.Write(file) //запись ответа на запрос в файл
		fmt.Println("Страница " + urlog + " успешно сохранена")
		file.Close()
		counter++
	}
}
