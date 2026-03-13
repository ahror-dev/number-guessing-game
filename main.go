package main 

import (
	"fmt"
	"bufio"
	"io"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"

)


// qiyinchilik darajasi
type DifficultyLevel int

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
	score, exists := hs.scores[level]
	return score, exists	
}


// baland natijalarni print qilish
func (hs *HighScores) Print(w io.Writer) {
	fmt.Println("Eng yaxshi natijalar: ")
	for _, level := range AllLevels(){
		if score, exists := hs.scores[level]; exists {
			fmt.Fprintf(w, " \n %-8s : %d urinish\n", level, score)

		}else {
			fmt.Fprintf(w, "\n%s Oyin hali oynalmagan!", level)
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
	return map[HintLevel]string{
		HintFire: "🔥 Juda Yaqin! Bir qadam qoldi.",
		HintHot: "♨️ Yaqinlashyapsan.",
		HintWarm: "😐 To'g'ri yo'ldasan",
		HintIce: "❄️ Ancha uzoqda",
		
	}[hintLevel(secret, guess)]
}




// O`quvchi qism
type Reader struct {
	r *bufio.Reader
}


func NewReader(input io.Reader) *Reader {
	return &Reader {r: bufio.NewReader(input)}
}


func (rd *Reader) ReaderLine() (string, error) {
	line, err := rd.r.ReadString('\n')
	return strings.TrimSpace(line), err 
}



func (rd *Reader) ReadInt(min, max int) (int, error) {
	line, err := rd.ReaderLine()
	if err != nil {
		return 0, err
	}

	n, err := strconv.Atoi(line)
	if err != nil {
		return 0, fmt.Errorf("invalid: %q (expected %d-%d)", line, min, max)
	}
	if n < min || n > max {
		return 0, fmt.Errorf("number out of range: %d (expected %d-%d)", n, min, max)
	}
	return n, nil
}


func (rd *Reader) ReadConfirm() bool {
	line, _ := rd.ReaderLine()
	line = strings.ToLower(line)
	return line == "ha" || line == "h" || line == "yes" || line == "y"
}




// ROUND

// RoundResult Structi
type RoundResult struct {
	Won bool
	Secret int
	Attempts int
	Elapsed time.Duration
	Level DifficultyLevel

}


// Round Structi
type Round struct {
	secret int
	level DifficultyLevel
	reader *Reader
	writer io.Writer

}




func NewRound(secret int, level DifficultyLevel, r *Reader, w io.Writer) *Round {
	return &Round{secret: secret, level: level, reader: r, writer: w}
}





func (r *Round) Play() RoundResult {
	fmt.Fprintf(r.writer, "Oyin boshlandi! (1-100 oraligidagi sonni toping)\n\n")

	start := time.Now()
	maxChances := r.level.Chances()

	for attempts := 1; attempts <= maxChances; attempts++ {
		remaining := maxChances - attempts + 1
		fmt.Fprintf(r.writer, "%d/%d urinishlaringiz qoldi!\n", remaining, maxChances)


		guess, err := r.reader.ReadInt(1, 100)

		if err != nil {
			fmt.Fprintln(r.writer, "Iltimos 1-100 oraligidagi butun sonlarni yozing!")
			attempts--
			continue
		}

		if guess == r.secret {
			elapsed := time.Since(start).Round(time.Second)
			fmt.Fprintf(r.writer, "\n🎉 To'g'ri! Son %d edi\n", r.secret)
			fmt.Fprintf(r.writer, "%d urinishda topdingiz(%v)\n", attempts, elapsed)
			return RoundResult{Won: true,Secret: r.secret,	Attempts: attempts,Elapsed: elapsed,Level: r.level}
		}

		if guess < r.secret {
			fmt.Fprintf(r.writer, "%d - bu son maxfiy sondan kichik!\n", guess)
		}else {
			fmt.Fprintf(r.writer, "%d - bu son maxfiy sondan katta\n", guess)
		}

		fmt.Fprintln(r.writer, Hint(r.secret, guess, remaining-1))
		fmt.Fprintln(r.writer)

	}

	elapsed := time.Since(start).Round(time.Second)
	fmt.Fprintf(r.writer, "\n💔 Urinishlaringiz tugadi, Togri son - %d. (%v)\n", r.secret, elapsed)
	return RoundResult{Won: false, Secret: r.secret, Attempts: maxChances, Elapsed: elapsed, Level: r.level}

}


// O`YIN
type Game struct {
	reader *Reader
	writer io.Writer
	scores *HighScores
}



func NewGame(input io.Reader, output io.Writer) *Game {
	return &Game {
		reader: NewReader(input),
		writer: output,
		scores: NewHighScores(),


	}
}

func (g *Game) Run() {
	g.printWelcome()


	for {
		level, err := g.selectDifficulty()
		if err != nil {
			continue
		}
	

		secret := rand.IntN(100) + 1
		result := NewRound(secret, level, 	g.reader, g.writer).Play() 

		if result.Won {
			g.scores.Update(result.Level, result.Attempts)
		}
		g.scores.Print(g.writer)

		fmt.Fprint(g.writer, "\nYana oynamoqchimisiz? (ha/yoq): ")

		if !g.reader.ReadConfirm() {
			break
		}

		fmt.Fprintln(g.writer, "\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	}

	fmt.Fprintln(g.writer, "Korishguncha ")
}

func (g *Game) selectDifficulty() (DifficultyLevel, error) {
	fmt.Fprintln(g.writer, "Qiyinchilik Darajasini Tanlang(sonlarda 1/2/3)")
	for _, level := 	range AllLevels() {
		fmt.Fprintf(g.writer, "  %d. %-8s (%d urinish)\n", int(level)+1, level, level.Chances())
	}

	n, err := g.reader.ReadInt(1,3)
	if err != nil {
		fmt.Fprintln(g.writer, "⚠️  1, 2 yoki 3 ni kiriting.")
		return 0, err
	}

	level := DifficultyLevel(n-1)
	fmt.Fprintf(g.writer, "%s tanlandi, %d urinish\n", level, level.Chances())
	return level, nil

}



func (g *Game) printWelcome() {
	fmt.Fprintln(g.writer, "╔══════════════════════════════════════╗")
	fmt.Fprintln(g.writer, "║      NUMBER GUESSING GAME             ║")
	fmt.Fprintln(g.writer, "╚══════════════════════════════════════╝")
	fmt.Fprintln(g.writer, "\n📋 QOIDALAR:")
	fmt.Fprintln(g.writer, "   1. Kompyuter 1–100 orasida son tanlaydi")
	fmt.Fprintln(g.writer, "   2. Qiyinchilik darajasini tanlaysiz")
	fmt.Fprintln(g.writer, "   3. O`yindan zavqlaning!")
	fmt.Fprintln(g.writer, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(g.writer, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

func main() {
	NewGame(os.Stdin, os.Stdout).Run()
}