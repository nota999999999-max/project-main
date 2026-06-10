package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/fogleman/gg"
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
var clients = []Client{
	{Name: "Muhammadjon S", MobileN: "+992000111111", Card: "5440333322221111"},
	{Name: "Parviz H", MobileN: "+992000222222", Card: "4000444433332222"},
	{Name: "Zarina A", MobileN: "+992000333333", Card: "4111555544443333"},
	{Name: "Bezhan Sh", MobileN: "+992000444444", Card: "4276666655554444"},
	{Name: "Guldofarin Kh", MobileN: "+992000555555", Card: "6211777766665555"},
}
var transactions = []Transaction{

	{ID: 1, Amount: 150.75, Currency: "TJS", PAN: "5440333322221111",
		Sender: "Nozim G", Provider: "Alif", Status: "SUCCESS",
		IP: "192.168.1.10", Country: "Tajikistan"},
	{ID: 2, Amount: 500.00, Currency: "USD", PAN: "4000444433332222",
		Sender: "John Tomson", Provider: "Wirecard", Status: "DECLINE",
		IP: "185.44.77.21", Country: "Germany"},
	{ID: 3, Amount: 320.00, Currency: "TJS", PAN: "4111555544443333",
		Sender: "Alijon R.", Provider: "Wallet", Status: "SUCCESS",
		IP: "46.101.12.55", Country: "UAE"},
	{ID: 4, Amount: 89.99, Currency: "TJS", PAN: "4276666655554444",
		Sender: "Van Li", Provider: "Binance", Status: "DECLINE",
		IP: "8.8.8.8", Country: "USA"},
	{ID: 5, Amount: 75.99, Currency: "EUR", PAN: "6211777766665555",
		Sender: "N/A", Provider: "PayPal", Status: "SUCCESS",
		IP: "185.4.4.4", Country: "China"},
}
var geoPoints = []GeoPoint{
	{
		Country: "Tajikistan",
		IP:      "185.1.1.1",
	},
	{
		Country: "Germany",
		IP:      "185.2.2.2",
	},
	{
		Country: "USA",
		IP:      "185.3.3.3",
	},
	{
		Country: "China",
		IP:      "185.4.4.4",
	},
}
var lastAlert string

type CardIdentifyModel struct {
	PAN string `json:"pan"`
}
type ReceivedData struct {
	PAN string `json:"pan"`
}
type Transaction struct {
	ID            int       `json:"id"`           //serial
	Amount        float64   `json:"money_amount"` //number
	Currency      string    `json:"currency"`     //text (TJS, USD)
	PAN           string    `json:"pan"`
	Receiver      string    `json:"receiver_name"` //text
	Sender        string    `json:"sender_name"`   //text (Name)
	Provider      string    `json:"host"`          //text (Binance)
	Status        string    `json:"status"`        //text (Success/Decline)
	Message       string    `json:"message"`       //text (...)
	Token         string    `json:"token"`
	IP            string    `json:"ip"` //ip country
	Country       string    `json:"country"`
	Alert         string    `json:"alert"`
	Time          time.Time `json:"time"`
	TimeFormatted string    `json:"time_formatted"`
}
type GeoPoint struct {
	Country string `json:"country"`
	IP      string `json:"ip"`
	Status  string `json:"status"`
}
type GeoLocation struct {
	Country string `json:"country"`
}

