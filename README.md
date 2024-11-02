# TicketHub

## Introdução
O desenvolvimento do TicketHub surge como uma solução inovadora para facilitar a integração entre companhias aéreas no processo de reserva de passagens, permitindo que um cliente reserve bilhetes de diferentes empresas em uma única transação e itinerário. Com o TicketHub, companhias conveniadas podem compartilhar trechos de rotas de maneira transparente, possibilitando a criação de itinerários de múltiplas empresas como se fossem parte de um sistema centralizado. Esse sistema é sustentado por uma API REST projetada especificamente para conectar os servidores das companhias aéreas, garantindo a comunicação segura e eficaz entre elas. Implementado em contêineres Docker, o TicketHub foi projetado com foco em escalabilidade e segurança, aproveitando as vantagens de sistemas distribuídos para evitar os riscos de uma abordagem centralizada. A utilização de relógios vetoriais garante a ordenação correta dos processos de compra, assegurando a consistência e a resolução de conflitos em um ambiente com múltiplos servidores interagindo simultaneamente

## Metodologia

### Arquitetura
Cada servidor possui uma arquitetura bem modularizada e semelhante. O diretório collections armazena as interfaces utilizadas por outros módulos do sistema. O diretório files contém os arquivos JSON que garantem a persistência dos dados. O módulo graphs é responsável pela leitura dos trechos aéreos nos arquivos JSON e pelo cálculo das rotas possíveis entre cidades. O módulo queues gerencia a fila de solicitações de compras do servidor. O módulo vectorClock cuida da concorrência nas compras em um sistema distribuído, utilizando relógios vetoriais (tópico que será abordado na seção de concorrência).

Por fim, o módulo passages é dividido em três camadas. A camada routes é responsável por estabelecer as rotas HTTP para a comunicação entre dois servidores ou entre o servidor e o cliente. A camada controllers intercepta as requisições, realiza o tratamento adequado e as encaminha para a camada inferior, chamada services. Esta última, services, cuida do processamento dos dados e retorna as respostas para as camadas superiores.

