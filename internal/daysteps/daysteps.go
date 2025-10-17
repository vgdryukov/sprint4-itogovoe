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

func parsePackage(data string) (int, time.Duration, error) {
	// TODO: реализовать функцию
	var durReturning time.Duration
	var errReturning error

	dataParts, err := spentcalories.DataParts(data, 3, 2)
	if err != nil {
		errReturning = fmt.Errorf("некорректная строка входящих данных '%v': %v", dataParts, err)
		return 0, 0, errReturning
	}

	steps, err := strconv.Atoi(dataParts[0])
	if err != nil {
		errReturning = fmt.Errorf("ошибка парсинга количества шагов '%s': %v", dataParts[0], err)
		return 0, 0, errReturning
	}
	if steps <= 0 {
		errReturning = errors.New("некорректное количество шагов")
		return 0, 0, errReturning
	}
	duration, err := time.ParseDuration(dataParts[1])
	if err != nil {
		errReturning = fmt.Errorf("ошибка парсинга продолжительности '%s': %v", dataParts[1], err)
		return 0, 0, errReturning
	}
	if duration <= 0 {
		errReturning = errors.New("некорректная продолжительность")
		return 0, 0, errReturning
	}
	durReturning = duration

	return steps, durReturning, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// TODO: реализовать функцию
	steps, duration, err := parsePackage(data)
	if err != nil {
		spentcalories.Logger(`dayAction.log`, "Восстановление в TrainingInfo() после паники в логгере: ", `dayAction `, "не удалось получить информацию о дневной активности: %v\n", err)
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
	strReturning := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calorics)

	return strReturning
}
