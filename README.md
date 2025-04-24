Desafio de Cotação de Dólar

Este projeto consiste em dois sistemas em Go: client.go e server.go. O servidor faz uma requisição para uma API externa para obter a cotação do dólar, 
armazena a cotação em um banco de dados SQLite e retorna o valor para o cliente. O cliente recebe a cotação e a salva em um arquivo de texto.
Estrutura do Projeto

/client-server-api
    ├── /client
    │   └── client.go          # Código do cliente
    ├── /server
    │   └── server.go          # Código do servidor
    └── README.md              # Este arquivo

Tecnologias Usadas

    Go (versão 1.18 ou superior)

    SQLite (para armazenar as cotações)

    Context (para controle de tempo limite nas requisições e banco de dados)

    HTTP (para comunicação entre cliente e servidor)

Pré-requisitos

Certifique-se de ter o Go e o SQLite instalados em sua máquina.
Instalar Go

    Baixar Go

Instalar SQLite

Para instalar o SQLite, siga as instruções no site oficial.
Como Rodar o Projeto
1. Executar o Servidor

Primeiro, inicie o servidor. Ele vai consumir a API externa para pegar a cotação do dólar e armazená-la no banco de dados SQLite.

    Navegue até o diretório do servidor:

cd client-server-api/server

    Execute o servidor com o comando:

go run server.go

O servidor estará disponível na porta 8080.
2. Executar o Cliente

Agora, execute o cliente para fazer a requisição ao servidor, obter a cotação do dólar e salvar em um arquivo cotacao.txt.

    Navegue até o diretório do cliente:

cd client-server-api/client

    Execute o cliente com o comando:

go run client.go

O cliente fará uma requisição ao servidor e salvará a cotação no arquivo cotacao.txt.
Funcionalidade

    Servidor:

        O servidor responde a requisições HTTP no endpoint /cotacao com a cotação atual do dólar (campo "bid").

        O servidor realiza a inserção da cotação recebida em um banco de dados SQLite.

        Os contextos são utilizados para garantir que os tempos de resposta da API externa e a inserção no banco de dados respeitem os limites de tempo definidos (200ms para a API e 10ms para o banco).

    Cliente:

        O cliente realiza uma requisição HTTP ao servidor para buscar a cotação do dólar.

        Se o tempo de resposta for superior a 300ms, o cliente retorna um erro de timeout.

        O cliente salva a cotação recebida em um arquivo cotacao.txt no formato: Dólar: {valor}.

Observações

    Banco de Dados: O banco de dados SQLite será criado automaticamente na primeira execução do servidor e as cotações serão inseridas na tabela cotacoes.

    Timeouts:

        O servidor tem um timeout de 200ms para chamar a API de cotação.

        O servidor tem um timeout de 10ms para inserir a cotação no banco de dados.

        O cliente tem um timeout de 300ms para receber a resposta do servidor.

Possíveis Melhorias

    Autenticação: O sistema poderia ser melhorado com autenticação, garantindo que apenas clientes autorizados possam acessar a cotação.

    Persistência em Arquivo: Em vez de usar o banco de dados SQLite, o cliente poderia armazenar as cotações em arquivos, se necessário.

    Validação de Dados: Melhorar a validação de dados ao lidar com a resposta da API externa.
