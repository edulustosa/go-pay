# Go Pay

Este projeto é uma implementação em Go do [desafio back-end do PicPay](https://github.com/PicPay/picpay-desafio-backend?tab=readme-ov-file), que propõe a criação de uma versão simplificada do PicPay, onde os usuários podem realizar depósitos e transferências de dinheiro. O sistema suporta dois tipos de usuários: comuns e lojistas, ambos com suas respectivas carteiras para movimentações financeiras entre si.

A arquitetura foi estruturada em camadas, adotando padrões como o *Repository Pattern* e o *Factory Pattern* para promover flexibilidade e manutenção do código. Além disso, foram implementados testes unitários para garantir a confiabilidade e eficiência dos serviços.

## Como testar

Acesse <https://go-pay.apidocumentation.com/reference> para ter acesso à documentação da API e testar facilmente através do website ou usando seu cliente HTTP preferido.

### Testar localmente

Você pode testar localmente facilmente usando o [Docker](https://docs.docker.com/engine/install/).

- Primeiro clone o repositório:

``` bash
git clone https://github.com/edulustosa/go-pay.git
```

- Dentro do projeto renomeie o arquivo `.env.example` para `.env` e ajuste conforme necessário.

- Suba os containers usando o Docker Compose:

```bash
docker compose up --build
```

- Teste seguindo a [documentação](https://go-pay.apidocumentation.com/reference).

- Você também pode rodar os testes unitários com o seguinte comando:

```bash
go test -v ./...
```

Para isso, será necessário ter o [Go instalado](https://go.dev/doc/install).

## Agradecimentos

Obrigado por verificar meu projeto. Espero que tenha atendido às expectativas do desafio.
