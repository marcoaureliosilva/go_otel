# Sistema de Temperatura por CEP

Sistema em Go que recebe um CEP, identifica a cidade e retorna o clima atual. É utilizado o OTEL (OpenTelemetry) e Zipkin para tracing distribuído.

## Passo a Passo para Rodar a Aplicação
Subindo os serviços com Docker:
docker-compose up --build

Com os serviços rodando, você pode testar a API do Serviço A, que encaminhará a solicitação para o Serviço B.
Exemplo de teste usando curl:
curl -X POST http://localhost:8080/cep -d '{"cep": "29902555"}' -H 'Content-Type: application/json'

## Testes
Certifique-se de ter o Zipkin rodando para visualizar os traces. Você pode executar o Zipkin localmente via Docker com o seguinte comando:
docker run -d -p 9411:9411 openzipkin/zipkin
Depois de rodar os serviços, acesse a interface do Zipkin em http://localhost:9411 para visualizar os traces.
