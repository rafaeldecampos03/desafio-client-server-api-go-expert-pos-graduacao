package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type CotacaoDolar struct {
	Dolar string `json:"dolar"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		fmt.Println("Erro ao criar requisição:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao realizar requisição:", err)
		return
	}
	//fechar a conexão
	defer resp.Body.Close()
	var responseData []byte
	if responseData, err = io.ReadAll(resp.Body); err != nil {
		panic(err)
	}
	//Cotação do dolar
	var cotacaoDolar CotacaoDolar
	err = json.Unmarshal(responseData, &cotacaoDolar)
	if err != nil {
		fmt.Println("Erro ao deserializar JSON:", err)
		return
	}
	fmt.Println("Cotação do dolar:", cotacaoDolar)
	cotacao := fmt.Sprintf("Dólar: %s", cotacaoDolar.Dolar)
	if err := os.WriteFile("cotacao.txt", []byte(cotacao), 0644); err != nil {
		fmt.Println("erro ao escrever no arquivo:", err)
		return
	}
	fmt.Println("Cotação salva no arquivo cotacao.txt")
}
