# Open Telemetry - Golang Expert

Para rodar a aplicação use o docker-compose com o comando abaixo:

```
docker-compose up -d
```

Para acessar a rota do servico, utilize alguem aplicativo para fazer um `POST` no seguinte endereço:

```
http://localhost:8080
```

Use um payload `JSON` com um cep, sem espaços, hifem ou pontuação: como no exemplo abaixo:

```
{
  "cep": "45208643"
}
```

Para acessar o zapkin basta usar o endereço:

```
http://localhost:9411/zipkin
```