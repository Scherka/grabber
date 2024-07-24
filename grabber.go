package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	//время начала программы
	start := time.Now()
	FlagParsing()
	//время завершения программы
	finish := time.Since(start)
	fmt.Println("Время выполнения программы:", finish)
}

// FlagParsing - обработка флагов
func FlagParsing() {
	//флаг файла
	srcfl := flag.String("src", "url.txt", "third")
	//флаг папки
	dirfl := flag.String("dst", "./dst/", "third")
	flag.Parse()
	fmt.Println("Открываем", *srcfl)
	//проверяем наличие каталога, создаём, если его нет
	CheckDir(*dirfl)
	//считывание содержимого файла построчно
	if !(strings.HasSuffix(*dirfl, "/")) {
		//добавление "/" к концу введённого каталога, если его нет
		ReadLines(*srcfl, fmt.Sprintf("%s/", *dirfl))
	} else {
		ReadLines(*srcfl, *dirfl)
	}
}

// CheckDir - проверка существования директории и её создание в случае отсутствия
func CheckDir(path string) {
	//проверка существования каталога
	if _, err := os.Stat(path); os.IsNotExist(err) {
		//создание каталога
		os.Mkdir(path, os.ModeDir|0755)
	}
}

// ReadLines - построчное чтение файла с url
func ReadLines(src string, dst string) {
	//создаём группу ожидания
	wg := sync.WaitGroup{}
	//открываем файл
	file, err := os.Open(src)
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(file)
	//проходим все строки документа
	for scanner.Scan() {
		//Увеличиваем размер группы
		wg.Add(1)
		//считываем строку из файла
		scan := scanner.Text()
		go func() {

			Parse(scan, dst)
			//уменьшаем размер группы
			defer wg.Done()
		}()

	}
	//ждём завершения всех горутин
	wg.Wait()
}

// Parse - get-запрос и запись ответа в .html-файл
func Parse(urlWithoutPrefix string, dst string) {
	var url string
	//проверка наличия "http:// в начале строки"
	if !(strings.HasPrefix(urlWithoutPrefix, "http://")) || !(strings.HasPrefix(urlWithoutPrefix, "https://")) {
		//приведение url  к нужному формату
		url = fmt.Sprintf("http://%s", urlWithoutPrefix)
	} else {
		url = urlWithoutPrefix
	}
	//отправка get запроса
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Не удалось открыть '%s' \r\n", urlWithoutPrefix)
	} else {
		//создание html-файла
		file, err := os.Create(fmt.Sprintf("%s%s.html", dst, strings.Replace(urlWithoutPrefix, "/", "|", -1)))
		if err != nil {
			fmt.Println("Ошибка при создании файла")
		} else {
			//запись ответа на запрос в файл
			resp.Write(file)
			fmt.Printf("Страница %s успешно сохранена \r\n", urlWithoutPrefix)
		}
		file.Close()
	}
}
