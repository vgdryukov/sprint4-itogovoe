package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

var (
	ErrIncomingData = errors.New("ошибка во входящей строке данных")
)

func parsePackage(data string) (int, time.Duration, error) {
	// TODO: реализовать функцию
	var buffer string = ""
	var durReturning time.Duration
	var errReturning error
	//Проверка, что введенные пользователем данные не пусты, и потенциально могут содержать 2 элемента
	dataParts, err := DataParts(data, 3, 2)
	if err != nil {
		errReturning = fmt.Errorf("Ошибка входящих данных '%v': %v\n", dataParts, err)
	}

	incomeDataSteps := []rune(dataParts[0])

	for i := 0; i < len(incomeDataSteps); i++ {
		if string(incomeDataSteps[i]) != " " {
			buffer += string(incomeDataSteps[i])
		} else {
			i++
		}
	}

	steps, err := strconv.Atoi(buffer)
	if err != nil {
		errReturning = fmt.Errorf("Ошибочный ввод количества шагов '%s': %v\n", buffer, err)
		return 0, 0, errReturning
	}

	incomeDataTime := []rune(dataParts[1])
	buffer = ""
	for i := 0; i < len(incomeDataTime); i++ {
		if string(incomeDataTime[i]) != " " {
			buffer += string(incomeDataTime[i])
		} else {
			i++
		}
	}

	duration, err := time.ParseDuration(buffer)
	if err != nil {
		errReturning = fmt.Errorf("Ошибка парсинга продолжительности '%s': %v\n", buffer, err)
		return 0, 0, errReturning
	}
	durReturning = duration

	return steps, durReturning, errReturning
}

func DataParts(data string, minLength, numberParts int) ([]string, error) {
	var errReturning error
	if len(data) < minLength {
		return []string{}, ErrIncomingData
	}

	dataParts := strings.Split(data, ",")

	if len(dataParts) != numberParts {
		return []string{}, ErrIncomingData
	}
	for _, i := range dataParts {
		if len(i) == 0 {
			return []string{}, ErrIncomingData
		}
	}

	return dataParts, errReturning
}

func DayActionInfo(data string, weight, height float64) string {
	// TODO: реализовать функцию
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if steps <= 0 {
		return ""
	}
	distance := (float64(steps) * stepLength) / mInKm
	calorics, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	strReturning := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2fкм.\nВы сожгли %.2fккал.", steps, distance, calorics)

	return strReturning
}
