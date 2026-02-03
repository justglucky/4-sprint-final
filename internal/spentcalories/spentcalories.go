package spentcalories

import (
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
	strSplit := strings.Split(data, ",")
	if len(strSplit) != 3 {
		return 0, "", 0, fmt.Errorf("Длина слайса: %d", len(strSplit))
	}
	steps, err := strconv.Atoi(strSplit[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка в преобразовании шагов")
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов: %d", steps)
	}
	activity := strSplit[1]
	time, err := time.ParseDuration(strSplit[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка в преобразовании времени")
	}
	if time <= 0 {
		return 0, "", 0, fmt.Errorf("количество пройденного времени: %v", time)
	}
	return steps, activity, time, nil
}

func distance(steps int, height float64) float64 {
	strideLength := height * stepLengthCoefficient
	distanceInM := float64(steps) * strideLength
	distanceInKm := distanceInM / mInKm
	return distanceInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distanceInKm := distance(steps, height)
	hours := duration.Hours()
	averageMoveSpeed := distanceInKm / hours
	return averageMoveSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, time, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	hours := time.Hours()
	switch activity {
	case "Ходьба":
		distanceInKm := distance(steps, height)
		averageMoveSpeed := meanSpeed(steps, height, time)
		calories, err := WalkingSpentCalories(steps, weight, height, time)
		if err != nil {
			return "", err
		}
		var information string = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, hours, distanceInKm, averageMoveSpeed, calories)
		return information, nil
	case "Бег":
		distanceInKm := distance(steps, height)
		averageMoveSpeed := meanSpeed(steps, height, time)
		calories, err := RunningSpentCalories(steps, weight, height, time)
		if err != nil {
			return "", err
		}
		information := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n", activity, hours, distanceInKm, averageMoveSpeed, calories)
		return information, nil
	}
	return "", fmt.Errorf("неизвестный тип тренировки")
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || duration <= 0 || weight <= 0 || height <= 0 {
		return 0, fmt.Errorf("количество шагов: %d. пройденное время: %v. вес: %f. рост: %f", steps, duration, weight, height)
	}
	averageMoveSpeed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()
	calories := (weight * averageMoveSpeed * minutes) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || duration <= 0 || weight <= 0 || height <= 0 {
		return 0, fmt.Errorf("количество шагов: %d. пройденное время: %v. вес: %f. рост: %f", steps, duration, weight, height)
	}
	averageMoveSpeed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()
	calories := (weight * averageMoveSpeed * minutes) / minInH
	walkingCalories := calories * walkingCaloriesCoefficient
	return walkingCalories, nil
}
