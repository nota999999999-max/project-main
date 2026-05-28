package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)
var droppers = []string{
	"Wirecard",
	"Binance",
}
var senders = []string{
	"John Tomson",
	"Van Li",
}
var resivers = []string{
	"John Tomson",
	"Van Li",
}
type CardIdentifyModel struct {
	PAN string `json:"pan"`
}
type ReceivedData struct {
	PAN string `json:"pan"`
}
type Transaction struct {
	ID       int     `json:"id"`           //serial
	Amount   float64 `json:"money_amount"` //number
	Currency string  `json:"currency"`     //text (TJS, USD)
	PAN      string  `json:"pan"`
	Receiver string  `json:"receiver_name"` //text
	Sender   string  `json:"sender_name"`   //text (Name)
	Provider string  `json:"host"`          //text (Binance)
	Status string `json:"status"` //text (Success/Decline)
}
type Sender struct {
	SenderName string `json:"sender_name"`
	Provider   string `json:"sender_host"`
}
type Resiver struct {
	ResiverName string `json:"resiver_name"`
	Provider    string `json:"resiver_host"`
	MobileN     string `json:"mobile_number"`
	Card        string `json:"card_pan"`
}

func main() {
	r := mux.NewRouter()
	
	r.HandleFunc("/card", CardIdentify).Methods("POST")
	r.HandleFunc("/format", CardFormat).Methods("POST")
	r.HandleFunc("/record", TransactionsRecord).Methods("POST")

	fmt.Println("8080")
	http.ListenAndServe(":8080", r)
}

func CardIdentify(w http.ResponseWriter, r *http.Request) {
	var userCard CardIdentifyModel
	json.NewDecoder(r.Body).Decode(&userCard)
	
	if len(userCard.PAN) != 16 {
		w.Write([]byte("Invalid PAN"))
		return
	}
	fmt.Println(userCard.PAN)
	w.Write([]byte("Fisical Card"))
}

func CardFormat(w http.ResponseWriter, r *http.Request) {
	var data ReceivedData
	json.NewDecoder(r.Body).Decode(&data)

	if data.PAN[:4] == "5440" || data.PAN[:4] == "9860" || data.PAN[:4] == "8600" {

		fmt.Println(data.PAN)
		w.Write([]byte("Its KortiMilli"))
	}
	if data.PAN[:1] == "4" {

		fmt.Println(data.PAN)
		w.Write([]byte("Its Visa"))
	}
	if data.PAN[:2] == "62" || data.PAN[:2] == "81" {

		fmt.Println(data.PAN)
		w.Write([]byte("Its UnionPay"))
	}
}

func TransactionsRecord(w http.ResponseWriter, r *http.Request) {
	var transaction Transaction
	json.NewDecoder(r.Body).Decode(&transaction)

	//тут проверка
	if CheckDroppers(transaction) == true {
		//тут процесс токенизации - подмена пана карты на токен
		rand.Seed(time.Now().UnixNano())

		var token string

		for i := 0; i < 16; i++ {
			token += fmt.Sprintf("%d", rand.Intn(10))
		}
		transaction.PAN = token

		fmt.Println("Обнаружен дроперский платёж!")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(transaction)
	}
	//тут отказ
	if CheckDroppers(transaction) == true {

		transaction.Status = "DECLINE"

		fmt.Println("Откланён по причине антифрода")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(transaction)
		return

	} else {
		PaymentProcess(transaction)

			transaction.Status = "SUCCESS"

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(transaction)
	}

}

func CheckDroppers(trn Transaction) bool {

	for _, dropper := range droppers {

		if trn.Provider == dropper {
			return true
		}
	}
	for _, sender := range senders {

		if trn.Sender == sender {
			return true
		}
	}
	for _, resiver := range resivers {

		if trn.Receiver == resiver {
			return true
		}
	}
	return false
}
func PaymentProcess(t Transaction) bool {

	if t.Amount <= 0 {
		return false
	}
	if len(t.PAN) != 16 {
		return false
	}

	fmt.Println("Перевод выполнен:", t.Amount)
	return true
}


