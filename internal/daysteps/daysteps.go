package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
	// количество минут в часе.
	minInH = 60
	// коэффициент для расчета калорий при ходьбе
	walkingCaloriesCoefficient = 0.5
	// коэффициент для расчета длины шага на основе роста
	stepLengthCoefficient = 0.45
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("Ожидаются два значения")
	}

	stepsStr := parts[0]
	durationStr := parts[1]

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("Не удалось преобразовать количество шагов в число: %w", err)
	}
	if steps <= 0 {
		return 0, 0, errors.New("Количество шагов должно быть больше 0")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("Не удалось парсить продолжительность прогулки: %w", err)
	}
	if duration <= 0 {
		return 0, 0, errors.New("Продолжительность прогулки должна быть положительным числом")
	}

	return steps, duration, nil
}

func distance(steps int, height float64) float64 {
	var calcStepLength float64

	if height > 0 {
		calcStepLength = height * stepLengthCoefficient
	} else {
		calcStepLength = stepLength
	}

	distanceMeters := float64(steps) * calcStepLength
	distanceKm := distanceMeters / mInKm

	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distanceKm := distance(steps, height)

	durationHours := duration.Hours()
	if durationHours == 0 {
		return 0
	}

	speed := distanceKm / durationHours

	return speed
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps < 0 {
		return 0, errors.New("Количество шагов должно быть положительным числом")
	}
	if weight <= 0 {
		return 0, errors.New("Вес должен быть положительным числом")
	}
	if height <= 0 {
		return 0, errors.New("Рост должен быть положительным числом")
	}
	if duration <= 0 {
		return 0, errors.New("Продолжительность должна быть положительной")
	}

	speed := meanSpeed(steps, height, duration)
	if speed == 0 {
		return 0, errors.New("не удалось рассчитать среднюю скорость (возможно, нулевая продолжительность или дистанция)")
	}

	durationMinutes := duration.Minutes()

	calories := (weight * speed * durationMinutes) / minInH

	calories *= walkingCaloriesCoefficient

	return calories, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	if steps <= 0 {
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	dist := distanceMeters / mInKm

	calories, errCalories := WalkingSpentCalories(steps, weight, height, duration)
	if errCalories != nil {
		log.Println(errCalories)
		return ""
	}

	result := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, dist, calories,
	)

	return result
}
