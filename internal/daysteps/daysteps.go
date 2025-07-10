package daysteps

import (
	"errors"
	"fmt"
	"log"
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

// parsePackage разбирает строку данных на количество шагов и продолжительность активности.
func parsePackage(data string) (int, time.Duration, error) {
	splitData := strings.Split(data, ",")
	if len(splitData) != 2 {
		return 0, 0, errors.New("неправильное количество параметров")
	}

	steps, err := strconv.Atoi(splitData[0])
	if err != nil {
		return 0, 0, err
	}

	if steps <= 0 {
		return 0, 0, errors.New("неверное значение шагов")
	}

	duration, err := time.ParseDuration(splitData[1])
	if err != nil {
		return 0, 0, err
	}

	if duration <= 0 {
		return 0, 0, errors.New("неверная продолжительность - ноль")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	if !spentcalories.CheckWeight(weight) {
		fmt.Println("неверное значение веса", weight)
		return ""
	}
	if !spentcalories.CheckHeight(height) {
		fmt.Println("неверное значение роста", height)
		return ""
	}

	steps, duration, err := parsePackage(data)

	if err != nil {
		log.Println(err)
		return ""
	}

	distance := (float64(steps) * stepLength) / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	if err != nil {
		log.Println(err)
		return ""
	}

	dayActionInfo := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distance, calories,
	)

	return dayActionInfo
}
