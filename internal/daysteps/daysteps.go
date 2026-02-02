package daysteps

import (
	"fmt"
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

func parsePackage(data string) (int, time.Duration, error) {
	strSplit := strings.Split(data, ",")
	if len(strSplit) != 2 {
		return 0, 0, fmt.Errorf("Длина слайса: %d", len(strSplit))
	}
	steps, err := strconv.Atoi(strSplit[0])
	if err != nil {
		return 0, 0, err
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов: %d", steps)
	}
	time, err := time.ParseDuration(strSplit[1])
	if err != nil {
		return 0, 0, err
	}
	return steps, time, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, time, err := parsePackage(data)
	if err != nil {
		fmt.Println("ошибка парсинга шагов и времени прогулки:", err)
		return ""
	}
	if steps <= 0 {
		return ""
	}
	var distance float64 = float64(steps) * stepLength
	distanceInKm := distance / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, time)
	if err != nil {
		fmt.Println("ошибка вычисления калорий", err)
		return ""
	}
	information := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceInKm, calories)
	return information
}
