package testutil

var names = []string{
	"Liam", "Noah", "Oliver", "Elijah", "William", "James", "Benjamin", "Lucas", "Henry", "Alexander", "Mason", "Michael", "Ethan", "Daniel", "Logan", "Jack", "Caden", "Wyatt", "Grayson", "Julian", "Levi", "Isaiah", "Samuel", "Owen", "Roman", "Josiah", "Weston", "Cooper", "Finn", "Asher", "Caleb", "Miles", "Jace", "Theodore", "Leo", "Sebastian", "Jackson", "Landon", "Adam", "Xander", "Hudson", "Aiden", "Declan", "Hunter", "Luca", "Jaxon", "Colton", "Elliot", "Brooks"}

var nameCh chan string

func init() {
	nameCh = make(chan string, len(names))
	for _, name := range names {
		nameCh <- name
	}
}

func GetName() string {
	return <-nameCh
}

func PutName(name string) {
	nameCh <- name
}
