# TicketHub

## Introdução
O desenvolvimento do **TicketHub** surge como uma solução inovadora para facilitar a integração entre companhias aéreas no processo de reserva de passagens. Com o TicketHub, companhias conveniadas podem compartilhar trechos de rotas, permitindo que um cliente reserve passagens de diferentes empresas em uma única transação, como se estivesse utilizando um sistema centralizado. Através de uma API REST projetada especificamente para esse propósito, o TicketHub conecta os servidores das companhias, possibilitando que o cliente selecione trechos de diversas empresas em um só itinerário. Implementado em contêineres Docker para garantir escalabilidade e segurança, o sistema evita os riscos de uma solução centralizada e é projetado para operar com eficiência e resiliência.

## Metodologia

Como estamos lidando com sistemas distribuidos e as requisições de compra de passagens podem chegar de qualquer um dos servidores, precisamos utilizar um método capaz "sincronizar" a ordem das solicitações e lidar com a concorrência. Desta forma foi escolhido o uso de relógios vetoriais. Relogios vetoriais são uma técnica de sincronização utilizada em sistemas distribuídos para ordenar eventos e tratar concorrência. Eles são uma evolução dos relógios lógicos propostos por Lamport, mas com uma abordagem mais robusta para capturar relações de causalidade entre eventos de processos diferentes.

### Conceito Básico de Relógios Vetoriais
Em um sistema distribuído, cada processo (servidor no nosso caso) mantém um vetor de relógios, onde cada posição representa o contador de eventos para um processo específico no sistema. Esse vetor de relógios é atualizado com base nas operações realizadas localmente e nas mensagens trocadas entre processos. Quando ocorre um evento em um processo, o contador desse processo é incrementado em seu vetor de relógios.

### Funcionamento
1. **Eventos Internos**: Quando um processo realiza um evento interno, ele incrementa sua posição no vetor. Isso reflete a passagem de tempo dentro do próprio processo.

2. **Envio de Mensagens**: Quando um processo envia uma mensagem a outro, ele anexa seu vetor de relógios à mensagem. Isso permite ao processo receptor entender a relação causal dos eventos que ocorreram no processo emissor até aquele momento.

3. **Recebimento de Mensagens**: Ao receber uma mensagem, o processo atualiza seu vetor de relógios para refletir o estado do vetor de quem enviou a mensagem. Isso é feito com uma operação de maximização dos valores de cada posição do vetor (elemento a elemento), onde o valor de cada posição é o maior entre o valor atual e o valor recebido. Em seguida, o processo incrementa seu próprio contador local.

### Ordem Parcial e Causalidade
Relógios vetoriais garantem uma **ordem parcial** entre eventos, ou seja, permitem identificar quando um evento A aconteceu antes de um evento B. Essa relação de causalidade é fundamental em sistemas distribuídos, pois ajuda a evitar conflitos de concorrência e a coordenar operações entre processos. Se, para dois eventos A e B, o vetor de A for menor ou igual que o vetor de B (em todas as posições), então A causou B. Caso contrário, os eventos são considerados concorrentes, sem relação causal direta.

### Vantagens e Aplicações
Os relógios vetoriais são especialmente úteis para detectar eventos concorrentes e resolver conflitos, principalmente em sistemas de armazenamento distribuído e controle de versões. Por exemplo, ao registrar e ordenar operações de leitura e escrita de dados, os relógios vetoriais podem assegurar que a execução seja consistente e sem interferências entre processos.

![Diagrama de Concorrência](images/controle-concorrencia)
## Resultados e Discussões

## Conclusão

## Equipe
- [José Gabriel](https://github.com/juserrrrr)
- [Thiago Sena](https://github.com/ThiagoPPSena)
 
## Referências

