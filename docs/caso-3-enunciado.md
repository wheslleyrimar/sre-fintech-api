# Enunciado do Desafio – Integridade e Consistência  
**Cenário do Caso 3**

## Exercício de Integridade e Consistência 👋

Vamos continuar aprofundando o desenvolvimento de aplicações confiáveis. Neste encontro você começará um novo desafio para trabalhar os conceitos de **Integridade** e **Consistência**.

Há sistemas que gerenciam informação crítica: bases de dados, sistemas financeiros, sistemas de controle industrial. Nesses contextos, consistência e integridade são fundamentais. Não é aceitável, por exemplo, que o saldo de uma conta mude sem que existam movimentos correspondentes. Da mesma forma, é crítico garantir que os dados permaneçam íntegros depois de gerados.

### O que é Integridade e Consistência?

No marco de **Site Reliability Engineering (SRE)**:

- **Integridade** refere-se à imutabilidade e à proteção dos dados e dos sistemas contra alterações não autorizadas, garantindo que operem conforme o esperado.
- **Consistência** garante que, em ambientes distribuídos ou com múltiplos processos, a informação permaneça uniforme e sem discrepâncias entre o que foi solicitado e o que está efetivamente armazenado ou exposto.

Juntas, essas qualidades são essenciais para manter a confiabilidade e a previsibilidade dos sistemas em ambientes complexos.

---

## Contexto (fintech)

Uma das funções centrais do nosso produto é a **atualização de tarifas** (taxa mensal, *monthly_fee*) que os titulares das contas solicitam para suas contas. Quando um deles solicita um **ajuste de tarifa**, o pedido é processado de forma assíncrona: uma aplicação (backend) valida e aplica o novo valor e, em paralelo, o sistema registra o histórico de ajustes da conta. Depois que o processamento é concluído, a consulta à conta deve refletir o **último valor de tarifa** solicitado e aprovado.

Ou seja: o fluxo envolve **uma conta** por vez; o desafio é garantir que, após um ajuste de tarifa ser aceito (POST com sucesso), a leitura da conta (GET) mostre sempre o **mesmo valor** que foi solicitado — ou seja, **consistência** entre o que foi pedido e o que está persistido.

---

## Desafio a resolver

O ajuste de tarifa é sobre **uma conta** por requisição. Comece lendo o **README** da aplicação deste desafio e execute o script de validação para verificar que **alguns ajustes de tarifa ficam consistentes e outros não**.

O desafio é alcançar **100% de consistência** com um **throughput** (rendimento) de **mais de 25 ajustes de tarifa por minuto**.

Em outras palavras:

- **Consistência:** sempre que um POST em `/v1/accounts/{id}/tariff-adjustments` for aceito (status 200 ou 204), uma leitura posterior da conta em `/v1/accounts/{id}` deve retornar um `monthly_fee` igual ao `new_fee` enviado no ajuste.
- **Throughput:** o script de validação (`case_3.js`) executa um cenário de carga por 60 segundos; é necessário manter um ritmo de mais de 25 ajustes de tarifa por minuto nesse cenário.

**Importante:** para resolver este desafio **não é permitido aumentar a infraestrutura** da solução. É preciso resolver mantendo os processos em paralelo (ou as decisões de arquitetura já existentes). **Não é permitido alterar o script de validação** (`case_3.js`).

---

## Método científico: experimentação e validação 👋

Para resolver este tipo de problema é fundamental ter ou desenvolver um **mindset científico**. O processo pode ser resumido assim:

1. **Qual é o problema que quero resolver?**
2. **Replique o problema de forma experimental.**
3. **Proponha uma possível solução.**
4. **Verifique se essa solução resolve o problema.**
5. **Se não resolver, volte ao ponto 3.**
6. **Garanta que a solução não introduza novos problemas.**

Em cada encontro, a cooperação entre pares será fundamental para chegar à solução do exercício do dia. Se surgir algum obstáculo, você poderá recorrer à ajuda de um mentor, ferramentas de busca, ChatGPT e demais recursos que considerar necessários.

---

## Condições e restrições

- **Não aumentar a infraestrutura** da solução (sem escalar recursos, sem adicionar novos servidores ou instâncias).
- **Manter o processamento em paralelo** conforme o desenho atual; a solução deve garantir consistência dentro desses limites.
- **Não manipular o script de validação** (`case_3.js`). A validação do Caso 3 é feita exclusivamente por esse script.
- **Validação local:** o desafio é validado localmente com o script de validação (k6) e os binários gerados por `./install.sh`, conforme descrito no README do repositório.

---

## Meta quantitativa

- **Consistência:** 100% das requisições de ajuste de tarifa aceitas (POST 200/204) devem resultar em conta com `monthly_fee` igual ao `new_fee` solicitado quando a conta for consultada (GET).
- **Throughput:** mais de 25 ajustes de tarifa por minuto no cenário de 60 segundos do `case_3.js`.

O desafio é considerado superado quando o script `case_3.js` atinge **100% de checks de consistência** (“account monthly_fee matches requested new_fee”) e o throughput permanece acima de 25 ajustes por minuto.

---

## Processo de resolução sugerido

1. **Preparação:** leia o README do repositório, execute `./install.sh` e rode o script de validação do Caso 3 para reproduzir o problema (inconsistências e throughput atual).
2. **Baseline:** anote quantos ajustes ficam consistentes e qual o throughput atual (antes da solução).
3. **Hipótese:** identifique por que alguns ajustes não refletem o `new_fee` na leitura da conta (condições de corrida, ordem de eventos, falta de sincronização, etc.) e proponha uma solução que preserve o paralelismo e não aumente a infraestrutura.
4. **Implementação e validação:** implemente a solução e execute novamente o `case_3.js` até alcançar 100% de consistência e throughput > 25 ajustes/min.
5. **Registro:** documente a hipótese, a alteração feita e um resultado de validação (antes/depois) para evidência.

Este desafio treina a capacidade de garantir **integridade** e **consistência** em um fluxo assíncrono de atualização de tarifas, sem sacrificar o rendimento do sistema.
