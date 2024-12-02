# GO Routines == Leilão

## Descrição
Este projeto Go que fechar automaticamente um leilão a partir de um tempo definido.

## Para Testar
Para subir os contêineres, utilize o seguinte comando:

docker-compose up --build

Isso irá construir a imagem da aplicação e iniciar tanto a aplicação quanto o MongoDB

## Acessando a Aplicação
Uma vez que os contêineres estejam em execução, você poderá acessar a aplicação na porta 8080 em http://localhost:8080.

## Testes
Os testes automatizados também podem ser executados dentro do ambiente Docker. Caso desejado, considere adicionar um serviço de teste no docker-compose.yml ou executar os testes manualmente com os contêineres em execução.

Para acessar os logs dos contêineres e diagnosticar problemas:

docker-compose logs
