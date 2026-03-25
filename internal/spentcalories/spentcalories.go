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
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("Ожидаются три значения")
	}

	stepsStr := parts[0]
	activity := parts[1]
	durationStr := parts[2]

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("Не удалось преобразовать количество шагов в число: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("Количество шагов должно быть больше 0")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("Не удалось парсить продолжительность активности: %w", err)
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("Продолжительность должна быть положительной")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	var stepLength float64

	if height > 0 {
		stepLength = height * stepLengthCoefficient
	} else {
		stepLength = lenStep
	}

	distanceMeters := float64(steps) * stepLength
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

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height)

	speed := meanSpeed(steps, height, duration)

	var calories float64
	var errCalories error

	switch activity {
	case "Ходьба":
		calories, errCalories = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, errCalories = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if errCalories != nil {
		log.Println(errCalories)
		return "", errCalories
	}

	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity, duration.Hours(), dist, speed, calories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
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
		return 0, errors.New("Не удалось рассчитать среднюю скорость")
	}

	durationMinutes := duration.Minutes()

	calories := (weight * speed * durationMinutes) / minInH

	return calories, nil
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
