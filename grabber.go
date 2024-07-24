package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	//время начала программы
	start := time.Now()
	var src, dst string
	src, dst, err := flagParsing()
	if err != nil {
		fmt.Printf("%v \r\n", err)
		return
	}
	err = checkOrCreateDir(dst)
	if err != nil {
		fmt.Printf("%v \r\n", err)
		return
	}
	urls, err := readLinesFromFile(src)
	if err != nil {
		fmt.Printf("%v \r\n", err)
		return
	}
	parseMakeHTML(urls, dst)
	//время завершения программы
	finish := time.Since(start).Truncate(10 * time.Millisecond).String()
	fmt.Println("Время выполнения программы:", finish)
}

// flagParsing - обработка флагов
func flagParsing() (string, string, error) {

	//флаг файла
	src := flag.String("src", "", "используйте флаг -stc для введения файла с URL.")
	//флаг папки

	dst := flag.String("dst", "", "используйте флаг -dst для введения каталога для сохраниня html.")
	flag.Parse()
	//проверка наличия флагов
	if len(*src) == 0 {
		flag.PrintDefaults()
		return "", "", fmt.Errorf("отстутствуют необходимые флаги: -src")
	}
	if len(*dst) == 0 {
		flag.PrintDefaults()
		return "", "", fmt.Errorf("отстутствуют необходимые флаги: -dst")
	}
	return *src, *dst, nil
}

// checkOrCreateDir - проверка существования директории и её создание в случае отсутствия
func checkOrCreateDir(path string) error {
	//проверка существования каталога
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("ошибка при создании каталога %s: %v", path, err)
	}
	if os.IsNotExist(err) {
		//создание каталога
		os.Mkdir(path, os.ModeDir|0755)
	}
	return nil
}

// readLinesFromFile - построчное чтение url из файла
func readLinesFromFile(src string) ([]string, error) {
	//открываем файл
	file, err := os.Open(src)
	if err != nil {
		return []string{}, fmt.Errorf("ошибка при открытии файла с url: %v", err)
	}
	scanner := bufio.NewScanner(file)
	//создаём срез url-в
	var urls []string
	for scanner.Scan() {
		//добавляем в срез текущую строку файла
		urls = append(urls, formatURL(scanner.Text()))
	}
	defer file.Close()
	return urls, nil
}

// parseMakeHTML - парсинг по url из среза и создание html-файла в случае успеха
func parseMakeHTML(urls []string, dst string) {
	for i := 0; i < len(urls); i++ {
		//попытка парсинга по текущему элементу среза
		resp, err := parse(urls[i])
		if err != nil {
			fmt.Printf("ошибка при парсинге с url %s: %v \r\n", urls[i], err)
		} else {
			//попытка создания html-файла для текущего элемента среза
			err = createHTML(resp, dst, urls[i])
			if err != nil {
				fmt.Printf("ошибка при создании HTML-файла страницы %s: %v\r\n", urls[i], err)
			}
		}

	}
}

// formatURL - проверка наличия "http://" в начале строки
func formatURL(urlWithoutPrefix string) string {
	var url string
	if !(strings.HasPrefix(urlWithoutPrefix, "http://")) || !(strings.HasPrefix(urlWithoutPrefix, "https://")) {
		//приведение url  к нужному формату
		url = fmt.Sprintf("http://%s", urlWithoutPrefix)
	} else {
		url = urlWithoutPrefix
	}
	return url
}

// parse - get-запрос и возврат тела ответа
func parse(url string) ([]byte, error) {

	//отправка get запроса
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии url %s: %v", url, err)
	}
	//чтение тела ответа
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return b, fmt.Errorf("ошибка при получении тела запроса %s: %v", url, err)
	}
	defer resp.Body.Close()
	return []byte(b), err
}

// createHTML - создание HTML на основе полученного ответа
func createHTML(resp []byte, dst, url string) error {
	nameHTML := fmt.Sprintf("%s%s.html", dst, strings.Replace(url, "/", "|", -1))
	//создание html-файла
	file, err := os.Create(nameHTML)
	if err != nil {
		return fmt.Errorf("ошибка при cоздании файла %s: %v", nameHTML, err)
	}
	//запись ответа на запрос в файл
	file.Write([]byte(resp))
	fmt.Printf("Страница %s успешно сохранена \r\n", url)
	defer file.Close()
	return nil
}
