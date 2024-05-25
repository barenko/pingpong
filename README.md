# PingPong

Documentação pode ser encontrada [aqui](docs/pingpong.pdf).
Os testes de performance, [aqui](docs/performance_results.pdf).

> Os testes de performance devem ser analisados de forma comparativa e não absoluta.

Caso queira, utilize os facilitadores que estão configurados no [Makefile](Makefile).

## Compilar

    go build

## Executar com um servidor

    ./pingpong &

    http :3000/ping/4

## Executar com dois servidores

    PORT=3000 PONG=http://0.0.0.0:3001 ./pingpong &
    PORT=3001 PING=http://0.0.0.0:3000 ./pingpong &

    http :3000/ping/4
