package main

import (
	"bufio"
	"flag"
	"fmt"
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
		fmt.Println(err)
		return
	}
	err = checkOrCreateDir(dst)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = readLinesFromFileRunParseMakeHTML(src, dst)
	if err != nil {
		fmt.Println(err)
		return
	}
	//время завершения программы
	finish := time.Since(start).Truncate(10 * time.Millisecond).String()
	fmt.Println("Время выполнения программы:", finish)
}

// FlagParsing - обработка флагов
func flagParsing() (string, string, error) {

	//флаг файла
	src := flag.String("src", "", "Используйте флаг -stc для введения файла с URL.")
	//флаг папки
	dst := flag.String("dst", "", "Используйте флаг -dst для введения каталога для сохраниня html.")
	flag.Parse()
	if len(*src) == 0 {
		return *src, *dst, fmt.Errorf("ошибка: используйте флаг -src для введения файла с URL")
	}
	if len(*dst) == 0 {
		return *src, *dst, fmt.Errorf("ошибка: используйте флаг -dst для введения папки, куда будут сохраняться HTML")
	}
	return *src, *dst, nil
}

// CheckDir - проверка существования директории и её создание в случае отсутствия
func checkOrCreateDir(path string) error {
	//проверка существования каталога
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		//fmt.Println(err)
		return fmt.Errorf("ошибка при создании каталога %s: %v", path, err)
	}
	if os.IsNotExist(err) {
		//создание каталога
		os.Mkdir(path, os.ModeDir|0755)
	}
	return nil
}

// readLinesFromFileRunParseMakeHTML - построчное чтение url из файла
func readLinesFromFileRunParseMakeHTML(src string, dst string) error {
	//создаём группу ожидания
	//открываем файл
	file, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("ошибка при открытии файла с url: %v", err)
	}
	scanner := bufio.NewScanner(file)
	//проходим все строки документа
	for scanner.Scan() {
		scan := formatURL(scanner.Text())
		func() {
			resp, err := parse(scan)
			if err != nil {
				fmt.Println(err)

			} else {
				err = createHTML(*resp, dst, scan)
				if err != nil {
					fmt.Println(err)
				}
			}
		}()
	}
	return nil
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

// Parse - get-запрос и возврат ответа-
func parse(url string) (*http.Response, error) {

	//отправка get запроса
	resp, err := http.Get(url)
	if err != nil {
		//fmt.Printf("Не удалось открыть '%s' \r\n", url)
		return nil, fmt.Errorf("ошибка при открытии url %s: %v", url, err)
	}

	return resp, err
}

// createHTML - создание HTML на основе полученного ответа
func createHTML(resp http.Response, dst, url string) error {
	nameHTML := fmt.Sprintf("%s%s.html", dst, strings.Replace(url, "/", "|", -1))
	//создание html-файла
	file, err := os.Create(nameHTML)
	if err != nil {
		return fmt.Errorf("ошибка при cоздании файла %s: %v", nameHTML, err)
	}
	//запись ответа на запрос в файл
	resp.Write(file)
	fmt.Printf("Страница %s успешно сохранена \r\n", url)
	defer file.Close()
	return nil
}
