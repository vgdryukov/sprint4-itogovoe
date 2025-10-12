package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/daysteps"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

var (
	ErrIncomingData = errors.New("ошибка во входящих данных")
	ErrCalcCalorics = errors.New("ошибка расчета каллорий")
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// TODO: реализовать функцию
	var steps int
	var errReturning error
	var buffer string = ""
	var activity string = ""
	var durReturning time.Duration = 0

	dataParts, err := daysteps.DataParts(data, 10, 3)
	if err != nil {
		errReturning = fmt.Errorf("Ошибка ввода данных '%v': %v\n", dataParts, err)
	}

	for i := 0; i < 3; i++ {
		incomeData := []rune(dataParts[i])

		for k := 0; k < len(incomeData); k++ {
			if string(incomeData[k]) != " " {
				buffer += string(incomeData[k])
			} else {
				k++
			}
		}
		switch i {
		case 0:
			num, err := strconv.Atoi(buffer)
			if err != nil {
				errReturning = fmt.Errorf("Ошибочный ввод количества шагов '%s': %v\n", buffer, err)
				break
			}
			steps = num
		case 1:
			activity = buffer
		case 2:
			duration, err := time.ParseDuration(buffer)
			if err != nil {
				errReturning = fmt.Errorf("Ошибка парсинга продолжительности '%s': %v\n", buffer, err)
				break
			}
			durReturning = duration
		}
	}
	return steps, activity, durReturning, errReturning
}

func distance(steps int, height float64) float64 {
	// TODO: реализовать функцию
	distance := float64(steps) * height * stepLengthCoefficient / mInKm
	return distance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// TODO: реализовать функцию
	if duration <= 0 {
		return 0
	}
	speed := distance(steps, height) / duration.Hours()
	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// TODO: реализовать функцию
	var errReturning error
	var calorics float64

	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		errReturning = fmt.Errorf("Не получилось получить информацию о тренировке: %v", err)
		//		log.Printf("не получилось получить информацию о тренировке: %v", err)
	}

	switch activity {
	case "Ходьба":
		calorics, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", ErrCalcCalorics
		}
	case "Бег":
		calorics, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", ErrCalcCalorics
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	distance := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	strReturning := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2fч.\nДистанция: %.2fкм.\nСкорость: %.2fкм/ч\nСожгли калорий: %.2f", activity, duration, distance, speed, calorics)
	return strReturning, errReturning
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	if steps < 1 || weight < 10 || height < 1 || duration.Minutes() < 1 {
		return 0, ErrIncomingData
	}
	calorics := (weight * meanSpeed(steps, height, duration) * duration.Minutes()) / minInH

	return calorics, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	if steps < 1 || weight < 10 || height < 1 || duration.Minutes() < 1 {
		return 0, ErrIncomingData
	}
	calorics := walkingCaloriesCoefficient * (weight * meanSpeed(steps, height, duration) * duration.Minutes()) / minInH

	return calorics, nil
}
