# LojaGin API

API RESTful para um sistema de e-commerce simplificado, construída com Go e o framework Gin. O projeto inclui funcionalidades de autenticação de usuários, gerenciamento de perfis e uma estrutura modular para futuras expansões.

## Features

-   **Autenticação JWT**: Sistema completo de registro e login com tokens JWT em cookies.
-   **Gerenciamento de Usuários (CRUD)**: Operações completas de Criação, Leitura, Atualização e Deleção para os usuários da plataforma.
-   **Banco de Dados SQLite**: Utilização do GORM para interagir com um banco de dados SQLite, facilitando a configuração e o desenvolvimento.
-   **Documentação da API com Swagger**: Documentação completa e interativa dos endpoints da API, gerada a partir de um arquivo `openapi.yaml`.
-   **Testes End-to-End**: Suíte de testes completa para garantir a estabilidade e o correto funcionamento de todos os endpoints da API.
-   **Estrutura Modular**: O código é organizado em módulos (auth, user), facilitando a manutenção e a adição de novas funcionalidades.

## Tecnologias Utilizadas

-   **Go**: Linguagem de programação principal.
-   **Gin**: Framework web para construção da API.
-   **GORM**: ORM para interação com o banco de dados.
-   **SQLite**: Banco de dados relacional.
-   **JWT-go**: Para geração e validação de tokens JWT.
-   **Testify**: Biblioteca para asserções em testes.
-   **Swagger**: Para documentação da API.

## Instalação e Setup

Siga os passos abaixo para configurar e rodar o projeto em seu ambiente local.

1.  **Clone o repositório:**

    ```bash
    git clone https://github.com/andrexime/lojagin.git
    cd lojagin
    ```

2.  **Instale as dependências:**

    ```bash
    go mod tidy
    ```

3.  **Crie as variáveis de ambiente:**
    Você pode defini-las no seu terminal ou usar um arquivo `.env`.

    ```bash
    export JWT_SECRET="seu_segredo_super_secreto"
    export DB_URL="lojagin.db"
    ```

4.  **Rode a aplicação em modo de desenvolvimento:**
    O projeto vem com `air` configurado para hot-reloading. Para iniciar, use o comando do `makefile`:

    ```bash
    make dev
    ```

    O servidor estará disponível em `http://localhost:8080`.

5.  **Build da aplicação:**
    Para compilar o binário da aplicação, utilize:

    ```bash
    make build
    ```

    O executável será gerado no diretório `bin/`.

## Endpoints da API

A API está estruturada com as seguintes rotas:

-   `POST /api/auth/register`: Registra um novo usuário.
-   `POST /api/auth/login`: Autentica um usuário e retorna um token JWT.
-   `GET /api/users/`: Lista todos os usuários (rota protegida).
-   `GET /api/users/:id`: Obtém os detalhes de um usuário específico (rota protegida).
-   `PUT /api/users/:id`: Atualiza os dados de um usuário (rota protegida).
-   `DELETE /api/users/:id`: Deleta um usuário (rota protegida).

### Documentação Interativa

Para uma visão completa de todos os endpoints, parâmetros e corpos de requisição, acesse a documentação do Swagger enquanto a aplicação estiver rodando:

[http://localhost:8080/swagger/index.html]()

## Executando os Testes

Para rodar a suíte completa de testes end-to-end, utilize o seguinte comando:

```bash
make test
```
