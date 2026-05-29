package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

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
	Status   string  `json:"status"`        //text (Success/Decline)
	Message  string  `json:"message"`       //text (...)
	Token    string  `json:"token"`
}
type TransactionResponse struct {
	Token   string `json:"token"`
	Status  string `json:"status"`
	Message string `json:"message"`
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

	fmt.Println("List 8080")
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

	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Transaction:", transaction)

	//тут оригинальный PAN
	originalPAN := transaction.PAN

	//тут проверка антифрода
	if CheckDroppers(transaction) {

		//тут сохранён оригинальный PAN
		originalPAN := transaction.PAN

		//это вывод в терминал
		fmt.Println("Original PAN:", originalPAN)

		//тут токенизация
		var token string

		for i := 0; i < 16; i++ {
			token += strconv.Itoa(rand.Intn(10))
		}

		//тут вывод токена и уведомление в терминал
		fmt.Println("Generated token:", token)

		fmt.Println("Обнаружен дроперский платёж!")

		//тут ответ отправителю-sender
		response := TransactionResponse{
			Token:   MaskToken(token),
			Status:  "DECLINE",
			Message: "Неправильный номер карты",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		return
	}

	//платёж
	if PaymentProcess(transaction) {

		transaction.Status = "SUCCESS"
		transaction.Message = "Платёж выполнен"

		fmt.Println("Платёж выполнен")

	} else {

		transaction.Status = "DECLINE"
		transaction.Message = "Ошибка проведения платежа"

		fmt.Println("Ошибка проведения платежа")
	}

	//тут маскировка пана
	transaction.PAN = MaskPAN(originalPAN)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)

}

func CheckDroppers(trn Transaction) bool {

	fmt.Println("Transaction:", trn)

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

	fmt.Println("Success:", t.Amount)
	return true
}
func MaskPAN(pan string) string {
	if len(pan) != 16 {
		return pan
	}
	return pan[:4] + "********" + pan[12:]
}
func MaskToken(token string) string {

	if len(token) != 16 {
		return token
	}

	return token[:4] + "********" + token[12:]
}
