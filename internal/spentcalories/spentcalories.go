package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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
	ErrIncomingData      = errors.New("ошибка во входящей строке данных")
	ErrIncomingDataParts = errors.New("ошибка в количестве частей входящей строки")
	ErrIncomingDataSteps = errors.New("ошибка в количестве шагов")
	ErrIncomingDataDur   = errors.New("ошибка в продолжительности")
	ErrCalcCalorics      = errors.New("ошибка расчета каллорий")
)

func DataParts(data string, minLength, numberParts int) ([]string, error) {
	//fmt.Println("------------func DataParts (data, minLength, numberParts): ", data, minLength, numberParts)
	if len(data) < minLength {
		return []string{}, ErrIncomingData
	}

	dataParts := strings.Split(data, ",")
	//fmt.Println("------------func DataParts slise dataParts, len(dataParts): ", dataParts, len(dataParts))
	if len(dataParts) != numberParts {
		//fmt.Println("---------------func DataParts RETURN-0: ", []string{}, ErrIncomingDataParts)
		return []string{}, ErrIncomingDataParts
	}
	if len(dataParts[0]) == 0 {
		//fmt.Println("---------------func DataParts RETURN-1: ", []string{}, ErrIncomingDataParts)
		return []string{}, ErrIncomingDataSteps
	}
	if len(dataParts[1]) == 0 {
		//fmt.Println("---------------func DataParts RETURN-2: ", []string{}, ErrIncomingDataParts)
		return []string{}, ErrIncomingDataDur
	}
	//fmt.Println("---------------func DataParts RETURN-NORM: ", dataParts, nil)
	return dataParts, nil
}

func parseTraining(data string) (int, string, time.Duration, error) {
	// TODO: реализовать функцию
	//fmt.Println("---------func parseTraining(data): ", data)
	var steps int = 0
	var errReturning error = nil
	var activity string = ""
	var durReturning time.Duration = 0

	dataParts, err := DataParts(data, 10, 3)
	//fmt.Println("---------func parseTraining slice dataParts: ", dataParts)
	if err != nil {
		errReturning = fmt.Errorf("ошибка входящих данных '%v': %v", dataParts, err)
		return steps, activity, durReturning, errReturning
	}

	for i := 0; i < 3; i++ {
		incomeData := []rune(dataParts[i])
		var buffer string = ""
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
				errReturning = fmt.Errorf("ошибочный ввод количества шагов '%s': %v", buffer, err)
				return steps, activity, durReturning, errReturning
			}
			steps = num
		case 1:
			activity = buffer
		case 2:
			//fmt.Println("------------func parseTraining duration buffer: ", buffer)
			duration, err := time.ParseDuration(buffer)
			//fmt.Println("---------------func parseTraining duration: ", duration)
			if err != nil {
				errReturning = fmt.Errorf("ошибка парсинга продолжительности '%s': %v", buffer, err)
				return 0, "", 0, errReturning
			}
			durReturning = duration
		}
	}
	return steps, activity, durReturning, nil
}

func distance(steps int, height float64) float64 {
	// TODO: реализовать функцию
	distance := float64(steps) * height * stepLengthCoefficient / mInKm
	return distance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// TODO: реализовать функцию
	//fmt.Println("------------func meanSpeed (steps, height, duration):", steps, height, duration)
	if duration <= 0 {
		return 0
	}
	speed := distance(steps, height) / duration.Hours()
	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// TODO: реализовать функцию
	//fmt.Println("------func TrainingInfo (data, weight, height): ", data, weight, height)
	var errReturning error
	var calorics float64

	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		errReturning = fmt.Errorf("не удалось получить информацию о тренировке: %v", err)

		flog, err := os.OpenFile(`server.log`, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Восстановление в TrainingInfo() после паники в логгере: ", r)
				}
			}()
		}
		defer flog.Close()
		mylog := log.New(flog, `serv `, log.LstdFlags|log.Lshortfile)
		mylog.Printf("не удалось получить информацию о тренировке: %v\n", err)

		return "", errReturning
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
	durationInHours := duration.Hours()
	distance := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	strReturning := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2fч.\nДистанция: %.2fкм.\nСкорость: %.2fкм/ч\nСожгли калорий: %.2f\n", activity, durationInHours, distance, speed, calorics)
	return strReturning, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	//fmt.Println("------------func RunningSpentCalories (steps, weight, height, duration): ", steps, weight, height, duration)
	if steps < 1 || weight < 10 || height < 1 || duration.Minutes() < 1 {
		return 0, ErrIncomingData
	}
	calorics := (weight * meanSpeed(steps, height, duration) * duration.Minutes()) / minInH

	return calorics, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	//fmt.Println("------------func WalkingSpentCalories (steps, weight, height, duration): ", steps, weight, height, duration)
	if steps < 1 || weight < 10 || height < 1 || duration.Minutes() < 1 {
		return 0, ErrIncomingData
	}
	calorics := walkingCaloriesCoefficient * (weight * meanSpeed(steps, height, duration) * duration.Minutes()) / minInH

	return calorics, nil
}
