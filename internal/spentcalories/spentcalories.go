package spentcalories

import (
	"errors"
	"fmt"
	"log"
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

	minWeight = 2.0   // Минимальный допустимый вес (кг)
	maxWeight = 635.0 // Максимальный допустимый вес (кг)
	minHeight = 0.50  // Минимальный допустимый рост (м)
	maxHeight = 2.75  // Максимальный допустимый рост (м)
)

func parseTraining(data string) (int, string, time.Duration, error) {
	splitData := strings.Split(data, ",")
	if len(splitData) != 3 {
		return 0, "", 0, errors.New("неправильное количество параметров")
	}

	steps, err := strconv.Atoi(splitData[0])
	if err != nil {
		return 0, "", 0, err
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("неверное значение шагов")
	}

	duration, err := time.ParseDuration(splitData[2])
	if err != nil {
		return 0, "", 0, err
	}

	if duration == 0 {
		fmt.Println("неверная продолжительность - ноль")
		return 0, "", 0, err
	}

	return steps, splitData[1], duration, nil
}

func distance(steps int, height float64) float64 {

	if !ValidateHeight(height) {
		fmt.Println("неверное значение роста:", height)
		return 0
	}

	if steps <= 0 {
		fmt.Println("неверное значение шагов:", steps)
		return 0
	}

	stepLength := height * stepLengthCoefficient
	distance := (float64(steps) * stepLength) / mInKm
	return distance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0 {
		fmt.Println("неверное значение продолжительности:", duration)
		return 0
	}

	distance := distance(steps, height)
	speed := distance / duration.Hours()

	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {

	if !ValidateWeight(weight) {
		return "", errors.New("неверное значение веса")
	}

	if !ValidateHeight(height) {
		return "", errors.New("неверное значение роста")
	}

	steps, typeActivity, duration, err := parseTraining(data)

	if err != nil {
		log.Println(err)
		return "", err
	}

	distance := distance(steps, height)
	meanSpeed := meanSpeed(steps, height, duration)

	var caloriesBurned float64

	switch typeActivity {
	case "Бег":
		caloriesBurned, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "Ходьба":
		caloriesBurned, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	trainingInfo := fmt.Sprintf(
		"Тип тренировки: %s\n"+
		"Длительность: %.2f ч.\n"+
		"Дистанция: %.2f км.\n"+
		"Скорость: %.2f км/ч\n"+
		"Сожгли калорий: %.2f",
		typeActivity, duration.Hours(), distance, meanSpeed, caloriesBurned,
	)

	return trainingInfo, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 {
		return 0, errors.New("неверное значение шагов")
	}

	if !ValidateWeight(weight) {
		return 0, errors.New("неверное значение веса")
	}

	if !ValidateHeight(height) {
		return 0, errors.New("неверное значение роста")
	}

	if duration <= 0 {
		return 0, errors.New("неверное значение продолжительности:")
	}

	speed := meanSpeed(steps, height, duration)
	caloriesBurned := (weight * speed * duration.Minutes()) / minInH

	return caloriesBurned, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	caloriesBurned, err := RunningSpentCalories(steps, weight, height, duration)

	if err != nil {
		return 0, err
	}

	walkingCalories := caloriesBurned * walkingCaloriesCoefficient

	return walkingCalories, nil
}

// Функция проверки веса (в килограммах)
func ValidateWeight(weight float64) bool {

	if weight < minWeight {
		return false
	}
	if weight > maxWeight {
		return false
	}
	return true
}

// Функция проверки роста (в метрах)
func ValidateHeight(height float64) bool {

	if height < minHeight {
		return false
	}
	if height > maxHeight {
		return false
	}

	return true
}
