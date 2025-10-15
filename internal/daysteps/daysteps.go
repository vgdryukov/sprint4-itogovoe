package daysteps

import (
	"errors"
	"fmt"
	"strconv"
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
	//fmt.Println("---------func parsePackage (data): ", data)
	var buffer string = ""
	var durReturning time.Duration
	var errReturning error

	dataParts, err := spentcalories.DataParts(data, 3, 2)
	if err != nil {
		errReturning = fmt.Errorf("ошибка входящих данных '%v': %v", dataParts, err)
		return 0, 0, errReturning
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
		errReturning = fmt.Errorf("ошибочный ввод количества шагов '%s': %v", buffer, err)
		return 0, 0, errReturning
	}

	incomeDataDur := []rune(dataParts[1])
	buffer = ""
	for i := 0; i < len(incomeDataDur); i++ {
		if string(incomeDataDur[i]) != " " {
			buffer += string(incomeDataDur[i])
		} else {
			i++
		}
	}

	duration, err := time.ParseDuration(buffer)
	if err != nil {
		errReturning = fmt.Errorf("ошибка парсинга продолжительности '%s': %v", buffer, err)
		return 0, 0, errReturning
	}
	durReturning = duration

	return steps, durReturning, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// TODO: реализовать функцию
	//fmt.Println("------func DayActionInfo (data, weight, height): ", data, weight, height)
	steps, duration, err := parsePackage(data)
	if err != nil {
		return ""
	}
	if steps <= 0 {
		return ""
	}
	distance := (float64(steps) * stepLength) / mInKm
	calorics, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return ""
	}
	strReturning := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2fкм.\nВы сожгли %.2fккал.\n", steps, distance, calorics)

	return strReturning
}
