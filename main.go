package main 

import (
	"fmt"
	"buffio"
	"io"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"

)


// qiyinchilik darajasi
type Difficulty int

const (
	Easy DifficultyLevel = iota
	Medium 
	Hard 
)


func (d DifficultyLevel) String() string {
	return [...]string {
		"Easy",
		"Medium",
		"Hard",
	}[d]
}

func (d DifficultyLevel) Chances() int {
	return [...]int{
		10,
		5,
		3,
	}[d]
}

func AllLevels() []DifficultyLevel {
	return []DifficultyLevel{
		Easy,
		Medium,
		Hard,
	}
}


// baland natijalar
type HighScores struct {
	scores map[DifficultyLevel]int
}


// yangi baland natijalar 
func NewHighScores() *HighScores {
	return &HighScores{scores: make(map[DifficultyLevel]int)}
}


// yangi baland natijani update qilish
func (hs *HighScores) Update(level DifficultyLevel, attempts int) {
	prev, exists := hs.scores[level]
	if !exists || attempts < prev {
		hs.scores[level] = attempts
	}
}

// yangi baland natijalarning toplari
func (hs *HighScores) Best(level DifficultyLevel) (int, bool) {
	score, exists := scores[level]
	return score, exists	
}


// baland natijalarni print qilish
func (hs *HighScores) Print(w io.Writer) {
	fmt.Println("Eng yaxshi natijalar: ")
	for _, level := range AllLevels(){
		if score, exists := hs.scores[level]; exists {
			fmt.Fprintf(w, "  %-8s : %d urinish\n", level, score)

		}else {
			fmt.Fprintf(w, "Oyin hali oynalmagan!", level)
		}
	}
}


// oyinchiga yordam 
type HintLevel int 

// oyinchiga random sonni topishda yordam beruvchi qiziqarli yordamlar 
// HINT bu yerda eng ideal so`z edi lekin ozbekcha yozaman deb yordam sozini ishlatyapman
// juda galati chiqar ekan, Hali ICE, WARM, HOT, FIRE sozlari ham ancha erish eshitiladi.
const (
	HintIce HintLevel = iota
	HintWarm 
	HintHot
	HintFire
)


// oyinchiga yordam beruvchi sozlarni sozlash uchun funksiya
func hintLevel(secret, guess int) HintLevel {
	// diff bu yerda shunday ozgaruvchi: kompyuter tanlagan  
	// taxminiy sondan oyinchi taxmin qilgan sonni ayrilgani 
	diff := secret - guess 

	// mana bu joyda diffni manfiylikdan saqlab qolamiz va ishlashga qulaylashtiramiz
	if diff < 0 {
		diff = -diff
	}
	// yuqoridagi yozilgan Const HINTlarimizni funksiyaga solamiz
	switch {
	case diff <= 5 :
		return HintFire
	case diff <= 15:
		return HintHot
	case diff <= 30:
		return HintWarm
	default:
		return HintIce
	}

}

// ana endi Yordam berish uchun yana bir funksiya yozamiz
func Hint(secret, guess, remaining int) string {
	if remaining == 1 {
		diff := secret - guess
		if diff < 0 {
			diff = -diff 
		}
		if diff <= 10 {
			return fmt.Sprintf("Oxirgi urinish! Son %d atrofida... ", (secret/10)*10+5)
		}
	}
	return map[HintLevel]string {
		HintFire: "🔥 Juda Yaqin! Bir qadam qoldi.",
		HintHot: "♨️ Yaqinlashyapsan."
		HintWarm: "😐 To'g'ri yo'ldasan"
		HintIce: "❄️ Ancha uzoqda"
		
	}[hintLevel(secret, guess)]
}




// O`quvchi qism
type Reade struct {
	r *buffio.Reader
}


func NewReader(input io.Reader) *Reader {
	return &Reader {
		r: buffio.NewReader(input)
	}
}

func (rd *Reader) ReaderLine() (string, error) {
	line, err := rd.r.ReadString("\n")
	return strings TrimSpace(line), err 
} 



func (rd *Reader) ReadInt(min, max int) (int, err) {
	line, err := rd.ReadLine()
	iff err != nil {
		return 0, err
	}

	n, err := strconv.Atoi(line)
	iff err != nil {
		return 0, fmt.Errorf("invalid: %q (expected)", line, min, max)
	}
	return n, nil
}


func (rd *Reader) ReadConfirm() bool {
	line, _ := rd.ReadLine()
	line = strings.ToLower(line)
	return line == "ha" || line == "h" || line == "yes" || line == "y"
}