type TransactionResponse struct {
	Status  string `json:"status"`
	Token   string `json:"token"`
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
type Client struct {
	Name    string `json:"user_name"`
	MobileN string `json:"mobile_number"`
	Card    string `json:"card"`
}
type CardInfo struct {
	CardNumber string `json:"card_number"`
	Format     string `json:"format"`
	ClientName string `json:"client_name"`
	Phone      string `json:"mobile"`
}
type CardPage struct {
	CardInfo CardInfo
	History  []Transaction
}
type Dashboard struct {
	FraudScore   int
	Velocity     string
	Transactions int
	TxList       []Transaction
	GeoPoints    []GeoPoint

	Alert string
	PAN   string
	Token string
}

func main() {

	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/card", CardIdentify).Methods("POST")
	r.HandleFunc("/format", CardFormat).Methods("POST")
	r.HandleFunc("/record", TransactionsRecord).Methods("POST")
	r.HandleFunc("/data", CardData).Methods("GET")

	// Dashboard
	r.HandleFunc("/dashboard", DashboardHandler).Methods("GET")

	// STATIC FILES (карта, css, изображения)
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static")),
		),
	)

	fmt.Println("STATIC PATH SHOULD BE: ./static")
	fmt.Println("Dashboard: http://localhost:8080/dashboard")
	fmt.Println("Static: http://localhost:8080/static/map.png")

	log.Fatal(http.ListenAndServe(":8080", r))

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

		w.Write([]byte("Its KortiMilli"))

	} else if data.PAN[:1] == "4" {

		w.Write([]byte("Its Visa"))

	} else if data.PAN[:2] == "62" || data.PAN[:2] == "81" {

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

	transaction.Time = time.Now()

	transaction.TimeFormatted = transaction.Time.Format("2006-01-02 15:04:05")

	fmt.Println("Transaction:", transaction)

	//тут проверка антифрода и результат в терминал
	if CheckDroppers(transaction) {

		transaction.Status = "DECLINE"
		transaction.Alert = "⚠ DROP DETECTED"
		transaction.Message = "Неправильный номер карты"

		lastAlert = "⚠ DROP ATTACK DETECTED"

		// оригинальный PAN
		originalPAN := transaction.PAN
		fmt.Println("Original PAN:", originalPAN)

		// генерация токена
		var token string
		for i := 0; i < 16; i++ {
			token += strconv.Itoa(rand.Intn(10))
		}

		fmt.Println("Generated token:", token)

		// alert
		alert := "Попытка Дроп Платежа! Проверить."
		fmt.Println(alert)

		// маскировка

		maskedToken := MaskToken(token)

		transaction.Message = "Dropper detected"

		//тут ответ отправителю
		response := TransactionResponse{

			Status:  "DECLINE",
			Token:   maskedToken,
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

	//тут Геомонитринг в терминал
	country := GeoMonitoring(transaction.IP)

	fmt.Println("Sender IP:", transaction.IP)
	fmt.Println("Sender Country:", country)

	// 👉 ВОТ ЭТО ДОБАВЛЯЕМ
	geoPoints = append(geoPoints, GeoPoint{
		Country: country,
		IP:      transaction.IP,
		Status:  transaction.Status,
	})

	// сохраняем транзакцию для Dashboard
	transactions = append(transactions, transaction)

	GenerateMap(geoPoints)
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

func CardData(w http.ResponseWriter, r *http.Request) {

	pan := r.URL.Query().Get("pan")
	if pan == "" {
		http.Error(w, "missing pan", http.StatusBadRequest)
		return
	}

	var format string

	if len(pan) >= 4 &&
		(pan[:4] == "5440" || pan[:4] == "9860" || pan[:4] == "8600") {
		format = "KortiMilli"

	} else if len(pan) >= 1 && pan[:1] == "4" {
		format = "Visa"

	} else if len(pan) >= 2 &&
		(pan[:2] == "62" || pan[:2] == "81") {
		format = "UnionPay"

	} else {
		format = "Unknown"
	}

	// ищем клиента
	var cardInfo CardInfo
	found := false

	for _, client := range clients {
		if client.Card == pan {

			cardInfo = CardInfo{
				CardNumber: client.Card,
				Format:     format,
				ClientName: client.Name,
				Phone:      client.MobileN,
			}

			found = true
			break
		}
	}

	if !found {
		http.Error(w, "card not found", http.StatusNotFound)
		return
	}

	// история транзакций
	history := []Transaction{}

	for _, tx := range transactions {
		if tx.PAN == pan {
			history = append(history, tx)
		}
	}

	// структура страницы
	page := CardPage{
		CardInfo: cardInfo,
		History:  history,
	}

	tmpl, err := template.ParseFiles("templates/card.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, page)
}

func GeoMonitoring(ip string) string {

	url := fmt.Sprintf(
		"http://ip-api.com/json/%s?fields=country",
		ip,
	)

	resp, err := http.Get(url)
	if err != nil {
		return "UNKNOWN"
	}
	defer resp.Body.Close()

	var geo GeoLocation

	err = json.NewDecoder(resp.Body).Decode(&geo)
	if err != nil {
		return "UNKNOWN"
	}

	return geo.Country
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/dashboard.html")
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {

	if err := GenerateMap(geoPoints); err != nil {
		fmt.Println("GenerateMap error:", err)
	}

	tx := transactions

	// показываем только последние 5
	if len(tx) > 5 {
		tx = tx[len(tx)-5:]
	}

	data := Dashboard{
		FraudScore:   78,
		Velocity:     "HIGH",
		Transactions: len(tx),
		TxList:       tx,
		GeoPoints:    geoPoints,
		Alert:        lastAlert,
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GenerateMap(points []GeoPoint) error {

	const W = 1000
	const H = 500

	dc := gg.NewContext(W, H)

	// фон (простая “карта”)
	dc.SetRGB(0.05, 0.07, 0.12)
	dc.Clear()

	// заголовок
	dc.SetRGB(1, 1, 1)
	dc.DrawString("Geo Monitoring Map", 20, 30)

	// рисуем точки (пока псевдо-расположение)
	for _, p := range points {

		fmt.Println("DEBUG:", p.Country, p.Status)

		var x, y int

		switch p.Country {

		case "Tajikistan":
			x, y = 650, 220

		case "Germany":
			x, y = 500, 170

		case "USA":
			x, y = 180, 180

		case "China":
			x, y = 760, 190

		case "Russia":
			x, y = 620, 130

		case "UAE":
			x, y = 590, 240

		default:
			x, y = 500, 250
		}

		// цвет по статусу
		if p.Status == "SUCCESS" {
			dc.SetRGB(0, 1, 0)
		} else {
			dc.SetRGB(1, 0, 0)
		}

		dc.DrawCircle(float64(x), float64(y), 8)
		dc.Fill()

		// подпись страны
		dc.SetRGB(1, 1, 1)
		dc.DrawString(p.Country, float64(x-20), float64(y+20))
	}

	return dc.SavePNG("static/map.png")

}
