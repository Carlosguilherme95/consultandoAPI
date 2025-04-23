package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type CurrencyInfo struct {
	Bid string `json:"bid"`
}

type CurrencyResponse struct {
	USDBRL CurrencyInfo `json:"USDBRL"`
}

// Função para configurar o banco de dados e criar a tabela
func sqlDatabase() (*sql.DB, error) {
	// Configuração do banco de dados SQLite
	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		return nil, err
	}
	// Criação da tabela se ela não existir
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacoes 
	( id INTEGER PRIMARY KEY AUTOINCREMENT, 
	 dolar TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP)`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Função para inserir a cotação no banco de dados
func sqlInsert(db *sql.DB, cotacao string) error {
	// Preparando o comando SQL para inserir a cotação no banco
	stmt, err := db.Prepare("INSERT INTO cotacoes (dolar) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Executando o comando SQL
	_, err = stmt.Exec(cotacao)
	if err != nil {
		return err
	}

	fmt.Println("Cotação salva com sucesso no banco de dados")
	return nil
}

// Função para fazer a requisição HTTP e pegar a cotação
func fetchCotacao() (CurrencyResponse, error) {
	// Fazendo a requisição para a API de economia
	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		return CurrencyResponse{}, err
	}
	defer req.Body.Close()

	var currency CurrencyResponse
	// Decodificando o JSON retornado pela API
	err = json.NewDecoder(req.Body).Decode(&currency)
	if err != nil {
		return CurrencyResponse{}, err
	}

	// Exibindo a cotação no console
	fmt.Println("Cotação do dólar: ", currency.USDBRL.Bid)
	return currency, nil
}
func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlDatabase()
	if err != nil {
		http.Error(w, "Erro ao configurar o banco de dados", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer db.Close()

	// Obtendo a cotação da API
	currency, err := fetchCotacao()
	if err != nil {
		http.Error(w, "Erro ao buscar a cotação", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Inserindo a cotação no banco de dados
	err = sqlInsert(db, currency.USDBRL.Bid)
	if err != nil {
		http.Error(w, "Erro ao salvar a cotação no banco", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Respondendo ao cliente com sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Retornando a cotação para o cliente
	response := map[string]string{
		"cotacao": currency.USDBRL.Bid,
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	// inicializando servidor HTTP
	http.HandleFunc("/cotacao", cotacaoHandler)
	fmt.Println("Servidor rodando na porta 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
	// Inicializando o banco de dados
	db, err := sqlDatabase()
	if err != nil {
		log.Fatal("Erro ao configurar o banco de dados: ", err)
	}
	defer db.Close()

	// Obtendo a cotação da API
	currency, err := fetchCotacao()
	if err != nil {
		log.Fatal("Erro ao buscar a cotação: ", err)
	}

	// Inserindo a cotação no banco de dados
	err = sqlInsert(db, currency.USDBRL.Bid)
	if err != nil {
		log.Fatal("Erro ao salvar a cotação no banco: ", err)
	}
}
