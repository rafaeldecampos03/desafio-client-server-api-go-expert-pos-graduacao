package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type CotacaoDolar struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	//root:<SENHA-DB>@tcp(<HOST>:<PORTA>)/<NOME-DB>
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS rates (id INT AUTO_INCREMENT PRIMARY KEY, rate VARCHAR(10), timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
		defer cancel()

		rate, err := fetchCotacaoDolar(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = logCotacaoDolar(db, rate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Cotação do dólar")
		fmt.Println(rate)
		//mapping chave e valor da resposta
		cotacaoDolarResponse := map[string]string{"dolar": rate}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cotacaoDolarResponse)
	})

	fmt.Println("Server is running on port 8080...")
	panic(http.ListenAndServe(":8080", nil))
}

func fetchCotacaoDolar(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var cotacaoDolar CotacaoDolar
	if err := json.NewDecoder(resp.Body).Decode(&cotacaoDolar); err != nil {
		return "", err
	}

	return cotacaoDolar.USDBRL.Bid, nil
}

func logCotacaoDolar(db *sql.DB, rate string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := db.ExecContext(ctx, "INSERT INTO rates (rate) VALUES (?)", rate)
	return err
}