![arquitetura](https://github.com/user-attachments/assets/f9c9fd65-97b8-4b8a-b818-7eee5d3f9f74)

### Protocolo de comunicação
A api desenvolvida para comunicação entre servidores e clientes utiliza métodos HTTP. Os métodos utilizados são GET para pegar os trechos disponíveis e o POST para requisição de compra.

#### Rotas
- http://{host}:{porta do servidor}/passages/routes
- http://{host}:{porta do servidor}/passages/flights
- http://{host}:{porta do servidor}/passages/buy
  
### Roteamento
Quando o usuário seleciona a origem e o destino em um servidor, este servidor realiza uma solicitação para obter todos os trechos disponíveis nos outros servidores por meio de um método HTTP GET. Em seguida, ele combina esses trechos com os próprios trechos da companhia. Dessa forma, é gerado um grafo abrangente que inclui todos os trechos de todas as companhias aéreas. O sistema, então, aplica um algoritmo de busca para identificar as rotas possíveis entre a origem e o destino selecionados.

### Concorrência
Como estamos lidando com sistemas distribuidos e as requisições de compra de passagens podem chegar de qualquer um dos servidores, precisamos utilizar um método capaz "sincronizar" a ordem das solicitações e lidar com a concorrência. Desta forma foi escolhido o uso de relógios vetoriais. Relógios vetoriais são uma técnica de sincronização utilizada em sistemas distribuídos para ordenar eventos e tratar concorrência. Eles são uma evolução dos relógios lógicos propostos por Lamport (LAMPORT, 1978), mas com uma abordagem mais robusta para capturar relações de causalidade entre eventos de processos diferentes.

#### Conceito Básico de Relógios Vetoriais
Em um sistema distribuído, cada processo mantém um vetor de relógios, onde cada posição representa o contador de eventos para um processo específico no sistema. Esse vetor de relógios é atualizado com base nas operações realizadas localmente e nas mensagens trocadas entre processos. Quando ocorre um evento em um processo, o contador desse processo é incrementado em seu vetor de relógios (COULOURIS et al., 2011).

#### Funcionamento
1. **Eventos Internos**: Quando um processo realiza um evento interno, ele incrementa sua posição no vetor. Isso reflete a passagem de tempo dentro do próprio processo.

2. **Envio de Mensagens**: Quando um processo envia uma mensagem a outro, ele anexa seu vetor de relógios à mensagem. Isso permite ao processo receptor entender a relação causal dos eventos que ocorreram no processo emissor até aquele momento.

3. **Recebimento de Mensagens**: Ao receber uma mensagem, o processo atualiza seu vetor de relógios para refletir o estado do vetor de quem enviou a mensagem. Isso é feito com uma operação de maximização dos valores de cada posição do vetor (elemento a elemento), onde o valor de cada posição é o maior entre o valor atual e o valor recebido. Em seguida, o processo incrementa seu próprio contador local (BERNSTEIN; HADZILACOS; GOODMAN, 1987).

#### Ordem Parcial e Causalidade
Relógios vetoriais garantem uma **ordem parcial** entre eventos, ou seja, permitem identificar quando um evento A aconteceu antes de um evento B. Essa relação de causalidade é fundamental em sistemas distribuídos, pois ajuda a evitar conflitos de concorrência e a coordenar operações entre processos. Se, para dois eventos A e B, o vetor de A for menor que o vetor de B (em todas as posições), então A causou B. Caso contrário, os eventos são considerados concorrentes, sem relação causal direta (COULOURIS et al., 2011).

#### Vantagens e Aplicações
Os relógios vetoriais são especialmente úteis para detectar eventos concorrentes e resolver conflitos, principalmente em sistemas de armazenamento distribuído e controle de versões. Por exemplo, ao registrar e ordenar operações de leitura e escrita de dados, os relógios vetoriais podem assegurar que a execução seja consistente e sem interferências entre processos.

#### Uso no sistema 
O sistema utiliza relógios vetoriais para ordenar os processos na fila de compras de cada servidor. Na imagem abaixo, mostramos o que ocorre com um servidor (A). Quando uma solicitação de compra chega a esse servidor (seja de outro servidor ou de um cliente), ela é associada a um relógio vetorial. Dessa forma, o servidor pode comparar o relógio dessa solicitação com os dos demais itens na fila, ordenando-os corretamente. A rotina responsável pelo processamento da compra retira a primeira solicitação da fila, a processa e retorna uma resposta com o status da operação de compra daquela solicitação.

![Diagrama de Concorrência](images/controle-concorrencia)

### Confiabilidade (Falta essa parte)
## Resultados e Discussões

## Conclusão
O TicketHub provou ser uma solução eficiente e escalável para integrar companhias aéreas em um processo de venda de passagens que envolve múltiplos servidores e empresas. A implementação de relógios vetoriais se mostrou crucial para a ordenação correta dos eventos de compra, permitindo que as requisições sejam tratadas de maneira sincronizada, mesmo em um ambiente distribuído. Esse controle preciso de concorrência garante que os clientes possam realizar reservas de diferentes companhias em um único itinerário, melhorando a experiência de compra e ampliando as possibilidades de itinerário. Com isso, o TicketHub se posiciona como uma ferramenta eficaz para a gestão de reservas em um sistema distribuido, oferecendo uma infraestrutura capaz de crescer e se adaptar às demandas do setor aéreo.

## Equipe
- [José Gabriel](https://github.com/juserrrrr)
- [Thiago Sena](https://github.com/ThiagoPPSena)
 
## Referências
- COULOURIS, George; DOLLIMORE, Jean; KINDBERG, Tim; BLAIR, Gordon. Distributed Systems: Concepts and Design. 5th ed. Boston: Addison-Wesley, 2011.
- LAMPORT, Leslie. Time, Clocks, and the Ordering of Events in a Distributed System. Communications of the ACM, v. 21, n. 7, p. 558-565, 1978.
- BERNSTEIN, Philip A.; HADZILACOS, Vassos; GOODMAN, Nathan. Concurrency Control and Recovery in Database Systems. Boston: Addison-Wesley, 1987.
