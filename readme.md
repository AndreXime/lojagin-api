# LojaGin API

API RESTful para um sistema de e-commerce simplificado, construída com Go e o framework Gin. O projeto inclui funcionalidades completas de autenticação, gerenciamento de produtos, categorias, e um carrinho de compras funcional.

# Features

-   Autenticação JWT: Sistema completo de registro e login com tokens JWT armazenados em cookies.
-   Gerenciamento de Módulos (CRUD): Operações completas de Criação, Leitura, Atualização e Deleção para:

    -   Usuários: Gerenciamento de perfis de clientes.
    -   Produtos: Cadastro e controle de itens da loja.
    -   Categorias: Organização de produtos em seções.

-   Carrinho de Compras:

    -   Adicionar e remover produtos.
    -   Limpar o carrinho.
    -   Finalizar compra (Checkout).

-   Banco de Dados SQLite: Utilização do GORM para interagir com um banco de dados SQLite, facilitando a configuração e o desenvolvimento.

-   Migrations e Seeding: Sistema de migrações para o esquema do banco de dados e seeding para popular o banco com dados iniciais.

-   Documentação da API com Swagger: Documentação completa e interativa dos endpoints da API, gerada a partir de um arquivo openapi.yaml.

-   Testes End-to-End: Suíte de testes completa para garantir a estabilidade e o correto funcionamento de todos os endpoints da API.

-   Estrutura Modular: O código é organizado em módulos (auth, user, product, category, cart), facilitando a manutenção e a adição de novas funcionalidades.

## Instalação e Setup

Siga os passos abaixo para configurar e rodar o projeto em seu ambiente local.

1.  **Clone o repositório:**

    ```bash
    git clone https://github.com/AndreXime/lojagin-api.git
    cd lojagin-api
    ```

2.  **Instale as dependências:**

    ```bash
    go mod tidy
    ```

3.  **Crie as variáveis de ambiente:**
    Você pode defini-las no seu terminal ou usar um arquivo `.env`.

    ```bash
    export JWT_SECRET="segredo_secreto"
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

## Executando os Testes

Para rodar a suíte completa de testes end-to-end, utilize o seguinte comando:

```bash
make test
```

## Endpoints da API

### Documentação Interativa

Para uma visão completa de todos os endpoints, parâmetros e corpos de requisição, acesse a documentação do Swagger enquanto a aplicação estiver rodando:

[http://localhost:8080/swagger/index.html]()

A API está estruturada com as seguintes rotas:

-   Autenticação:

    -   POST /api/auth/register: Registra um novo usuário.

    -   POST /api/auth/login: Autentica um usuário e retorna um token JWT.

-   Usuários (protegido):

    -   GET /api/users/: Lista todos os usuários.

    -   GET /api/users/:id: Obtém os detalhes de um usuário.

    -   PUT /api/users/:id: Atualiza os dados de um usuário.

    -   DELETE /api/users/:id: Deleta um usuário.

-   Categorias:

    -   GET /api/categories/: Lista todas as categorias.

    -   GET /api/categories/:id: Obtém uma categoria específica.

    -   POST /api/categories/ (protegido): Cria uma nova categoria.

    -   PUT /api/categories/:id (protegido): Atualiza uma categoria.

    -   DELETE /api/categories/:id (protegido): Deleta uma categoria.

-   Produtos:

    -   GET /api/products/: Lista todos os produtos.

    -   GET /api/products/:id: Obtém um produto específico.

    -   POST /api/products/ (protegido): Cria um novo produto.

    -   PUT /api/products/:id (protegido): Atualiza um produto.

    -   DELETE /api/products/:id (protegido): Deleta um produto.

-   Carrinho (protegido):

    -   GET /api/cart/: Visualiza o carrinho.

    -   POST /api/cart/add: Adiciona um item ao carrinho.

    -   POST /api/cart/remove: Remove um item do carrinho.

    -   DELETE /api/cart/clear: Esvazia o carrinho.

    -   POST /api/cart/checkout: Finaliza a compra.
