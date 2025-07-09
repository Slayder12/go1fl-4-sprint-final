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

// parseTraining разбирает строку с данными о тренировке
// Формат данных: "steps,activityType,duration" (например: "5000,Бег,1h30m")
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

	if duration <= 0 {
		return 0, "", 0, errors.New("неверная продолжительность - ноль")
	}

	return steps, splitData[1], duration, nil
}

// distance рассчитывает пройденную дистанцию в километрах
func distance(steps int, height float64) float64 {
	if !CheckHeight(height) {
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

// meanSpeed рассчитывает среднюю скорость в км/ч
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		fmt.Println("неверное значение продолжительности:", duration)
		return 0
	}

	distance := distance(steps, height)
	speed := distance / duration.Hours()

	return speed
}

// TrainingInfo формирует отчет о тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
	if !CheckWeight(weight) {
		return "", errors.New("неверное значение веса")
	}

	if !CheckHeight(height) {
		return "", errors.New("неверное значение роста")
	}

	steps, activityType, duration, err := parseTraining(data)

	if err != nil {
		log.Println(err)
		return "", err
	}

	distance := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	caloriesBurned, err := calculateCalories(activityType, steps, weight, height, duration)
	if err != nil {
		return "", fmt.Errorf("ошибка расчета калорий: %w", err)
	}

	return formatTrainingInfo(activityType, duration, distance, speed, caloriesBurned), nil
}

// RunningSpentCalories рассчитывает количество потраченных калорий при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("неверное значение шагов")
	}

	if !CheckWeight(weight) {
		return 0, errors.New("неверное значение веса")
	}

	if !CheckHeight(height) {
		return 0, errors.New("неверное значение роста")
	}

	if duration <= 0 {
		return 0, errors.New("неверное значение продолжительности:")
	}

	speed := meanSpeed(steps, height, duration)
	caloriesBurned := (weight * speed * duration.Minutes()) / minInH

	return caloriesBurned, nil
}

// WalkingSpentCalories рассчитывает количество потраченных калорий при ходьбе
// Использует RunningSpentCalories с понижающим коэффициентом
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	caloriesBurned, err := RunningSpentCalories(steps, weight, height, duration)

	if err != nil {
		return 0, err
	}

	walkingCalories := caloriesBurned * walkingCaloriesCoefficient

	return walkingCalories, nil
}

// CheckWeight проверяет корректность веса
func CheckWeight(weight float64) bool {
	if weight < minWeight {
		return false
	}
	if weight > maxWeight {
		return false
	}
	return true
}

// CheckHeight проверяет корректность роста
func CheckHeight(height float64) bool {
	if height < minHeight {
		return false
	}
	if height > maxHeight {
		return false
	}
	return true
}

// calculateCalories рассчитывает количество потраченных калорий в зависимости от типа активности.
func calculateCalories(activityType string, steps int, weight, height float64, duration time.Duration) (float64, error) {
	switch activityType {
	case "Бег":
		return RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		return WalkingSpentCalories(steps, weight, height, duration)
	default:
		return 0, errors.New("неизвестный тип тренировки")
	}
}

// formatTrainingInfo форматирует информацию о тренировке в читаемую строку.
func formatTrainingInfo(activityType string, duration time.Duration, distance, speed, calories float64) string {
	return fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		activityType, duration.Hours(), distance, speed, calories,
	)
}
