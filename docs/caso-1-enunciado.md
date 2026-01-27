# Enunciado do Desafio – Availability  
**Cenário do Caso 1**

## Desafio de availability 👋

Antes de você se concentrar em resolver o desafio, vale revisar o conceito central de **availability**.

No contexto de **Site Reliability Engineering (SRE)**, *availability* refere-se ao grau em que um sistema é acessível e capaz de executar suas funções necessárias de forma eficaz. Ela costuma ser medida como a porcentagem do tempo em que o sistema está operacional e acessível aos usuários em relação ao tempo total considerado (ou como a porcentagem de vezes em que responde com sucesso em relação ao número de solicitações recebidas). **Availability** é um indicador essencial de desempenho e confiabilidade e deve ser maximizada dentro de limites práticos e econômicos.

**Estar disponível e responder corretamente às solicitações dos usuários é fundamental para o sucesso das soluções fintech.** Não basta só funcionar corretamente; é essencial que o sistema o faça de forma confiável e consistente na maior parte das vezes possível.

---

## Cenário (fintech)

Temos duas aplicações interconectadas: a **SRE API**, que é a principal (expõe contas, relatórios, busca e ajustes tarifários), e o **fintech-api-failures**, do qual ela depende para funcionar corretamente.

O problema é que o **fintech-api-failures** falha em aproximadamente **50% das consultas**, o que afeta diretamente a estabilidade da **SRE API** e, portanto, a experiência do usuário final.

## Objetivo

Melhorar a **availability** (uptime) da **SRE API** para que supere a do serviço do qual depende, **sem modificar a infraestrutura subjacente nem escalar recursos**.

O objetivo é aproximar o **uptime** o máximo possível de 100%, garantindo que a experiência do usuário permaneça fluida e estável mesmo quando o serviço externo falha.

**Meta quantitativa:** você deve obter uma **error rate** inferior a 5%. O desafio é considerado superado quando esse resultado é alcançado, sem afetar de forma significativa o tempo de resposta (**latency** P50/P95). Você poderá acompanhar essas informações no resultado do script de validação local.

Antes de começar, é importante que você **leia o README** do repositório em que trabalhará e **execute o script de validação** localmente.

---

## Método científico: experimentação e validação 👋

Para resolver os tipos de desafios apresentados neste curso, é muito importante ter ou desenvolver um **mindset científico**. Isso pode ser resumido no seguinte processo:

1. **Qual é o problema que quero resolver?**
2. **Replicar o problema experimentalmente.**
3. **Apresentar uma possível solução.**
4. **Verificar se essa solução resolve o problema.**
5. **Caso contrário, volte ao ponto 3.**
6. **Certifique-se de que a solução não crie novos problemas.**

---

## Condições e restrições

- **Não modificar** o serviço externo (fintech-api-failures).
- **Não escalar** infraestrutura nem adicionar recursos. A solução deve alcançar **error rate** inferior a 5% **sem alterar a infraestrutura**.
- **Validação 100% local:** não é necessário fazer **deploy**. Codifique e teste tudo localmente; os testes e a validação são inteiramente locais.
- **Script de validação:** para este desafio, **apenas o script `case_1.js`** deve ser considerado. Os outros scripts serão usados em desafios futuros.
- **Tempo:** o desafio foi planejado para ser resolvido em **dois encontros**.
- O foco é aproximar o resultado o máximo possível da meta (**error rate** inferior a 5%), e não obter um número exato.

---

## Processo de resolução sugerido

1. **Preparação e checagens:** prepare o ambiente e verifique se as validações rodam localmente sem erros, com condições comparáveis para todas as provas.
2. **Baseline:** execute o caso de validação (`case_1.js`) e registre **success rate**, P50/P95 e erros como referência do “antes”.
3. **Hipótese de contingência:** enumere alternativas para elevar a **availability** respeitando as restrições do desafio; defina critérios simples de comparação (impacto, complexidade, reversão).
4. **Experimento mínimo:** escolha uma alternativa e configure-a com parâmetros claros e plano de **rollback**. Priorize capacidades já existentes antes de construir do zero.
5. **Prova sob load e comparação:** repita a validação com o mesmo perfil de **load**. Compare antes/depois e analise **trade-offs** entre **availability** e tempos de resposta (**latency**).
6. **Documentação:** registre evidência mínima: uma captura ou saída da validação antes/depois (**success rate** e P50/P95), uma mensagem/commit curto com a hipótese e a alteração que você testou (ou **flags/toggles**), e o comando para reproduzir a mesma **load** usada na prova.

Este desafio treina a capacidade de pensar em termos de **resilience** e **availability**. É preciso ser criativo, sem cair em soluções forçadas ou custosas.

---

## Compartilhe sua solução (opcional)

Ao final do desafio, você pode compartilhar no Whatsapp da turma uma captura ou um breve resumo do antes→depois, para inspirar outros(as) colegas, contrastar abordagens e explorar diferentes formas de implementar e validar contingências.
