package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type CurrencyInfo struct {
	Bid string `json:"bid"`
}

type CurrencyResponse struct {
	USDBRL CurrencyInfo `json:"USDBRL"`
}

func sqlDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		dolar TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func sqlInsert(ctx context.Context, db *sql.DB, cotacao string) error {
	stmt, err := db.Prepare("INSERT INTO cotacoes (dolar) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, cotacao)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Timeout ao inserir a cotação no banco de dados.")
		}
		return err
	}

	fmt.Println("Cotação salva com sucesso no banco de dados")
	return nil
}

func fetchCotacao(ctx context.Context) (CurrencyResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return CurrencyResponse{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Timeout ao buscar a cotação da API externa.")
		}
		return CurrencyResponse{}, err
	}
	defer resp.Body.Close()

	var currency CurrencyResponse
	err = json.NewDecoder(resp.Body).Decode(&currency)
	if err != nil {
		return CurrencyResponse{}, err
	}

	fmt.Println("Cotação do dólar recebida:", currency.USDBRL.Bid)
	return currency, nil
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlDatabase()
	if err != nil {
		http.Error(w, "Erro ao configurar o banco de dados", http.StatusInternalServerError)
		log.Println("Erro ao configurar o banco:", err)
		return
	}
	defer db.Close()

	// Contexto com timeout de 200ms para a API externa
	apiCtx, cancelAPI := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancelAPI()

	currency, err := fetchCotacao(apiCtx)
	if err != nil {
		http.Error(w, "Erro ao buscar a cotação", http.StatusInternalServerError)
		log.Println("Erro ao buscar a cotação:", err)
		return
	}

	// Contexto com timeout de 10ms para inserção no banco
	dbCtx, cancelDB := context.WithTimeout(r.Context(), 10*time.Millisecond)
	defer cancelDB()

	err = sqlInsert(dbCtx, db, currency.USDBRL.Bid)
	if err != nil {
		http.Error(w, "Erro ao salvar a cotação no banco", http.StatusInternalServerError)
		log.Println("Erro ao salvar no banco:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"cotacao": currency.USDBRL.Bid,
	})
}

func main() {
	http.HandleFunc("/cotacao", cotacaoHandler)
	fmt.Println("Servidor rodando na porta 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}
}
