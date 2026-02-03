# Enunciado do Desafio – Eficiência e Performance  
**Cenário do Caso 2**

## Desafio de eficiência e performance 👋

Antes de você se concentrar em resolver o desafio, vale revisar os conceitos centrais de **eficiência** e **performance**.

No contexto de **Site Reliability Engineering (SRE)**, *eficiência* e *performance* referem-se à capacidade de um sistema de executar tarefas de forma ótima e atender à demanda, usando os recursos da maneira mais custo-efetiva possível. Enquanto a **performance** se concentra na velocidade, na capacidade e no tempo de resposta do sistema, a **eficiência** avalia como esses resultados são obtidos em relação aos recursos consumidos — buscando maximizar o que se entrega e minimizar custos. Em outras palavras, **eficiência e performance** tratam de fazer as coisas fluídas e sustentáveis ao mesmo tempo.

**Esses conceitos são fundamentais para uma boa experiência do usuário em fintech.** Quanto maior o tempo de resposta das aplicações, maior a probabilidade de frustração e perda de usuários. É preciso que o sistema não só responda corretamente (availability), mas que responda de forma ágil e eficiente.

---

## Cenário (fintech)

Assim como no desafio anterior, existem duas aplicações interconectadas: a **SRE API** (busca, relatórios, contas e ajustes) e o **fintech-api-failures**, que expõe o catálogo de contas do qual a SRE API depende.

A **SRE API** entrega um **relatório** (`/v1/report`) que permite obter a quantidade de contas disponíveis por tipo (equivalente a “produtos por categoria”). Também permite gerar listas ordenadas, como os **100 primeiros por taxa** (fee). Esse relatório hoje é **ineficiente** em relação ao desejado: sob carga, a quantidade de relatórios que o sistema consegue atender no mesmo intervalo de tempo fica limitada.

## Objetivo

Comece lendo o **README** do repositório deste desafio e execute o script de validação para verificar que a aplicação permite atingir um **máximo de cerca de 15 relatórios** no cenário de prova.

O desafio é fazer com que a quantidade de relatórios atendidos com sucesso seja de **pelo menos 70** no mesmo cenário de validação, **sem aumentar a infraestrutura** da solução.

**Meta quantitativa:** você deve obter **pelo menos 70 relatórios** (requisições ao `/v1/report` concluídas com sucesso) na execução do script `case_2.js`. O desafio é considerado superado quando esse resultado é alcançado. Você poderá acompanhar a quantidade de requisições bem-sucedidas no resultado do script de validação local.

**Importante:** para resolver este desafio é necessário ter resolvido antes o desafio de **disponibilidade** (Caso 1). A solução deve alcançar a meta **sem escalar infraestrutura**.

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
- **Não escalar** infraestrutura nem adicionar recursos. A solução deve alcançar **pelo menos 70 relatórios** **sem alterar a infraestrutura**.
- **Pré-requisito:** ter resolvido o Caso 1 (disponibilidade).
- **Validação 100% local:** não é necessário fazer **deploy**. Codifique e teste tudo localmente; os testes e a validação são inteiramente locais.
- **Script de validação:** para este desafio, **apenas o script `case_2.js`** deve ser considerado para a validação do caso 2. Os outros scripts são usados em outros desafios.
- **Tempo:** o desafio foi planejado para ser resolvido em **dois encontros** (Encontro 3 em diante).
- O foco é aproximar o resultado o máximo possível da meta (≥ 70 relatórios), e não obter um número exato.

---

## Processo de resolução sugerido

1. **Preparação e checagens:** prepare o ambiente e verifique se as validações do caso 2 rodam localmente, com condições comparáveis para todas as provas.
2. **Baseline:** execute o caso de validação (`case_2.js`) e registre quantos relatórios são atendidos com sucesso (e, se útil, success rate, latência) como referência do “antes”.
3. **Hipótese:** enumere alternativas para tornar o relatório mais **eficiente** e melhorar a **performance** (throughput) respeitando as restrições (cache, otimizações, redução de trabalho redundante, etc.); defina critérios simples de comparação.
4. **Experimento mínimo:** escolha uma alternativa e configure-a com parâmetros claros e plano de **rollback**. Priorize capacidades já existentes antes de construir do zero.
5. **Prova sob load e comparação:** repita a validação com o mesmo perfil de **load**. Compare antes/depois e analise **trade-offs** entre throughput, latência, disponibilidade e consistência dos dados.
6. **Documentação:** registre evidência mínima: uma captura ou saída da validação antes/depois (quantidade de relatórios, success rate), uma mensagem/commit curto com a hipótese e a alteração testada, e o comando para reproduzir a mesma **load** usada na prova.

Este desafio treina a capacidade de pensar em termos de **eficiência**, **performance** e **throughput**. É preciso equilibrar quantidade de relatórios atendidos, tempo de resposta e uso de recursos.

---

## Compartilhe sua solução (opcional)

Ao final do desafio, você pode compartilhar no Whatsapp da turma uma captura ou um breve resumo do antes→depois (quantidade de relatórios e, se quiser, success rate / latência), para inspirar outros(as) colegas, contrastar abordagens e explorar diferentes formas de implementar e validar melhorias de eficiência e performance.
